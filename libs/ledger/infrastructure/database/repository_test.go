package database_test

import (
	"database/sql"
	"embed"
	postgres "packages/accounting"
	"packages/accounting/dbgen"
	"packages/accounting/domain"
	testutil "packages/shared/test-util"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	testcontainers_postgres "github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestRepository_CreateTransactions(t *testing.T) {
	testCases := []struct {
		testName string
		given    []domain.BankTransaction
		postings []dbgen.Posting
		entries  []dbgen.Entry
	}{
		{
			testName: "Credit Card Debit",
			given: []domain.BankTransaction{
				{
					TransactionType: domain.TransactionSourceCreditCard,
					Date:            time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
					Description:     "CHICK-FIL-A SINGAPORE SG",
					Debit:           decimal.NewFromFloat(77.7),
					Credit:          decimal.Zero,
				},
			},
			postings: []dbgen.Posting{{
				ID:           1,
				Description:  pgtype.Text{String: "CHICK-FIL-A SINGAPORE SG", Valid: true},
				SystemNotes:  pgtype.Text{String: "", Valid: false},
				TransactedAt: time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
			}},
			entries: []dbgen.Entry{{
				ID:               1,
				Description:      pgtype.Text{String: "CHICK-FIL-A SINGAPORE SG", Valid: true},
				SystemNotes:      pgtype.Text{String: "", Valid: false},
				PostingsID:       1,
				LedgerAccountsID: 4000,
				DebitMicrosgd:    77_700_000,
				CreditMicrosgd:   0,
			}},
		},
		{
			testName: "Credit Card Credit (Refund)",
			given: []domain.BankTransaction{
				{
					TransactionType: domain.TransactionSourceCreditCard,
					Date:            time.Date(2025, 12, 22, 0, 0, 0, 0, time.UTC),
					Description:     "SHOPEE SINGAPORE MP SINGAPORE SG",
					Debit:           decimal.Zero,
					Credit:          decimal.NewFromFloat(36.65),
				},
			},
			postings: []dbgen.Posting{{
				ID:           1,
				Description:  pgtype.Text{String: "SHOPEE SINGAPORE MP SINGAPORE SG", Valid: true},
				SystemNotes:  pgtype.Text{String: "", Valid: false},
				TransactedAt: time.Date(2025, 12, 22, 0, 0, 0, 0, time.UTC),
			}},
			entries: []dbgen.Entry{{
				ID:               1,
				Description:      pgtype.Text{String: "SHOPEE SINGAPORE MP SINGAPORE SG", Valid: true},
				SystemNotes:      pgtype.Text{String: "", Valid: false},
				PostingsID:       1,
				LedgerAccountsID: 4000,
				DebitMicrosgd:    0,
				CreditMicrosgd:   36_650_000,
			}},
		},
		{
			testName: "Bank Statement Debit",
			given: []domain.BankTransaction{
				{
					TransactionType: domain.TransactionSourceBank,
					Date:            time.Date(2026, 2, 19, 0, 0, 0, 0, time.UTC),
					Description:     "OLD CHANG KEE Singapore SG",
					Debit:           decimal.NewFromFloat(1.9),
					Credit:          decimal.Zero,
				},
			},
			postings: []dbgen.Posting{{
				ID:           1,
				Description:  pgtype.Text{String: "OLD CHANG KEE Singapore SG", Valid: true},
				SystemNotes:  pgtype.Text{String: "", Valid: false},
				TransactedAt: time.Date(2026, 2, 19, 0, 0, 0, 0, time.UTC),
			}},
			entries: []dbgen.Entry{{
				ID:               1,
				Description:      pgtype.Text{String: "OLD CHANG KEE Singapore SG", Valid: true},
				SystemNotes:      pgtype.Text{String: "", Valid: false},
				PostingsID:       1,
				LedgerAccountsID: 4000,
				DebitMicrosgd:    1_900_000,
				CreditMicrosgd:   0,
			}},
		},
		{
			testName: "Bank Statement Credit (Salary)",
			given: []domain.BankTransaction{
				{
					TransactionType: domain.TransactionSourceBank,
					Date:            time.Date(2025, 12, 4, 0, 0, 0, 0, time.UTC),
					Description: `GIRO - SALARY
SALARY                            ALLEGIS GROUP SINGASALA`,
					Debit:  decimal.Zero,
					Credit: decimal.NewFromInt(8517),
				},
			},
			postings: []dbgen.Posting{{
				ID: 1,
				Description: pgtype.Text{String: `GIRO - SALARY
SALARY                            ALLEGIS GROUP SINGASALA`, Valid: true},
				SystemNotes:  pgtype.Text{String: "", Valid: false},
				TransactedAt: time.Date(2025, 12, 4, 0, 0, 0, 0, time.UTC),
			}},
			entries: []dbgen.Entry{{
				ID: 1,
				Description: pgtype.Text{String: `GIRO - SALARY
SALARY                            ALLEGIS GROUP SINGASALA`, Valid: true},
				SystemNotes:      pgtype.Text{String: "", Valid: false},
				PostingsID:       1,
				LedgerAccountsID: 3000,
				DebitMicrosgd:    0,
				CreditMicrosgd:   8_517_000_000,
			}},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			// Arrange
			t.Parallel()
			ctx := t.Context()

			postgresContainer := NewPostgresContainer(t)

			connStr, err := postgresContainer.ConnectionString(ctx)
			require.NoError(t, err)

			pool, err := pgxpool.New(ctx, connStr)
			require.NoError(t, err)

			repo := postgres.NewRepository(pool)

			// sanity checks
			q := dbgen.New(pool)
			postings, err := q.ListPostings(ctx)
			require.NoError(t, err)
			require.Empty(t, postings)
			entries, err := q.ListEntries(ctx)
			require.NoError(t, err)
			require.Empty(t, entries)

			// Act
			err = repo.CreateTransactions(ctx, tt.given)

			// Assert
			require.NoError(t, err)

			q = dbgen.New(pool)

			postings, err = q.ListPostings(ctx)
			assert.NoError(t, err)
			postingsDiff := cmp.Diff(postings, tt.postings, testutil.CmpOptions()...)
			assert.Emptyf(t, postingsDiff, postingsDiff)

			entries, err = q.ListEntries(ctx)
			entriesDiff := cmp.Diff(entries, tt.entries, testutil.CmpOptions()...)
			assert.NoError(t, err)
			assert.Emptyf(t, entriesDiff, entriesDiff)
		})
	}
}

//go:embed migrations/*.sql
var embedMigrations embed.FS

//go:embed seed/*.sql
var embedSeeds embed.FS

func NewPostgresContainer(t *testing.T) *testcontainers_postgres.PostgresContainer {
	t.Helper()

	ctx := t.Context()

	postgresContainer, err := testcontainers_postgres.Run(ctx, "postgres:18-alpine",
		testcontainers_postgres.WithDatabase(t.Name()),
		testcontainers_postgres.WithDatabase("username"),
		testcontainers_postgres.WithDatabase("p@ssw0rd!123"),
		testcontainers_postgres.BasicWaitStrategies(),
	)
	require.NoError(t, err)
	t.Cleanup(func() {
		if err := testcontainers.TerminateContainer(postgresContainer); err != nil {
			t.Logf("error terminating postgres container: %v", err)
		}
	})

	connStr, err := postgresContainer.ConnectionString(ctx)
	require.NoError(t, err)

	pool, err := pgxpool.New(ctx, connStr)
	require.NoError(t, err)

	goose.SetBaseFS(embedMigrations)
	err = goose.SetDialect("postgres")
	require.NoError(t, err)

	db := sql.OpenDB(stdlib.GetPoolConnector(pool))
	err = goose.Up(db, "migrations")
	require.NoError(t, err)

	goose.SetBaseFS(embedSeeds)
	err = goose.SetDialect("postgres")
	require.NoError(t, err)

	db = sql.OpenDB(stdlib.GetPoolConnector(pool))
	err = goose.Up(db, "seed", goose.WithNoVersioning())
	require.NoError(t, err)

	return postgresContainer
}
