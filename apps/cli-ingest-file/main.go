package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	postgres "packages/accounting"
	"packages/ingestion/parsers"
	usecase "packages/ingestion/use-case"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()
	cfg, err := setup()
	if err != nil {
		slog.ErrorContext(ctx, "failed to setup application", slog.Any("error", err))
		os.Exit(1)
	}
	slog.SetDefault(cfg.Logger)

	slog.InfoContext(ctx, "starting...")
	parsers := []usecase.TransactionParser{
		parsers.NewDbsCreditCardCsvParser(),
		parsers.NewOcbcStatementCsvParser(),
	}

	pool, err := pgxpool.New(ctx, cfg.PostgresDSN)
	if err != nil {
		slog.ErrorContext(ctx, "failed to connect to postgres db", slog.Any("error", err))
		os.Exit(1)
	}

	repo := postgres.NewRepository(pool)
	service := usecase.NewIngestionService(repo, parsers)

	file, err := os.Open(os.Args[1])
	if err != nil {
		slog.Error("failed to open file", slog.Any("filepath", os.Args[1]), slog.Any("error", err))
		os.Exit(1)
	}

	err = service.IngestFile(ctx, file)
	if err != nil {
		slog.ErrorContext(ctx, "failed to ingest file", slog.Any("error", err))
		os.Exit(1)
	}

	slog.InfoContext(ctx, "completed successfully")
}

func setup() (Config, error) {
	slogger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug}))

	err := godotenv.Load()
	if err != nil {
		return Config{}, fmt.Errorf("failed to load environment variables, %w", err)
	}

	if len(os.Args) < 1 {
		return Config{}, errors.New("no filepath in args")
	}

	return Config{
		Logger:      slogger,
		PostgresDSN: os.Getenv("DB_URL"),
		FilePath:    os.Args[1],
	}, nil
}

type Config struct {
	Logger      *slog.Logger
	PostgresDSN string
	FilePath    string
}
