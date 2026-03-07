package application_test

import (
	"libs/ledger/application"
	"libs/ledger/application/commands"
	"libs/ledger/domain"
	testutil "packages/shared/test-util"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func TestDBSImportProfile(t *testing.T) {
	const (
		ExpensesAccountID          = 4000
		IncomeAccountID            = 3000
		LiabilitiesAccountID int64 = 2001
	)

	testCases := []struct {
		name string

		given commands.RawTransaction
		want  domain.Posting
	}{
		{
			name: "record money in as refund",
			given: commands.RawTransaction{
				ID:          "DBS Card Type 4119-1234-1234-1234|22 Jan 2026|23 Jan 2026|[REFUND] THE PALACE KO SINGAPORE SG|REFUND & CREDITS|Others|Settled||16.35",
				SourceName:  "DBS Card Type 4119-1234-1234-1234",
				Date:        time.Date(2026, 1, 23, 0, 0, 0, 0, time.UTC),
				Description: "[REFUND] THE PALACE KO SINGAPORE SG",
				Amount:      -16_350_000,
			},
			want: domain.Posting{
				ID:          "DBS Card Type 4119-1234-1234-1234|22 Jan 2026|23 Jan 2026|[REFUND] THE PALACE KO SINGAPORE SG|REFUND & CREDITS|Others|Settled||16.35",
				Date:        time.Date(2026, 1, 23, 0, 0, 0, 0, time.UTC),
				Description: "[REFUND] THE PALACE KO SINGAPORE SG",
				Entries: []domain.Entry{
					{AccountID: LiabilitiesAccountID, Description: "[REFUND] THE PALACE KO SINGAPORE SG", DebitMicroSGD: 16_350_000, CreditMicroSGD: 0},
					{AccountID: ExpensesAccountID, Description: "[REFUND] THE PALACE KO SINGAPORE SG", DebitMicroSGD: 0, CreditMicroSGD: 16_350_000},
				},
			},
		},
		{
			name: "record money out as expense",
			given: commands.RawTransaction{
				ID:          "DBS Card Type 4119-1234-1234-1234|23 Jan 2026|23 Jan 2026|WHIZ COMMUNICATIONS SINGAPORE SGP|PURCHASE|Manual Key Entry|Settled|26.0|",
				SourceName:  "DBS Card Type 4119-1234-1234-1234",
				Date:        time.Date(2026, 1, 23, 0, 0, 0, 0, time.UTC),
				Description: "WHIZ COMMUNICATIONS SINGAPORE SGP",
				Amount:      26_000_000,
			},
			want: domain.Posting{
				ID:          "DBS Card Type 4119-1234-1234-1234|23 Jan 2026|23 Jan 2026|WHIZ COMMUNICATIONS SINGAPORE SGP|PURCHASE|Manual Key Entry|Settled|26.0|",
				Date:        time.Date(2026, 1, 23, 0, 0, 0, 0, time.UTC),
				Description: "WHIZ COMMUNICATIONS SINGAPORE SGP",
				Entries: []domain.Entry{
					{AccountID: LiabilitiesAccountID, Description: "WHIZ COMMUNICATIONS SINGAPORE SGP", DebitMicroSGD: 0, CreditMicroSGD: 26_000_000},
					{AccountID: ExpensesAccountID, Description: "WHIZ COMMUNICATIONS SINGAPORE SGP", DebitMicroSGD: 26_000_000, CreditMicroSGD: 0},
				},
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			t.Parallel()
			profile := application.NewDBSImportProfile(ExpensesAccountID, LiabilitiesAccountID)

			// Act
			have, err := profile.NewPosting(tt.given)

			// Assert
			require.NoError(t, err)
			require.Empty(t, cmp.Diff(have, tt.want, testutil.CmpOptions()...))
		})
	}
}

func TestOCBCImportProfile(t *testing.T) {
	const (
		AssetsAccountID         = 1    // "Assets"
		ExpensesAccountID int64 = 4000 // "Expenses:Uncategorized"
		IncomeAccountID         = 3000 // "Income:Uncategorized"
	)

	testCases := []struct {
		name string

		given commands.RawTransaction
		want  domain.Posting
	}{
		{
			name: "record money in as income (debit 'assets', credit 'income')",
			given: commands.RawTransaction{
				ID:          "J + K 609-123412-001|02/01/2026|02/01/2026|GIRO - SALARY\n00014                             EMPLOYER PTE LTD.  SALA||6,865.83",
				SourceName:  "J + K 609-123412-001",
				Date:        time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC),
				Description: "GIRO - SALARY\n00014                             EMPLOYER PTE LTD.  SALA",
				Amount:      6865_830_000,
			},
			want: domain.Posting{
				ID:          "J + K 609-123412-001|02/01/2026|02/01/2026|GIRO - SALARY\n00014                             EMPLOYER PTE LTD.  SALA||6,865.83",
				Date:        time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC),
				Description: "GIRO - SALARY\n00014                             EMPLOYER PTE LTD.  SALA",
				Entries: []domain.Entry{
					{AccountID: IncomeAccountID, Description: "GIRO - SALARY\n00014                             EMPLOYER PTE LTD.  SALA", DebitMicroSGD: 0, CreditMicroSGD: 6865_830_000},
					{AccountID: AssetsAccountID, Description: "GIRO - SALARY\n00014                             EMPLOYER PTE LTD.  SALA", DebitMicroSGD: 6865_830_000, CreditMicroSGD: 0},
				},
			},
		},
		{
			name: "record money out as expense paid out of assets",
			given: commands.RawTransaction{
				ID:          "J + K 609-123412-001|31/01/2026|31/01/2026|NETS QR\nAAMA BROTHER'S                    NETS QR PURCHASE   92041401|114.80|",
				SourceName:  "J + K 609-123412-001",
				Date:        time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC),
				Description: "NETS QR\nAAMA BROTHER'S                    NETS QR PURCHASE   92041401",
				Amount:      -114_800_000,
			},
			want: domain.Posting{
				ID:          "J + K 609-123412-001|31/01/2026|31/01/2026|NETS QR\nAAMA BROTHER'S                    NETS QR PURCHASE   92041401|114.80|",
				Date:        time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC),
				Description: "NETS QR\nAAMA BROTHER'S                    NETS QR PURCHASE   92041401",
				Entries: []domain.Entry{
					{AccountID: ExpensesAccountID, Description: "NETS QR\nAAMA BROTHER'S                    NETS QR PURCHASE   92041401", DebitMicroSGD: 114_800_000, CreditMicroSGD: 0},
					{AccountID: AssetsAccountID, Description: "NETS QR\nAAMA BROTHER'S                    NETS QR PURCHASE   92041401", DebitMicroSGD: 0, CreditMicroSGD: 114_800_000},
				},
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			t.Parallel()
			profile := application.NewOCBCStatementProfile(AssetsAccountID, ExpensesAccountID, IncomeAccountID)

			// Act
			have, err := profile.NewPosting(tt.given)

			// Assert
			require.NoError(t, err)
			require.Empty(t, cmp.Diff(have, tt.want, testutil.CmpOptions()...))
		})
	}
}
