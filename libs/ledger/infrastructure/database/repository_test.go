package database_test

import (
	"libs/ledger/domain"
	"libs/ledger/infrastructure/database"
	dbgen "libs/ledger/infrastructure/database/gen"
	testutil "packages/shared/test-util"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepository_CreateTransactions(t *testing.T) {
	testCases := []struct {
		testName string
		given    []domain.Posting
		postings []dbgen.Posting
		entries  []dbgen.Entry
	}{
		{
			testName: "DBS Credit Card - Debit",
			given: []domain.Posting{
				{
					ID:          "DBS Card Type 4119-1234-1234-1234|23 Jan 2026|23 Jan 2026|WHIZ COMMUNICATIONS SINGAPORE SGP|PURCHASE|Manual Key Entry|Settled|26.0|",
					Date:        time.Date(2026, 01, 23, 0, 0, 0, 0, time.UTC),
					Description: "WHIZ COMMUNICATIONS SINGAPORE SGP",
					Entries: []domain.Entry{
						{AccountID: 2, Description: "WHIZ COMMUNICATIONS SINGAPORE SGP", DebitMicroSGD: 0, CreditMicroSGD: 26_000_000},
						{AccountID: 4000, Description: "WHIZ COMMUNICATIONS SINGAPORE SGP", DebitMicroSGD: 26_000_000, CreditMicroSGD: 0},
					},
				},
			},
			postings: []dbgen.Posting{
				{
					ID:           1,
					Description:  pgtype.Text{String: "WHIZ COMMUNICATIONS SINGAPORE SGP", Valid: true},
					SystemNotes:  pgtype.Text{String: "", Valid: false},
					TransactedAt: time.Date(2026, 01, 23, 0, 0, 0, 0, time.UTC),
					SourceHash:   pgtype.Text{String: "DBS Card Type 4119-1234-1234-1234|23 Jan 2026|23 Jan 2026|WHIZ COMMUNICATIONS SINGAPORE SGP|PURCHASE|Manual Key Entry|Settled|26.0|", Valid: true},
				},
			},
			entries: []dbgen.Entry{
				{
					ID:               1,
					Description:      pgtype.Text{String: "WHIZ COMMUNICATIONS SINGAPORE SGP", Valid: true},
					SystemNotes:      pgtype.Text{String: "", Valid: false},
					PostingsID:       1,
					LedgerAccountsID: 2,
					DebitMicrosgd:    0,
					CreditMicrosgd:   26_000_000,
				},
				{
					ID:               2,
					Description:      pgtype.Text{String: "WHIZ COMMUNICATIONS SINGAPORE SGP", Valid: true},
					SystemNotes:      pgtype.Text{String: "", Valid: false},
					PostingsID:       1,
					LedgerAccountsID: 4000,
					DebitMicrosgd:    26_000_000,
					CreditMicrosgd:   0,
				},
			},
		},
		{
			testName: "DBS Credit Card - Credit",
			given: []domain.Posting{
				{
					ID:          "DBS Card Type 4119-1234-1234-1234|22 Jan 2026|23 Jan 2026|[REFUND] THE PALACE KO SINGAPORE SG|REFUND AND CREDITS|Others|Settled||26.0",
					Date:        time.Date(2026, 01, 22, 0, 0, 0, 0, time.UTC),
					Description: "[REFUND] THE PALACE KO SINGAPORE SG",
					Entries: []domain.Entry{
						{AccountID: 4000, Description: "[REFUND] THE PALACE KO SINGAPORE SG", DebitMicroSGD: 0, CreditMicroSGD: 16_350_000},
						{AccountID: 2, Description: "[REFUND] THE PALACE KO SINGAPORE SG", DebitMicroSGD: 16_350_000, CreditMicroSGD: 0},
					},
				},
			},
			postings: []dbgen.Posting{
				{
					ID:           1,
					Description:  pgtype.Text{String: "[REFUND] THE PALACE KO SINGAPORE SG", Valid: true},
					SystemNotes:  pgtype.Text{String: "", Valid: false},
					TransactedAt: time.Date(2026, 01, 22, 0, 0, 0, 0, time.UTC),
					SourceHash:   pgtype.Text{String: "DBS Card Type 4119-1234-1234-1234|22 Jan 2026|23 Jan 2026|[REFUND] THE PALACE KO SINGAPORE SG|REFUND AND CREDITS|Others|Settled||26.0", Valid: true},
				},
			},
			entries: []dbgen.Entry{
				{
					ID:               1,
					Description:      pgtype.Text{String: "[REFUND] THE PALACE KO SINGAPORE SG", Valid: true},
					SystemNotes:      pgtype.Text{String: "", Valid: false},
					PostingsID:       1,
					LedgerAccountsID: 4000,
					DebitMicrosgd:    0,
					CreditMicrosgd:   16_350_000,
				},
				{
					ID:               2,
					Description:      pgtype.Text{String: "[REFUND] THE PALACE KO SINGAPORE SG", Valid: true},
					SystemNotes:      pgtype.Text{String: "", Valid: false},
					PostingsID:       1,
					LedgerAccountsID: 2,
					DebitMicrosgd:    16_350_000,
					CreditMicrosgd:   0,
				},
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			// Arrange
			t.Parallel()
			ctx := t.Context()

			postgresContainer := testutil.NewPostgresContainer(t, database.EmbedMigrations, database.EmbedSeed)

			connStr, err := postgresContainer.ConnectionString(ctx)
			require.NoError(t, err)

			pool, err := pgxpool.New(ctx, connStr)
			require.NoError(t, err)

			repo := database.NewRepository(pool)

			// Act
			err = repo.CreatePostings(ctx, tt.given)

			// Assert
			require.NoError(t, err)

			q := dbgen.New(pool)

			postings, err := q.ListPostings(ctx)
			assert.NoError(t, err)
			postingsDiff := cmp.Diff(postings, tt.postings, testutil.CmpOptions()...)
			assert.Emptyf(t, postingsDiff, postingsDiff)

			entries, err := q.ListEntries(ctx)
			entriesDiff := cmp.Diff(entries, tt.entries, testutil.CmpOptions()...)
			assert.NoError(t, err)
			assert.Emptyf(t, entriesDiff, entriesDiff)
		})
	}
}
