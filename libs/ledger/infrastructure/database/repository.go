package database

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"libs/ledger/domain"
	dbgen "libs/ledger/infrastructure/database/gen"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// https://www.postgresql.org/docs/current/errcodes-appendix.html
const PgErrCodeUniqueViolation = "23505"

var ErrDuplicateTransaction = errors.New("transaction already exists")

//go:embed migrations/*.sql
var EmbedMigrations embed.FS

//go:embed seed/*.sql
var EmbedSeed embed.FS

type Repository struct {
	db   *dbgen.Queries
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		db:   dbgen.New(pool),
		pool: pool,
	}
}

func (repo *Repository) CreatePostings(ctx context.Context, rawPostings []domain.Posting) error {
	tx, err := repo.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to create postgres transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	qtx := repo.db.WithTx(tx)

	for _, rawPosting := range rawPostings {
		postingParams := dbgen.CreatePostingParams{
			Description:  pgtype.Text{String: rawPosting.Description, Valid: true},
			SystemNotes:  pgtype.Text{String: "", Valid: false},
			TransactedAt: rawPosting.Date,
			SourceHash:   pgtype.Text{String: rawPosting.ID, Valid: true},
		}
		posting, err := qtx.CreatePosting(ctx, postingParams)
		if err != nil {
			return fmt.Errorf("failed to create posting %v: %w", postingParams, err)
		}

		for _, rawEntry := range rawPosting.Entries {
			entryParams := dbgen.CreateEntryParams{
				Description:      posting.Description,
				SystemNotes:      pgtype.Text{String: "", Valid: false},
				PostingsID:       posting.ID,
				LedgerAccountsID: rawEntry.AccountID,
				DebitMicrosgd:    rawEntry.DebitMicroSGD,
				CreditMicrosgd:   rawEntry.CreditMicroSGD,
			}

			_, err := qtx.CreateEntry(ctx, entryParams)
			if err != nil {
				return fmt.Errorf("failed to create entry (%v): %w", entryParams, err)
			}
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit postgres transaction: %w", err)
	}

	return nil
}
