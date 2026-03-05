package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"packages/accounting/dbgen"
	"packages/accounting/domain"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

const (
	AccountIDAssets                = 1
	AccountIDLiabilities           = 2
	AccountIDIncome                = 3
	AccountIDIncomeUncategorized   = 3000
	AccountIDExpenses              = 4
	AccountIDExpensesUncategorized = 4000
	AccountIDEquity                = 5
)

// https://www.postgresql.org/docs/current/errcodes-appendix.html
const PgErrCodeUniqueViolation = "23505"

type Repository struct {
	pool *pgxpool.Pool
	db   *dbgen.Queries
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
		db:   dbgen.New(pool),
	}
}

var ErrDuplicateTransaction = errors.New("transaction already exists")

// TODO: attempt to categorize
func (repo *Repository) CreateTransactions(ctx context.Context, transactions []domain.BankTransaction) error {
	tx, err := repo.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to create postgres transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	qtx := repo.db.WithTx(tx)

	for _, t := range transactions {
		posting, err := qtx.CreatePosting(ctx, NewCreatePostingParams(t))
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			if pgErr.Code == PgErrCodeUniqueViolation {
				return fmt.Errorf("%w: %s", ErrDuplicateTransaction, pgErr.Detail)
			}
		} else if err != nil {
			slog.ErrorContext(ctx, "failed to create posting", slog.Any("error", err), slog.Any("transaction", t))
			return fmt.Errorf("failed to create posting: %w", err)
		}

		switch t.TransactionType {
		case domain.TransactionSourceBank:
			accountID := AccountIDExpensesUncategorized
			if t.IsCredit() {
				accountID = AccountIDIncomeUncategorized
			}

			createEntryParams := NewCreateEntryParams(t)
			createEntryParams.LedgerAccountsID = int64(accountID)
			createEntryParams.PostingsID = int64(posting.ID)

			_, err = qtx.CreateEntry(ctx, createEntryParams)
			if err != nil {
				return fmt.Errorf("failed to create entry %w", err)
			}
		case domain.TransactionSourceCreditCard:
			createEntryParams := NewCreateEntryParams(t)
			createEntryParams.LedgerAccountsID = int64(AccountIDExpensesUncategorized)
			createEntryParams.PostingsID = int64(posting.ID)

			_, err = qtx.CreateEntry(ctx, createEntryParams)
			if err != nil {
				return fmt.Errorf("failed to create entry %w", err)
			}
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit postgres transaction: %w", err)
	}

	return nil
}

func NewCreatePostingParams(tx domain.BankTransaction) dbgen.CreatePostingParams {
	return dbgen.CreatePostingParams{
		Description:  pgtype.Text{String: tx.Description, Valid: true},
		SystemNotes:  pgtype.Text{String: "", Valid: false},
		TransactedAt: tx.Date,
		SourceHash: pgtype.Text{String: strings.Join(
			[]string{tx.SourceName, tx.RawRow}, "|",
		), Valid: true},
	}
}

func NewCreateEntryParams(tx domain.BankTransaction) dbgen.CreateEntryParams {
	mil := decimal.NewFromInt(1_000_000)

	return dbgen.CreateEntryParams{
		Description:    pgtype.Text{String: tx.Description, Valid: true},
		SystemNotes:    pgtype.Text{String: "", Valid: false},
		DebitMicrosgd:  tx.Debit.Mul(mil).IntPart(),
		CreditMicrosgd: tx.Credit.Mul(mil).IntPart(),
	}
}
