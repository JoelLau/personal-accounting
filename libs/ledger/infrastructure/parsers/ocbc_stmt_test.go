package parsers_test

import (
	"os"
	"packages/accounting/domain"
	"packages/ingestion/parsers"
	testutil "packages/shared/test-util"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOcbcStatementCsvParser(t *testing.T) {
	t.Parallel()

	given := `Account details for:,J + K 609-123412-001
Available Balance,"12,123.12"
Ledger Balance,"12,123.12"

Transaction History
Transaction date,Value date,Description,Withdrawals(SGD),Deposits(SGD)
31/01/2026,31/01/2026,"NETS QR
AAMA BROTHER'S                    NETS QR PURCHASE   92041401",114.80,
31/01/2026,31/01/2026,"NETS QR
GLAM (85) PTE LTD                 NETS QR PURCHASE   23034119",2.30,
02/01/2026,02/01/2026,"GIRO - SALARY
00014                             EMPLOYER PTE LTD.  SALA",,"6,865.83"
`

	want := []domain.BankTransaction{
		{
			TransactionType: domain.TransactionSourceBank,
			SourceName:      "J + K 609-123412-001",
			Date:            time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC),
			Description:     "NETS QR\nAAMA BROTHER'S                    NETS QR PURCHASE   92041401",
			Debit:           decimal.NewFromFloat(114.8),
			Credit:          decimal.Zero,
			RawRow:          "31/01/2026|31/01/2026|NETS QR\nAAMA BROTHER'S                    NETS QR PURCHASE   92041401|114.80|",
		},
		{
			TransactionType: domain.TransactionSourceBank,
			SourceName:      "J + K 609-123412-001",
			Date:            time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC),
			Description:     "NETS QR\nGLAM (85) PTE LTD                 NETS QR PURCHASE   23034119",
			Debit:           decimal.NewFromFloat(2.3),
			Credit:          decimal.Zero,
			RawRow:          "31/01/2026|31/01/2026|NETS QR\nGLAM (85) PTE LTD                 NETS QR PURCHASE   23034119|2.30|",
		},
		{
			TransactionType: domain.TransactionSourceBank,
			SourceName:      "J + K 609-123412-001",
			Date:            time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC),
			Description:     "GIRO - SALARY\n00014                             EMPLOYER PTE LTD.  SALA",
			Debit:           decimal.Zero,
			Credit:          decimal.NewFromFloat(6865.83),
			RawRow:          "02/01/2026|02/01/2026|GIRO - SALARY\n00014                             EMPLOYER PTE LTD.  SALA||6,865.83",
		},
	}

	// Arrange
	tempFile := testutil.TemporaryFile(t, given)
	file, err := os.Open(tempFile.Name())
	require.NoError(t, err)

	parser := parsers.NewOcbcStatementCsvParser()

	// Act
	have, err := parser.Parse(file)

	// Assert
	assert.NoError(t, err)

	diff := cmp.Diff(have, want, testutil.CmpOptions()...)
	require.Emptyf(t, diff, diff)
}
