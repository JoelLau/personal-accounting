package services

import (
	"context"
	"fmt"
	"libs/ledger/application/commands"
	"libs/ledger/domain"
)

type ImportTransactionService interface {
	Import(context.Context, commands.ImportTransactionsCommand) error
}

type PostingsRepository interface {
	CreatePostings(ctx context.Context, postings []domain.Posting) error
}

type ImportTransactionsServiceImpl struct {
	repo PostingsRepository
}

func NewImportTransactionsService(repo PostingsRepository) *ImportTransactionsServiceImpl {
	return &ImportTransactionsServiceImpl{
		repo: repo,
	}
}

func (s *ImportTransactionsServiceImpl) Import(ctx context.Context, cmd commands.ImportTransactionsCommand) error {
	rawTransactions, err := cmd.Parser.Parse(cmd.Reader)
	if err != nil {
		return fmt.Errorf("failed to parse file: %w", err)
	}

	postings := make([]domain.Posting, len(rawTransactions))
	for idx, t := range rawTransactions {
		posting, err := cmd.Profile.NewPosting(t)
		if err != nil {
			return fmt.Errorf("failed to parse raw transaction %v at index %d: %w", t, idx, err)
		}
		postings[idx] = posting
	}

	err = s.repo.CreatePostings(ctx, postings)
	if err != nil {
		return fmt.Errorf("failed to persist transactions: %w", err)
	}

	return nil
}
