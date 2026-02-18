package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"packages/accounting/dbgen"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jszwec/csvutil"
	"github.com/shopspring/decimal"

	"github.com/joho/godotenv"
)

var MicroSGDMulFactor = decimal.NewFromInt(1_000_000)

const (
	Account_Liabilities = 2002
	Account_Expenses    = 4
)

func main() {
	ctx := context.Background()
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug})))

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		slog.ErrorContext(ctx, "failed to get DB_URL from environment variables")
		os.Exit(1)
		return
	}

	args := os.Args[1:]
	slog.InfoContext(ctx, "starting...", slog.Any("args", os.Args))

	if len(args) < 1 {
		slog.ErrorContext(ctx, "invalid file name ''")
		os.Exit(1)
		return
	}

	fileName := args[0]
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	r := csv.NewReader(file)

	for i := range linesToSkip {
		if _, err := r.Read(); err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("failed skipping headers at line %d: %v", i+1, err))
			os.Exit(1)
			return
		}
	}

	dec, err := csvutil.NewDecoder(r)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create csv util decoder", slog.Any("error", err))
	}

	var transactions []Transaction
	for {
		var tx Transaction
		if err := dec.Decode(&tx); err == io.EOF {
			break
		} else if err != nil {
			slog.WarnContext(ctx, "error decoding transaction", slog.Any("error", err))
			continue
		}
		transactions = append(transactions, tx)
	}

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		slog.ErrorContext(ctx, "failed to connect to database", slog.Any("error", err))
		os.Exit(1)
		return
	}
	defer pool.Close()

	querier := dbgen.New(pool)

	// TODO: use db tx / rollback / commit
	for _, tx := range transactions {
		ts := time.Now().Format(time.RFC3339)
		systemNotes := fmt.Sprintf("[%s] script parse DBS CC K csv on %s", tx.TransactionType, ts)

		posting, err := querier.CreatePosting(ctx, dbgen.CreatePostingParams{
			Description:  pgtype.Text{String: tx.Description, Valid: true},
			SystemNotes:  pgtype.Text{String: systemNotes, Valid: true},
			TransactedAt: tx.Date.Time,
		})
		if err != nil {
			slog.ErrorContext(ctx, "error creating new posting", slog.Any("error", err))
			continue
		}
		slog.InfoContext(ctx, "succesfully created posting", slog.Any("posting", posting))

		txCredit := decimal.NewFromInt(0)
		if strings.TrimSpace(tx.Credit) != "" {
			txCredit, err = decimal.NewFromString(tx.Credit)
			if err != nil {
				slog.ErrorContext(ctx, "could not parse tx credit string", slog.Any("error", err), slog.Any("tx", tx))
				continue
			}
		}

		txDebit := decimal.NewFromInt(0)
		if strings.TrimSpace(tx.Debit) != "" {
			txDebit, err = decimal.NewFromString(tx.Debit)
			if err != nil {
				slog.ErrorContext(ctx, "could not parse tx deFbit string", slog.Any("error", err), slog.Any("tx", tx))
				continue
			}
		}

		liability, err := querier.CreateEntry(ctx, dbgen.CreateEntryParams{
			Description:      posting.Description,
			SystemNotes:      pgtype.Text{String: systemNotes, Valid: true},
			PostingsID:       posting.ID,
			LedgerAccountsID: Account_Liabilities,
			DebitMicrosgd:    txCredit.Mul(MicroSGDMulFactor).IntPart(),
			CreditMicrosgd:   txDebit.Mul(MicroSGDMulFactor).IntPart(),
		})
		if err != nil {
			slog.ErrorContext(ctx, "failed to create liability entry", slog.Any("error", err), slog.Any("tx", tx))
			continue
		}
		slog.InfoContext(ctx, "succesfully created liability entry", slog.Any("liability", liability))

		expense, err := querier.CreateEntry(ctx, dbgen.CreateEntryParams{
			Description:      posting.Description,
			SystemNotes:      pgtype.Text{String: systemNotes, Valid: true},
			PostingsID:       posting.ID,
			LedgerAccountsID: Account_Expenses,
			DebitMicrosgd:    txDebit.Mul(MicroSGDMulFactor).IntPart(),
			CreditMicrosgd:   txCredit.Mul(MicroSGDMulFactor).IntPart(),
		})
		if err != nil {
			slog.ErrorContext(ctx, "failed to create liability entry", slog.Any("error", err), slog.Any("tx", tx))
			continue
		}
		slog.InfoContext(ctx, "succesfully created expense entry", slog.Any("expense", expense))
	}

	os.Exit(0)
}

const linesToSkip = 6

// Custom type to handle DBS date format: "12 Feb 2026"
type DBSTime struct {
	time.Time
}

func (t *DBSTime) UnmarshalCSV(data []byte) error {
	const layout = "02 Jan 2006"
	parsed, err := time.Parse(layout, string(data))
	if err != nil {
		return err
	}
	t.Time = parsed
	return nil
}

type Transaction struct {
	Date            DBSTime `csv:"Transaction Date"`
	Description     string  `csv:"Transaction Description"`
	TransactionType string  `csv:"Transaction Type"`
	PaymentType     string  `csv:"Payment Type"`
	Status          string  `csv:"Transaction Status"`
	Debit           string  `csv:"Debit Amount"`
	Credit          string  `csv:"Credit Amount"`
}
