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

	want := []commands.RawTransaction{
		{
			ID:          "J + K 609-123412-001|31/01/2026|31/01/2026|NETS QR\nAAMA BROTHER'S                    NETS QR PURCHASE   92041401|114.80|",
			SourceName:  "J + K 609-123412-001",
			Date:        time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC),
			Description: "NETS QR\nAAMA BROTHER'S                    NETS QR PURCHASE   92041401",
			Amount:      -114_800_000,
		},
		{
			ID:          "J + K 609-123412-001|31/01/2026|31/01/2026|NETS QR\nGLAM (85) PTE LTD                 NETS QR PURCHASE   23034119|2.30|",
			SourceName:  "J + K 609-123412-001",
			Date:        time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC),
			Description: "NETS QR\nGLAM (85) PTE LTD                 NETS QR PURCHASE   23034119",
			Amount:      -2_300_000,
		},
		{
			ID:          "J + K 609-123412-001|02/01/2026|02/01/2026|GIRO - SALARY\n00014                             EMPLOYER PTE LTD.  SALA||6,865.83",
			SourceName:  "J + K 609-123412-001",
			Date:        time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC),
			Description: "GIRO - SALARY\n00014                             EMPLOYER PTE LTD.  SALA",
			Amount:      6865_830_000,
		},
	}

	// Arrange
	tempFile := testutil.TemporaryFile(t, given)
	file, err := os.Open(tempFile.Name())
	require.NoError(t, err)

	parser := parsers.NewOcbcStatementCsvParser(2026, 01)

	// Act
	have, err := parser.Parse(file)

	// Assert
	assert.NoError(t, err)

	diff := cmp.Diff(have, want, testutil.CmpOptions()...)
	require.Emptyf(t, diff, diff)
}
