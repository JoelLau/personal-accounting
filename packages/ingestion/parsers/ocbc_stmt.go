package parsers

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"packages/accounting/domain"
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

func (p *OcbcStatementCsvParser) Parse(file io.Reader) ([]domain.BankTransaction, error) {
	transactions, err := ParseOcbcStatementCsv(file)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ocbc transactions: %w", err)
	}

	bankTransactions := make([]domain.BankTransaction, len(transactions))
	for idx, tx := range transactions {
		bankTx, err := NewBankTxFromOcbcTx(tx)
		if err != nil {
			return nil, fmt.Errorf("failed to normalize ocbc transaction: %w", err)
		}

		bankTransactions[idx] = bankTx
	}

	return bankTransactions, nil
}

func ParseOcbcStatementCsv(r io.Reader) ([]OcbcTransaction, error) {
	csvReader := csv.NewReader(r)
	csvReader.FieldsPerRecord = -1

	var (
		transactions     []OcbcTransaction
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
		case strings.Contains(firstCol, "Account details for:"):
			if len(record) > 1 && record[1] != "" {
				currCard = record[1]
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
				BankAccountName: currCard,
				TransactionDate: date,
				ValueDate:       valueDate,
				Description:     getCol(record, headerMap, "Description"),
				WithdrawalsSGD:  withdrawalsSgd,
				DepositsSGD:     depositsSgd,
				RawRow:          strings.Join(record, "|"),
			}
			transactions = append(transactions, newTx)

		}
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

func NewBankTxFromOcbcTx(ocbcTx OcbcTransaction) (domain.BankTransaction, error) {
	debitAmount := decimal.Zero
	if ocbcTx.WithdrawalsSGD.Valid {
		debitAmount = ocbcTx.WithdrawalsSGD.Decimal
	}

	creditAmount := decimal.Zero
	if ocbcTx.DepositsSGD.Valid {
		creditAmount = ocbcTx.DepositsSGD.Decimal
	}

	return domain.BankTransaction{
		TransactionType: domain.TransactionSourceBank,
		SourceName:      ocbcTx.BankAccountName,
		Date:            ocbcTx.TransactionDate,
		Description:     ocbcTx.Description,
		Debit:           debitAmount,
		Credit:          creditAmount,
		RawRow:          ocbcTx.RawRow,
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
