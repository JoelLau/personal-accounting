package domain_test

import (
	domain "packages/accounting/domain"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBankTransaction(t *testing.T) {
	testCases := []struct {
		testName string

		transactionType domain.TransactionSource
		date            time.Time
		description     string
		debit           decimal.Decimal
		credit          decimal.Decimal

		err error
	}{
		{
			testName:        "given 0 credit and debit, return error",
			transactionType: domain.TransactionSourceBank,
			date:            time.Date(2026, 02, 28, 0, 0, 0, 0, time.UTC),
			description:     "some description",
			debit:           decimal.Zero,
			credit:          decimal.Zero,
			err:             domain.ErrEntryIsNeitherDebitCredit,
		},
		{
			testName:        "given debit, return success",
			transactionType: domain.TransactionSourceBank,
			date:            time.Date(2026, 02, 28, 0, 0, 0, 0, time.UTC),
			description:     "some description",
			debit:           decimal.NewFromInt(100),
			credit:          decimal.Zero,
			err:             nil,
		},
		{
			testName:        "given credit, return success",
			transactionType: domain.TransactionSourceCreditCard,
			date:            time.Date(2026, 02, 28, 0, 0, 0, 0, time.UTC),
			description:     "some description",
			debit:           decimal.Zero,
			credit:          decimal.NewFromFloat(12.12),
			err:             nil,
		},
		{
			testName:        "given both credit and debit, return error",
			transactionType: domain.TransactionSourceCreditCard,
			date:            time.Date(2026, 02, 28, 0, 0, 0, 0, time.UTC),
			description:     "some description",
			debit:           decimal.NewFromFloat(13.13),
			credit:          decimal.NewFromInt(12),
			err:             domain.ErrEntryIsBothDebitCredit,
		},
		{
			testName:        "given invalid transaction type",
			transactionType: "invalid transaction type",
			date:            time.Date(2026, 02, 28, 0, 0, 0, 0, time.UTC),
			description:     "some description",
			debit:           decimal.Zero,
			credit:          decimal.NewFromInt(12),
			err:             domain.ErrInvalidTransactionType,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			t.Parallel()

			bankTransaction, err := domain.NewBankTransaction(tt.transactionType, "", tt.date, tt.description, tt.debit, tt.credit, "")

			assert.Equal(t, tt.transactionType, bankTransaction.TransactionType)
			assert.Truef(t, tt.date.Equal(bankTransaction.Date), "have %v, want %v", bankTransaction, tt.date)
			assert.Equal(t, tt.description, bankTransaction.Description)
			assert.Truef(t, tt.debit.Equal(bankTransaction.Debit), "have %v, want %v", bankTransaction.Debit, tt.debit)
			assert.Truef(t, tt.credit.Equal(bankTransaction.Credit), "have %v, want %v", bankTransaction.Credit, tt.credit)

			require.ErrorIs(t, err, tt.err)
		})
	}
}
