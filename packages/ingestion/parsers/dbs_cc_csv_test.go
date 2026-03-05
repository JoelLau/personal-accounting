package parsers_test

import (
	"os"
	domain "packages/accounting/domain"
	"packages/ingestion/parsers"
	testutil "packages/shared/test-util"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDbsCreditCardCsvParser(t *testing.T) {
	t.Parallel()

	given := `"Card Transaction Details For:","DBS Credit Card Type 1234-2345-3456-4567",,,,,,
"Transactions as at:","19 Feb 2026",,,,,,
"",,,,,,,
"Credit Limit:","SGD 5678",,,,,,
"Available Limit:","SGD 6789",,,,,,
"",,,,,,,
"Transaction Date","Transaction Posting Date","Transaction Description","Transaction Type","Payment Type","Transaction Status","Debit Amount","Credit Amount"
"23 Jan 2026","23 Jan 2026","SOME TELECOMMUNICATIONS SINGAPORE SGP","PURCHASE","Manual Key Entry","Settled","26.0",""
"22 Jan 2026","23 Jan 2026","KOREAN RESTAURANT SINGAPORE SG","PURCHASE","Contactless","Settled","16.35",""
"","","","","","","",""
"Supplementary Card:","","","","","","",""
"DBS Vantage Visa Infinite Card 4119-1100-6284-3393","","","","","","",""
"","","","","","","",""
"Transaction Date","Transaction Posting Date","Transaction Description","Transaction Type","Payment Type","Transaction Status","Debit Amount","Credit Amount"
"No transactions to view","","","","","","",""
"We didn't find any transactions for the selected date range.","","","","","","",""
"","","","","","","",""
"Blocked Card:","","","","","","",""
"DBS Vantage Visa Infinite Card 4119-1100-6224-8858","","","","","","",""
"","","","","","","",""
"Transaction Date","Transaction Posting Date","Transaction Description","Transaction Type","Payment Type","Transaction Status","Debit Amount","Credit Amount"
"16 Jan 2026","19 Jan 2026","VETS FOR PETS (LENGKOK SINGAPORE SG","PURCHASE","Contactless","Settled","77.7",""`

	want := []domain.BankTransaction{
		{
			TransactionType: domain.TransactionSourceCreditCard,
			SourceName:      "DBS Credit Card Type 1234-2345-3456-4567",
			Date:            time.Date(2026, 1, 23, 0, 0, 0, 0, time.UTC),
			Description:     "SOME TELECOMMUNICATIONS SINGAPORE SGP",
			Debit:           decimal.NewFromFloat(26.0),
			Credit:          decimal.NewFromFloat(0.0),
			RawRow:          "23 Jan 2026|23 Jan 2026|SOME TELECOMMUNICATIONS SINGAPORE SGP|PURCHASE|Manual Key Entry|Settled|26.0|",
		},
		{
			TransactionType: domain.TransactionSourceCreditCard,
			SourceName:      "DBS Credit Card Type 1234-2345-3456-4567",
			Date:            time.Date(2026, 1, 22, 0, 0, 0, 0, time.UTC),
			Description:     "KOREAN RESTAURANT SINGAPORE SG",
			Debit:           decimal.NewFromFloat(16.35),
			Credit:          decimal.NewFromFloat(0.0),
			RawRow:          "22 Jan 2026|23 Jan 2026|KOREAN RESTAURANT SINGAPORE SG|PURCHASE|Contactless|Settled|16.35|",
		},
		{
			TransactionType: domain.TransactionSourceCreditCard,
			SourceName:      "DBS Vantage Visa Infinite Card 4119-1100-6224-8858",
			Date:            time.Date(2026, 1, 16, 0, 0, 0, 0, time.UTC),
			Description:     "VETS FOR PETS (LENGKOK SINGAPORE SG",
			Debit:           decimal.NewFromFloat(77.7),
			Credit:          decimal.NewFromFloat(0.0),
			RawRow:          "16 Jan 2026|19 Jan 2026|VETS FOR PETS (LENGKOK SINGAPORE SG|PURCHASE|Contactless|Settled|77.7|",
		},
	}

	// Arrange
	tempFile := testutil.TemporaryFile(t, given)
	file, err := os.Open(tempFile.Name())
	require.NoError(t, err)

	parser := parsers.NewDbsCreditCardCsvParser()

	// Act
	transactions, err := parser.Parse(file)

	// Assert
	assert.NoError(t, err)

	diff := cmp.Diff(transactions, want, testutil.CmpOptions()...)
	require.Emptyf(t, diff, diff)
}
