package application

import (
	"context"
	"errors"
	"fmt"
	"io"
	domain "packages/accounting/domain"
)

type TransactionRepository interface {
	CreateTransactions(ctx context.Context, transactions []domain.BankTransaction) error
}

type TransactionParser interface {
	Parse(file io.Reader) ([]domain.BankTransaction, error)
}

type IngestionService struct {
	repo    TransactionRepository
	parsers []TransactionParser
}

func NewIngestionService(repo TransactionRepository, parsers []TransactionParser) *IngestionService {
	return &IngestionService{
		repo:    repo,
		parsers: parsers,
	}
}

func (s *IngestionService) IngestFile(ctx context.Context, file io.Reader) error {
	transactions, err := s.parseFile(file)
	if err != nil {
		return fmt.Errorf("failed to parse file: %w", err)
	}

	err = s.repo.CreateTransactions(ctx, transactions)
	if err != nil {
		return fmt.Errorf("failed to persist transactions: %w", err)
	}

	return nil
}

func (s *IngestionService) parseFile(file io.Reader) ([]domain.BankTransaction, error) {
	if len(s.parsers) == 0 {
		return nil, fmt.Errorf("no parsers")
	}

	seeker, isSeeker := file.(io.ReadSeeker)

	var errs []error
	for _, parser := range s.parsers {
		if isSeeker {
			if _, err := seeker.Seek(0, io.SeekStart); err != nil {
				return nil, fmt.Errorf("failed to reset file reader")
			}
		}

		txs, err := parser.Parse(file)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		if len(txs) == 0 {
			errs = append(errs, fmt.Errorf("no transactions parsed %T", parser))
			continue
		}

		return txs, err
	}

	return nil, fmt.Errorf("failed to parse transactions: %w", errors.Join(errs...))
}
