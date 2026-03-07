package application_test

import (
	"libs/ledger/application"
	"libs/ledger/application/commands"
	"libs/ledger/application/handlers"
	"libs/ledger/application/services"
	"libs/ledger/infrastructure/database"
	dbgen "libs/ledger/infrastructure/database/gen"
	"packages/ingestion/parsers"
	testutil "packages/shared/test-util"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var cmpOptPostings = []cmp.Option{
	cmpopts.SortSlices(func(a, b dbgen.Posting) int { return int(a.ID - b.ID) }),
	cmpopts.EquateApproxTime(0),
}

var cmpOptEntries = []cmp.Option{
	cmpopts.SortSlices(func(a, b dbgen.Entry) int { return int(a.ID - b.ID) }),
}

func TestImportTransactions(t *testing.T) {
	const (
		ExpensesAccountID          = 4000
		IncomeAccountID            = 3000
		LiabilitiesAccountID int64 = 2001
	)

	testCases := []struct {
		testName string
		given    string
		postings []dbgen.Posting
		entries  []dbgen.Entry
	}{
		{
			testName: "dbs credit card",
			given: `"Card Transaction Details For:","DBS Card Type 4119-1234-1234-1234",,,,,,
"Transactions as at:","19 Feb 2026",,,,,,
"",,,,,,,
"Credit Limit:","SGD 12345",,,,,,
"Available Limit:","SGD 12345",,,,,,
"",,,,,,,
"Transaction Date","Transaction Posting Date","Transaction Description","Transaction Type","Payment Type","Transaction Status","Debit Amount","Credit Amount"
"23 Jan 2026","23 Jan 2026","WHIZ COMMUNICATIONS SINGAPORE SGP","PURCHASE","Manual Key Entry","Settled","26.0",""
"22 Jan 2026","23 Jan 2026","THE PALACE KOREAN REST SINGAPORE SG","PURCHASE","Contactless","Settled","16.35",""
"22 Jan 2026","23 Jan 2026","[REFUND] THE PALACE KO SINGAPORE SG","REFUND & CREDITS","Others","Settled","","16.35"`,
			postings: []dbgen.Posting{
				{
					ID:           1,
					Description:  pgtype.Text{String: "WHIZ COMMUNICATIONS SINGAPORE SGP", Valid: true},
					TransactedAt: time.Date(2026, 01, 23, 0, 0, 0, 0, time.UTC),
					SourceHash:   pgtype.Text{String: "DBS Card Type 4119-1234-1234-1234|23 Jan 2026|23 Jan 2026|WHIZ COMMUNICATIONS SINGAPORE SGP|PURCHASE|Manual Key Entry|Settled|26.0|", Valid: true},
				},
				{
					ID:           2,
					Description:  pgtype.Text{String: "THE PALACE KOREAN REST SINGAPORE SG", Valid: true},
					TransactedAt: time.Date(2026, 01, 22, 0, 0, 0, 0, time.UTC),
					SourceHash:   pgtype.Text{String: "DBS Card Type 4119-1234-1234-1234|22 Jan 2026|23 Jan 2026|THE PALACE KOREAN REST SINGAPORE SG|PURCHASE|Contactless|Settled|16.35|", Valid: true},
				},
				{
					ID:           3,
					Description:  pgtype.Text{String: "[REFUND] THE PALACE KO SINGAPORE SG", Valid: true},
					TransactedAt: time.Date(2026, 01, 22, 0, 0, 0, 0, time.UTC),
					SourceHash:   pgtype.Text{String: "DBS Card Type 4119-1234-1234-1234|22 Jan 2026|23 Jan 2026|[REFUND] THE PALACE KO SINGAPORE SG|REFUND & CREDITS|Others|Settled||16.35", Valid: true},
				},
			},
			entries: []dbgen.Entry{
				// Transaction 1: Spend $26
				{ID: 1, Description: pgtype.Text{String: "WHIZ COMMUNICATIONS SINGAPORE SGP", Valid: true}, SystemNotes: pgtype.Text{String: "", Valid: false}, PostingsID: 1, LedgerAccountsID: LiabilitiesAccountID, DebitMicrosgd: 0, CreditMicrosgd: 26_000_000},
				{ID: 2, Description: pgtype.Text{String: "WHIZ COMMUNICATIONS SINGAPORE SGP", Valid: true}, SystemNotes: pgtype.Text{String: "", Valid: false}, PostingsID: 1, LedgerAccountsID: ExpensesAccountID, DebitMicrosgd: 26_000_000, CreditMicrosgd: 0},

				// Transaction 2: Spend $16.35
				{ID: 3, Description: pgtype.Text{String: "THE PALACE KOREAN REST SINGAPORE SG", Valid: true}, SystemNotes: pgtype.Text{String: "", Valid: false}, PostingsID: 2, LedgerAccountsID: LiabilitiesAccountID, DebitMicrosgd: 0, CreditMicrosgd: 16_350_000},
				{ID: 4, Description: pgtype.Text{String: "THE PALACE KOREAN REST SINGAPORE SG", Valid: true}, SystemNotes: pgtype.Text{String: "", Valid: false}, PostingsID: 2, LedgerAccountsID: ExpensesAccountID, DebitMicrosgd: 16_350_000, CreditMicrosgd: 0},

				// Transaction 3: Refund $16.35
				{ID: 5, Description: pgtype.Text{String: "[REFUND] THE PALACE KO SINGAPORE SG", Valid: true}, SystemNotes: pgtype.Text{String: "", Valid: false}, PostingsID: 3, LedgerAccountsID: LiabilitiesAccountID, DebitMicrosgd: 16_350_000, CreditMicrosgd: 0},
				{ID: 6, Description: pgtype.Text{String: "[REFUND] THE PALACE KO SINGAPORE SG", Valid: true}, SystemNotes: pgtype.Text{String: "", Valid: false}, PostingsID: 3, LedgerAccountsID: ExpensesAccountID, DebitMicrosgd: 0, CreditMicrosgd: 16_350_000},
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			t.Parallel()

			// Arrange
			ctx := t.Context()

			postgresContainer := testutil.NewPostgresContainer(t, database.EmbedMigrations, database.EmbedSeed)

			dsn, err := postgresContainer.ConnectionString(ctx)
			require.NoError(t, err)

			pool, err := pgxpool.New(ctx, dsn)
			require.NoError(t, err)
			defer pool.Close()

			reader := strings.NewReader(tt.given)

			cmd := commands.ImportTransactionsCommand{
				Reader:  reader,
				Parser:  parsers.NewDbsCreditCardCsvParser(),
				Profile: application.NewDBSImportProfile(ExpensesAccountID, LiabilitiesAccountID),
			}

			repo := database.NewRepository(pool)
			service := services.NewImportTransactionsService(repo)
			handler := handlers.NewImportTransactionsHandler(service)

			// Act
			err = handler.Handle(ctx, cmd)
			require.NoError(t, err)

			// Assert
			q := dbgen.New(pool)

			postings, err := q.ListPostings(ctx)
			require.NoError(t, err)
			assert.Empty(t, cmp.Diff(postings, tt.postings, cmpOptPostings...))

			entries, err := q.ListEntries(ctx)
			require.NoError(t, err)
			assert.Empty(t, cmp.Diff(entries, tt.entries, cmpOptEntries...))
		})
	}
}
