package usecase_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	domain "packages/accounting/domain"
	usecase "packages/ingestion/use-case"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIngestService_ParserErrors(t *testing.T) {
	testCases := []struct {
		testName           string
		givenParsers       []usecase.TransactionParser
		expectTransactions []domain.BankTransaction
		expectErr          bool
	}{
		{
			testName:           "given no parsers, return error",
			givenParsers:       make([]usecase.TransactionParser, 0),
			expectTransactions: nil,
			expectErr:          true,
		},
		{
			testName: "given 1 good parser sandwiched between 2 error parsers, return success",
			givenParsers: []usecase.TransactionParser{
				&MockTransactionParser{ParseFn: func(file io.Reader) ([]domain.BankTransaction, error) { return nil, errors.New("error A") }},
				&MockTransactionParser{ParseFn: func(file io.Reader) ([]domain.BankTransaction, error) {
					return []domain.BankTransaction{{}}, nil
				}},
				&MockTransactionParser{ParseFn: func(file io.Reader) ([]domain.BankTransaction, error) { return nil, errors.New("error B") }},
			},
			expectTransactions: []domain.BankTransaction{{}},
			expectErr:          false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			t.Parallel()

			// Arrange
			ctx := t.Context()
			repo := &MockTransactionRepository{}
			service := usecase.NewIngestionService(repo, tt.givenParsers)

			// Act
			err := service.IngestFile(ctx, nil)

			// Assert
			assert.Equalf(t, err != nil, tt.expectErr, fmt.Sprintf("unexpected error: %v", err))
			assert.Equal(t, tt.expectTransactions, repo.Transactions)
		})
	}
}

type MockTransactionRepository struct {
	Transactions []domain.BankTransaction
	Error        error
}

var _ usecase.TransactionRepository = &MockTransactionRepository{}

func (m *MockTransactionRepository) CreateTransactions(ctx context.Context, transactions []domain.BankTransaction) error {
	m.Transactions = append(m.Transactions, transactions...)

	return m.Error
}

type MockTransactionParser struct {
	ParseFn func(file io.Reader) ([]domain.BankTransaction, error)
}

var _ usecase.TransactionParser = &MockTransactionParser{}

func (m *MockTransactionParser) Parse(file io.Reader) ([]domain.BankTransaction, error) {
	return m.ParseFn(file)
}
