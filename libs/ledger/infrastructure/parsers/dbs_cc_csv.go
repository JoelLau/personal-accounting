package parsers

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"libs/ledger/application/commands"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

const DBSDateLayout = "02 Jan 2006"

type DbsCreditCardCsvParser struct{}

func NewDbsCreditCardCsvParser() *DbsCreditCardCsvParser {
	return &DbsCreditCardCsvParser{}
}

func (p *DbsCreditCardCsvParser) Parse(file io.Reader) ([]commands.RawTransaction, error) {
	records, err := p.extract(file)
	if err != nil {
		return nil, fmt.Errorf("error parsing raw rows: %w", err)
	}

	if len(records) == 0 {
		return nil, errors.New("no transactions")
	}

	return p.normalize(records), nil
}

func (p *DbsCreditCardCsvParser) extract(r io.Reader) ([]DbsTransaction, error) {
	csvReader := csv.NewReader(r)
	csvReader.FieldsPerRecord = -1

	var (
		transactions     []DbsTransaction
		headerMap        map[string]int
		currCard         string
		expectingCardNum bool
	)

	for {
		record, err := csvReader.Read()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, fmt.Errorf("csv read error: %w", err)
		}

		if len(record) == 0 || record[0] == "" {
			continue
		}

		firstCol := strings.TrimSpace(record[0])

		if expectingCardNum {
			currCard = firstCol
			expectingCardNum = false
			continue
		}

		switch {
		case firstCol == "Supplementary Card:", firstCol == "Blocked Card:":
			expectingCardNum = true
		case strings.Contains(firstCol, "Card Transaction Details For:"):
			if len(record) > 1 && record[1] != "" {
				currCard = record[1]
			}
		case firstCol == "Transaction Date":
			headerMap = make(map[string]int)
			for i, name := range record {
				headerMap[name] = i
			}
		case firstCol == "No transactions to view",
			firstCol == "We didn't find any transactions for the selected date range.":
			continue
		default:
			if headerMap == nil {
				continue
			}

			date, err := parseDbsDate(getCol(record, headerMap, "Transaction Date"))
			if err != nil {
				return nil, fmt.Errorf("error parsing transaction date: %w", err)
			}
			postingDate, err := parseDbsDate(getCol(record, headerMap, "Transaction Posting Date"))
			if err != nil {
				return nil, fmt.Errorf("error parsing transaction posting date: %w", err)
			}
			debitAmount, err := parseNullDecimal(getCol(record, headerMap, "Debit Amount"))
			if err != nil {
				return nil, fmt.Errorf("error parsing debit amount: %w", err)
			}
			creditAmount, err := parseNullDecimal(getCol(record, headerMap, "Credit Amount"))
			if err != nil {
				return nil, fmt.Errorf("error parsing credit amount: %w", err)
			}

			tx := DbsTransaction{
				CardName:     currCard,
				Date:         date,
				PostingDate:  postingDate,
				Description:  getCol(record, headerMap, "Transaction Description"),
				Type:         getCol(record, headerMap, "Transaction Type"),
				PaymentType:  getCol(record, headerMap, "Payment Type"),
				Status:       getCol(record, headerMap, "Status"),
				DebitAmount:  debitAmount,
				CreditAmount: creditAmount,
				RawRow:       strings.Join(append([]string{currCard}, record...), "|"),
			}
			transactions = append(transactions, tx)

		}
	}

	return transactions, nil
}

func (p *DbsCreditCardCsvParser) normalize(raw []DbsTransaction) []commands.RawTransaction {
	transactions := make([]commands.RawTransaction, len(raw))

	for idx, row := range raw {
		var amount int64 = 0

		if row.DebitAmount.Valid && row.DebitAmount.Decimal.GreaterThan(decimal.Zero) {
			amount = row.DebitAmount.Decimal.Mul(decimal.NewFromInt(1_000_000)).IntPart()
		}

		if row.CreditAmount.Valid && row.CreditAmount.Decimal.GreaterThan(decimal.Zero) {
			amount = row.CreditAmount.Decimal.Mul(decimal.NewFromInt(-1_000_000)).IntPart()
		}

		transactions[idx] = commands.RawTransaction{
			ID:          row.RawRow,
			SourceName:  row.CardName,
			Date:        row.Date,
			Description: row.Description,
			Amount:      amount,
		}
	}
	return transactions
}

// NOTE: struct tags are for personal reference (not used in code)
type DbsTransaction struct {
	CardName     string              `csv-h:"Card Transaction Details For:" csv-v:"Supplementary Card:,Blocked Card:"` // e.g. "DBS Vantage Visa Infinite Card 4119-1100-1234-1234"
	Date         time.Time           `csv:"Transaction Date"`                                                          // e.g. "17 Jan 2026"
	PostingDate  time.Time           `csv:"Transaction Posting Date"`                                                  // e.g. "21 Jan 2026"
	Description  string              `csv:"Transaction Description"`                                                   // e.g. "BUS/MRT 781915832 SINGAPORE SG"
	Type         string              `csv:"Transaction Type"`                                                          // e.g. "PURCHASE"
	PaymentType  string              `csv:"Payment Type"`                                                              // e.g. "Contactless"
	Status       string              `csv:"Status"`                                                                    // e.g. "Settled"
	DebitAmount  decimal.NullDecimal `csv:"Debit Amount"`                                                              // e.g. "2.98"
	CreditAmount decimal.NullDecimal `csv:"Credit Amount"`                                                             // e.g. ""
	RawRow       string
}

func parseNullDecimal(s string) (decimal.NullDecimal, error) {
	s = strings.TrimSpace(s)

	s = strings.ReplaceAll(s, "SGD", "")
	s = strings.ReplaceAll(s, ",", "")
	s = strings.TrimSpace(s)

	if s == "" {
		return decimal.NullDecimal{
			Decimal: decimal.Zero,
			Valid:   false,
		}, nil
	}

	d, err := decimal.NewFromString(s)
	if err != nil {
		return decimal.NullDecimal{
			Decimal: decimal.Zero,
			Valid:   false,
		}, fmt.Errorf("invalid string: %w", err)
	}

	return decimal.NewNullDecimal(d), nil
}

func parseDbsDate(s string) (time.Time, error) {
	t, err := time.Parse(DBSDateLayout, strings.TrimSpace(s))
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing time in dbs layout: %w", err)
	}
	return t, nil
}

func getCol(row []string, m map[string]int, key string) string {
	if idx, ok := m[key]; ok && idx < len(row) {
		return row[idx]
	}
	return ""
}
