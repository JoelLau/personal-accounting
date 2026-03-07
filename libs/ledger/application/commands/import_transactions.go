package commands

import (
	"io"
	"libs/ledger/domain"
	"time"
)

type RawTransaction struct {
	ID          string    // WARN: strictly for de-duplication! DO NOT use for db id
	SourceName  string    // name of the bank account / credit card
	Date        time.Time // date of transaction
	Description string    // human-readable description
	Amount      int64     // raw units, can be positive or negative
}

type TransactionFileParser interface {
	Parse(file io.Reader) ([]RawTransaction, error)
}

type ImportProfile interface {
	Name() string
	NewPosting(raw RawTransaction) (domain.Posting, error)
}

type ImportTransactionsCommand struct {
	Reader  io.Reader
	Parser  TransactionFileParser
	Profile ImportProfile
}
