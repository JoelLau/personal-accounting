package main

import (
	"context"
	"errors"
	"fmt"
	"libs/ledger/application"
	"libs/ledger/application/commands"
	"libs/ledger/application/handlers"
	"libs/ledger/application/services"
	"libs/ledger/infrastructure/database"
	"log/slog"
	"os"
	"packages/ingestion/parsers"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v3"
)

// TODO: move this to config file or sth
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
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug})))
	slog.InfoContext(ctx, "starting...")

	cmd := NewCLICommand()
	if err := cmd.Run(ctx, os.Args); err != nil {
		slog.ErrorContext(ctx, "unexpected error", slog.Any("error", err))
		os.Exit(1)
	}

	slog.InfoContext(ctx, "closing...")
}

var (
	CLIFlagFileType = cli.StringFlag{
		Name:     "file-type",
		Required: true,
		Usage:    "what type of file this is - 'dbs_cc' or 'ocbc_stmt'",
	}

	CLIFlagYearMonthFilter = cli.TimestampFlag{
		Name:     "month-filter",
		Required: true,
		Usage:    "ONLY ingest files for listed month in yyyy-mm e.g. 2026-03",
		Config: cli.TimestampConfig{
			Timezone: time.UTC,
			Layouts:  []string{"2006-01"},
		},
	}
)

func NewCLICommand() *cli.Command {
	return &cli.Command{
		Name:  "Ledger Import",
		Usage: "imports bank transaction files to postgres database at env var DB_URL",
		Flags: []cli.Flag{
			&CLIFlagFileType,
			&CLIFlagYearMonthFilter,
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			args, err := NewSuppliedArgs(cmd)
			if err != nil {
				return fmt.Errorf("argument error: %w", err)
			}

			pool, err := pgxpool.New(ctx, args.PostgresDSN)
			if err != nil {
				slog.ErrorContext(ctx, "failed to connect to postgres db", slog.Any("error", err))
				os.Exit(1)
			}

			repo := database.NewRepository(pool)
			service := services.NewImportTransactionsService(repo)
			handler := handlers.NewImportTransactionsHandler(service)

			file, err := os.Open(args.FilePath)
			if err != nil {
				slog.Error("failed to open file", slog.Any("filepath", args.FilePath), slog.Any("error", err))
				os.Exit(1)
			}
			defer file.Close()

			var parser commands.TransactionFileParser
			var profile commands.ImportProfile
			switch *&args.FileType {
			case "dbs_cc":
				parser = parsers.NewDbsCreditCardCsvParser(args.MonthFilter.Year(), int(args.MonthFilter.Month()))
				profile = application.NewDBSImportProfile(AccountIDs.ExpensesUncategorized, AccountIDs.LiabilitiesCreditCard)
			case "ocbc_stmt":
				parser = parsers.NewOcbcStatementCsvParser(args.MonthFilter.Year(), int(args.MonthFilter.Month()))
				profile = application.NewOCBCStatementProfile(
					AccountIDs.Assets,
					AccountIDs.ExpensesUncategorized,
					AccountIDs.IncomeUncategorized,
				)
			default:
				return fmt.Errorf("invalid file type: '%s' (must be dbs_cc or ocbc_stmt)", args.FileType)
			}

			err = handler.Handle(ctx, commands.ImportTransactionsCommand{
				Reader:  file,
				Parser:  parser,
				Profile: profile,
			})
			if err != nil {
				return fmt.Errorf("failed to ingest file: %w", err)
			}

			return nil
		},
	}
}

type SuppliedArgs struct {
	MonthFilter time.Time
	FilePath    string
	FileType    string // either "dbs_cc", "ocbc_stmt"
	PostgresDSN string
}

func NewSuppliedArgs(cmd *cli.Command) (*SuppliedArgs, error) {
	filePath := cmd.Args().First()
	if filePath == "" {
		return nil, errors.New("must supply path to file")
	}

	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading dotenv: %w", err)
	}

	return &SuppliedArgs{
		MonthFilter: cmd.Timestamp(CLIFlagYearMonthFilter.Name),
		FilePath:    filePath,
		FileType:    cmd.String(CLIFlagFileType.Name),
		PostgresDSN: os.Getenv("DB_URL"),
	}, nil
}
