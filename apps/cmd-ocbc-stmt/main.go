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
	Account_Assets   = 1
	Account_Expenses = 4
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
	r.FieldsPerRecord = -1
	r.LazyQuotes = true
	r.TrimLeadingSpace = true

	for {
		rec, err := r.Read()
		if err != nil {
			slog.WarnContext(ctx, "error skipping lines while reading csv file", slog.Any("error", err))
			os.Exit(1)
			return
		}
		if len(rec) == 1 && strings.TrimSpace(rec[0]) == "Transaction History" {
			break
		}
	}

	header, err := r.Read()
	if err != nil {
		slog.WarnContext(ctx, "error reading csv headers", slog.Any("error", err))
		os.Exit(1)
		return
	}

	dec, err := csvutil.NewDecoder(r, header...)
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
		systemNotes := fmt.Sprintf("script parse OCBC J+K STMT csv on %s", ts)

		posting, err := querier.CreatePosting(ctx, dbgen.CreatePostingParams{
			Description:  pgtype.Text{String: tx.Description, Valid: true},
			SystemNotes:  pgtype.Text{String: systemNotes, Valid: true},
			TransactedAt: tx.TransactionDate.Time,
		})
		if err != nil {
			slog.ErrorContext(ctx, "error creating new posting", slog.Any("error", err))
			continue
		}
		slog.InfoContext(ctx, "succesfully created posting", slog.Any("posting", posting))

		txWithdrawal := decimal.NewFromInt(0)
		if strings.TrimSpace(tx.Withdrawals) != "" {
			txWithdrawal, err = decimal.NewFromString(strings.ReplaceAll(tx.Withdrawals, ",", ""))
			if err != nil {
				slog.ErrorContext(ctx, "could not parse tx credit string", slog.Any("error", err), slog.Any("tx", tx))
				continue
			}
		}

		txDeposits := decimal.NewFromInt(0)
		if strings.TrimSpace(tx.Deposits) != "" {
			txDeposits, err = decimal.NewFromString(strings.ReplaceAll(tx.Deposits, ",", ""))
			if err != nil {
				slog.ErrorContext(ctx, "could not parse tx debit string", slog.Any("error", err), slog.Any("tx", tx))
				continue
			}
		}

		asset, err := querier.CreateEntry(ctx, dbgen.CreateEntryParams{
			Description:      posting.Description,
			SystemNotes:      pgtype.Text{String: systemNotes, Valid: true},
			PostingsID:       posting.ID,
			LedgerAccountsID: Account_Assets,
			DebitMicrosgd:    txDeposits.Mul(MicroSGDMulFactor).IntPart(),
			CreditMicrosgd:   txWithdrawal.Mul(MicroSGDMulFactor).IntPart(),
		})
		if err != nil {
			slog.ErrorContext(ctx, "failed to create asset entry", slog.Any("error", err), slog.Any("tx", tx))
			continue
		}
		slog.InfoContext(ctx, "succesfully created asset entry", slog.Any("liability", asset))

		expense, err := querier.CreateEntry(ctx, dbgen.CreateEntryParams{
			Description:      posting.Description,
			SystemNotes:      pgtype.Text{String: systemNotes, Valid: true},
			PostingsID:       posting.ID,
			LedgerAccountsID: Account_Expenses,
			DebitMicrosgd:    txWithdrawal.Mul(MicroSGDMulFactor).IntPart(),
			CreditMicrosgd:   txDeposits.Mul(MicroSGDMulFactor).IntPart(),
		})
		if err != nil {
			slog.ErrorContext(ctx, "failed to create expense entry", slog.Any("error", err), slog.Any("tx", tx))
			continue
		}
		slog.InfoContext(ctx, "succesfully created expense entry", slog.Any("expense", expense))
	}

	os.Exit(0)
}

// Custom type to handle OCBC date format: "16/02/2026"
type OCBCTime struct {
	time.Time
}

func (t *OCBCTime) UnmarshalCSV(data []byte) error {
	const layout = "02/01/2006"
	parsed, err := time.Parse(layout, string(data))
	if err != nil {
		return err
	}
	t.Time = parsed
	return nil
}

type Transaction struct {
	TransactionDate OCBCTime `csv:"Transaction date"`
	ValueDate       OCBCTime `csv:"Value date"`
	Description     string   `csv:"Description"`
	Withdrawals     string   `csv:"Withdrawals(SGD)"`
	Deposits        string   `csv:"Deposits(SGD)"`
}
