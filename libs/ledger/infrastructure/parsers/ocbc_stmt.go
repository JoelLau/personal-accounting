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

const (
	OCBCDateLayout         = "2/1/06"
	OCBCFullYearDateLayout = "2/1/2006"
)

type OcbcStatementCsvParser struct{}

func NewOcbcStatementCsvParser() *OcbcStatementCsvParser {
	return &OcbcStatementCsvParser{}
}

func (p *OcbcStatementCsvParser) Parse(file io.Reader) ([]commands.RawTransaction, error) {
	records, err := p.extract(file)
	if err != nil {
		return nil, fmt.Errorf("failed to parse raw rows: %w", err)
	}

	if len(records) == 0 {
		return nil, errors.New("no transactions")
	}

	return p.normalize(records)
}

func (p *OcbcStatementCsvParser) extract(r io.Reader) ([]OcbcTransaction, error) {
	csvReader := csv.NewReader(r)
	csvReader.FieldsPerRecord = -1

	var (
		transactions     []OcbcTransaction
		headerMap        map[string]int
		accountName      string
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
			accountName = firstCol
			expectingCardNum = false
			continue
		}

		switch {
		case strings.Contains(firstCol, "Account details for:"):
			if len(record) > 1 && record[1] != "" {
				accountName = record[1]
			}
		case firstCol == "Transaction date":
			headerMap = make(map[string]int)
			for i, name := range record {
				headerMap[name] = i
			}
		case firstCol == "Transaction History",
			firstCol == "Available Balance",
			firstCol == "Ledger Balance:":
			continue

		default:
			if headerMap == nil {
				continue
			}

			date, err := parseOcbcDate(getCol(record, headerMap, "Transaction date"))
			if err != nil {
				return nil, fmt.Errorf("error parsing transaction date: %w", err)
			}
			valueDate, err := parseOcbcDate(getCol(record, headerMap, "Value date"))
			if err != nil {
				return nil, fmt.Errorf("error parsing transaction posting date: %w", err)
			}
			withdrawalsSgd, err := parseNullDecimal(getCol(record, headerMap, "Withdrawals(SGD)"))
			if err != nil {
				return nil, fmt.Errorf("error parsing debit amount: %w", err)
			}
			depositsSgd, err := parseNullDecimal(getCol(record, headerMap, "Deposits(SGD)"))
			if err != nil {
				return nil, fmt.Errorf("error parsing credit amount: %w", err)
			}

			newTx := OcbcTransaction{
				BankAccountName: accountName,
				TransactionDate: date,
				ValueDate:       valueDate,
				Description:     getCol(record, headerMap, "Description"),
				WithdrawalsSGD:  withdrawalsSgd,
				DepositsSGD:     depositsSgd,
				RawRow:          strings.Join(append([]string{accountName}, record...), "|"),
			}
			transactions = append(transactions, newTx)

		}
	}

	return transactions, nil
}

func (p *OcbcStatementCsvParser) normalize(raw []OcbcTransaction) ([]commands.RawTransaction, error) {
	transactions := make([]commands.RawTransaction, len(raw))
	for idx, row := range raw {
		bankTx, err := NewBankTxFromOcbcTx(row)
		if err != nil {
			return nil, fmt.Errorf("failed to normalize ocbc transaction: %w", err)
		}

		transactions[idx] = bankTx
	}

	return transactions, nil
}

// NOTE: struct tags are for personal reference (not used in code)
type OcbcTransaction struct {
	BankAccountName string              `csv-h:"Account details for:"`
	TransactionDate time.Time           `csv:"Transaction date"`
	ValueDate       time.Time           `csv:"Value date"`
	Description     string              `csv:"Description"`
	WithdrawalsSGD  decimal.NullDecimal `csv:"Withdrawals(SGD)"`
	DepositsSGD     decimal.NullDecimal `csv:"Deposits(SGD)"`
	RawRow          string
}

func NewBankTxFromOcbcTx(row OcbcTransaction) (commands.RawTransaction, error) {
	var amount int64 = 0

	if row.WithdrawalsSGD.Valid && row.WithdrawalsSGD.Decimal.GreaterThan(decimal.Zero) {
		amount = row.WithdrawalsSGD.Decimal.Mul(decimal.NewFromInt(-1_000_000)).IntPart()
	}

	if row.DepositsSGD.Valid && row.DepositsSGD.Decimal.GreaterThan(decimal.Zero) {
		amount = row.DepositsSGD.Decimal.Mul(decimal.NewFromInt(1_000_000)).IntPart()
	}

	return commands.RawTransaction{
		ID:          row.RawRow,
		SourceName:  row.BankAccountName,
		Date:        row.TransactionDate,
		Description: row.Description,
		Amount:      amount,
	}, nil
}

func parseOcbcDate(s string) (time.Time, error) {
	t, err := time.Parse(OCBCDateLayout, strings.TrimSpace(s))
	if err == nil {
		return t, nil
	}

	t, err = time.Parse(OCBCFullYearDateLayout, strings.TrimSpace(s))
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse ocbc date (%s): %w", s, err)
	}

	return t, nil
}
