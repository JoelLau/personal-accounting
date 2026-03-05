package domain

import (
	"errors"
	"time"

	"github.com/shopspring/decimal"
)

var (
	ErrEntryIsBothDebitCredit    = errors.New("entry is both debit and credit")
	ErrEntryIsNeitherDebitCredit = errors.New("entry is neither debit or credit")
	ErrInvalidTransactionType    = errors.New("invalid transaction type")
)

type TransactionSource string

const (
	TransactionSourceBank       TransactionSource = "BANK"
	TransactionSourceCreditCard TransactionSource = "CREDIT_CARD"
)

func (t TransactionSource) IsValid() bool {
	switch t {
	case TransactionSourceBank, TransactionSourceCreditCard:
		return true
	default:
		return false
	}
}

type BankTransaction struct {
	TransactionType TransactionSource
	SourceName      string // e.g. credit card number / bank account number
	Date            time.Time
	Description     string
	Debit           decimal.Decimal
	Credit          decimal.Decimal
	RawRow          string
}

func NewBankTransaction(
	transactionType TransactionSource,
	sourceName string,
	date time.Time,
	description string,
	debit decimal.Decimal,
	credit decimal.Decimal,
	rawRow string,
) (BankTransaction, error) {
	e := BankTransaction{
		TransactionType: transactionType,
		SourceName:      sourceName,
		Date:            date,
		Description:     description,
		Debit:           debit,
		Credit:          credit,
		RawRow:          rawRow,
	}
	return e, e.Error()
}

func (e BankTransaction) IsDebit() bool {
	return e.Debit.GreaterThan(decimal.Zero)
}

func (e BankTransaction) IsCredit() bool {
	return e.Credit.GreaterThan(decimal.Zero)
}

func (e BankTransaction) IsValid() bool {
	return e.Error() == nil
}

// NOTE: consider combining errors (no use differentiating them)
func (e BankTransaction) Error() error {
	if e.IsDebit() && e.IsCredit() {
		return ErrEntryIsBothDebitCredit
	}

	if !e.IsDebit() && !e.IsCredit() {
		return ErrEntryIsNeitherDebitCredit
	}

	if !e.TransactionType.IsValid() {
		return ErrInvalidTransactionType
	}

	return nil
}
