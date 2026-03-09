package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"libs/ledger/application"
	"libs/ledger/application/commands"
	"libs/ledger/application/handlers"
	"libs/ledger/application/services"
	"libs/ledger/infrastructure/database"
	"log/slog"
	"os"
	"packages/ingestion/parsers"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var AccountIDs = struct {
	Assets                int64
	ExpensesUncategorized int64 // "Expenses:Uncategorized"
	IncomeUncategorized   int64
	LiabilitiesCreditCard int64 // "Liabilities:Credit Card"
}{
	Assets:                1,
	ExpensesUncategorized: 4000,
	IncomeUncategorized:   3000,
	LiabilitiesCreditCard: 2001,
}

func main() {
	ctx := context.Background()
	cfg, err := setup()
	if err != nil {
		slog.ErrorContext(ctx, "failed to setup application", slog.Any("error", err))
		os.Exit(1)
	}
	slog.SetDefault(cfg.Logger)

	slog.InfoContext(ctx, "starting...")

	pool, err := pgxpool.New(ctx, cfg.PostgresDSN)
	if err != nil {
		slog.ErrorContext(ctx, "failed to connect to postgres db", slog.Any("error", err))
		os.Exit(1)
	}

	repo := database.NewRepository(pool)
	service := services.NewImportTransactionsService(repo)
	handler := handlers.NewImportTransactionsHandler(service)

	file, err := os.Open(cfg.FilePath)
	if err != nil {
		slog.Error("failed to open file", slog.Any("filepath", cfg.FilePath), slog.Any("error", err))
		os.Exit(1)
	}
	defer file.Close()

	err = handler.Handle(ctx, commands.ImportTransactionsCommand{
		Reader:  file,
		Parser:  cfg.Parser,
		Profile: cfg.Profile,
	})
	if err != nil {
		slog.ErrorContext(ctx, "failed to ingest file", slog.Any("error", err))
		os.Exit(1)
	}

	slog.InfoContext(ctx, "completed successfully")
}

func setup() (Config, error) {
	slogger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug}))

	err := godotenv.Load(".env")
	if err != nil {
		return Config{}, fmt.Errorf("failed to load environment variables, %w", err)
	}

	fileType := flag.String("filetype", "", "The bank type (dbs_cc or ocbc_stmt)")
	flag.Parse()

	if len(flag.Args()) < 1 {
		return Config{}, errors.New("no filepath in args")
	}
	filePath := flag.Args()[0]

	var parser commands.TransactionFileParser
	var profile commands.ImportProfile
	switch *fileType {
	case "dbs_cc":
		parser = parsers.NewDbsCreditCardCsvParser()
		profile = application.NewDBSImportProfile(AccountIDs.ExpensesUncategorized, AccountIDs.LiabilitiesCreditCard)
	case "ocbc_stmt":
		parser = parsers.NewOcbcStatementCsvParser()
		profile = application.NewOCBCStatementProfile(
			AccountIDs.Assets,
			AccountIDs.ExpensesUncategorized,
			AccountIDs.IncomeUncategorized,
		)
	default:
		return Config{}, fmt.Errorf("invalid -filetype: '%s' (must be dbs_cc or ocbc_stmt)", *fileType)
	}

	return Config{
		Logger:      slogger,
		PostgresDSN: os.Getenv("DB_URL"),
		FilePath:    filePath,
		Parser:      parser,
		Profile:     profile,
	}, nil
}

type Config struct {
	Logger      *slog.Logger
	PostgresDSN string
	FilePath    string
	Parser      commands.TransactionFileParser
	Profile     commands.ImportProfile
}
