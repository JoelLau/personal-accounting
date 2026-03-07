package parsers_test

import (
	"libs/ledger/application/commands"
	"os"
	"packages/ingestion/parsers"
	testutil "packages/shared/test-util"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
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
"","","","","","","",""
"Supplementary Card:","","","","","","",""
"DBS Credit Card Type 4119-1100-6284-3393","","","","","","",""
"","","","","","","",""
"Transaction Date","Transaction Posting Date","Transaction Description","Transaction Type","Payment Type","Transaction Status","Debit Amount","Credit Amount"
"No transactions to view","","","","","","",""
"We didn't find any transactions for the selected date range.","","","","","","",""
"","","","","","","",""
"Blocked Card:","","","","","","",""
"DBS Credit Card Type 4419-4321-4321-4321","","","","","","",""
"","","","","","","",""
"Transaction Date","Transaction Posting Date","Transaction Description","Transaction Type","Payment Type","Transaction Status","Debit Amount","Credit Amount"
"22 Jan 2026","23 Jan 2026","KOREAN RESTAURANT SINGAPORE SG","PURCHASE","Contactless","Settled","16.35",""
"22 Jan 2026","23 Jan 2026","[Refund] KOREAN R SINGAPORE SG","REFUNDS & CREDITS","Others","Settled","","16.35"`

	want := []commands.RawTransaction{
		{
			ID:          "DBS Credit Card Type 1234-2345-3456-4567|23 Jan 2026|23 Jan 2026|SOME TELECOMMUNICATIONS SINGAPORE SGP|PURCHASE|Manual Key Entry|Settled|26.0|",
			SourceName:  "DBS Credit Card Type 1234-2345-3456-4567",
			Date:        time.Date(2026, 1, 23, 0, 0, 0, 0, time.UTC),
			Description: "SOME TELECOMMUNICATIONS SINGAPORE SGP",
			Amount:      26_000_000,
		},
		{
			ID:          "DBS Credit Card Type 4419-4321-4321-4321|22 Jan 2026|23 Jan 2026|KOREAN RESTAURANT SINGAPORE SG|PURCHASE|Contactless|Settled|16.35|",
			SourceName:  "DBS Credit Card Type 4419-4321-4321-4321",
			Date:        time.Date(2026, 1, 22, 0, 0, 0, 0, time.UTC),
			Description: "KOREAN RESTAURANT SINGAPORE SG",
			Amount:      16_350_000,
		},
		{
			ID:          "DBS Credit Card Type 4419-4321-4321-4321|22 Jan 2026|23 Jan 2026|[Refund] KOREAN R SINGAPORE SG|REFUNDS & CREDITS|Others|Settled||16.35",
			SourceName:  "DBS Credit Card Type 4419-4321-4321-4321",
			Date:        time.Date(2026, 1, 22, 0, 0, 0, 0, time.UTC),
			Description: "[Refund] KOREAN R SINGAPORE SG",
			Amount:      -16_350_000,
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
