package server

import (
	"apps/web-api/internal/webapi"
	"context"
	"fmt"
	dbgen "libs/ledger/infrastructure/database/gen"
	"log/slog"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oapi-codegen/runtime/types"
	"github.com/shopspring/decimal"
)

// TODO: refactor - server handles queries and validation, add repo layer to handle persistence
type Server struct {
	pool *pgxpool.Pool
}

var _ webapi.StrictServerInterface = &Server{}

func NewServer(pool *pgxpool.Pool) *Server {
	return &Server{
		pool: pool,
	}
}

func (s *Server) querier() *dbgen.Queries {
	return dbgen.New(s.pool)
}

// List all entries
// (GET /api/v1/accounting/entries)
func (s *Server) GetApiV1AccountingEntries(ctx context.Context, request webapi.GetApiV1AccountingEntriesRequestObject) (webapi.GetApiV1AccountingEntriesResponseObject, error) {
	accounts, err := s.querier().ListEntries(ctx)
	if err != nil {
		return nil, err
	}

	data := make([]webapi.Entry, len(accounts))
	for idx, account := range accounts {
		credit := decimal.NewFromInt(account.CreditMicrosgd).Div(decimal.NewFromInt(1_000_000))
		debit := decimal.NewFromInt(account.DebitMicrosgd).Div(decimal.NewFromInt(1_000_000))

		data[idx] = webapi.Entry{
			CreditAmount:     credit.String(),
			DebitAmount:      debit.String(),
			Description:      account.Description.String,
			Id:               strconv.FormatInt(account.ID, 10),
			LedgerAccountsId: strconv.FormatInt(account.LedgerAccountsID, 10),
			PostingsId:       strconv.FormatInt(account.PostingsID, 10),
			SystemNotes:      account.SystemNotes.String,
		}
	}

	return webapi.GetApiV1AccountingEntries200JSONResponse{Data: &data}, nil
}

// List all ledger accounts
// (GET /api/v1/accounting/ledger_accounts)
func (s *Server) GetApiV1AccountingLedgerAccounts(ctx context.Context, request webapi.GetApiV1AccountingLedgerAccountsRequestObject) (webapi.GetApiV1AccountingLedgerAccountsResponseObject, error) {
	accounts, err := s.querier().ListLedgerAccounts(ctx)
	if err != nil {
		return nil, err
	}

	data := make([]webapi.LedgerAccount, len(accounts))
	for idx, account := range accounts {
		var parentID *string
		if account.ParentID.Valid {
			val := strconv.FormatInt(account.ParentID.Int64, 10)
			parentID = &val
		}

		data[idx] = webapi.LedgerAccount{
			Id:            strconv.FormatInt(account.ID, 10),
			Name:          account.Name,
			QualifiedName: account.QualifiedName,
			Description:   account.Description.String,
			ParentId:      parentID,
		}
	}

	return webapi.GetApiV1AccountingLedgerAccounts200JSONResponse{Data: &data}, nil
}

// List all postings
// (GET /api/v1/accounting/postings)
func (s *Server) GetApiV1AccountingPostings(ctx context.Context, request webapi.GetApiV1AccountingPostingsRequestObject) (webapi.GetApiV1AccountingPostingsResponseObject, error) {
	postings, err := s.querier().ListPostings(ctx)
	if err != nil {
		return nil, err
	}

	data := make([]webapi.Posting, len(postings))
	for idx, posting := range postings {
		data[idx] = webapi.Posting{
			Id:           strconv.FormatInt(posting.ID, 10),
			Description:  posting.Description.String,
			SystemNotes:  posting.SystemNotes.String,
			TransactedAt: types.Date{Time: posting.TransactedAt},
		}
	}

	return webapi.GetApiV1AccountingPostings200JSONResponse{Data: &data}, nil
}

var mil = decimal.NewFromInt(1_000_000)

// (PUT /api/v1/accounting/entries/:entry_id)
func (s *Server) UpdateAccountingEntry(ctx context.Context, request webapi.UpdateAccountingEntryRequestObject) (webapi.UpdateAccountingEntryResponseObject, error) {
	slog.InfoContext(ctx, "/api/v1/accounting/entries/:entry_id", slog.Any("request", request))

	// validation
	if request.Body == nil {
		return webapi.UpdateAccountingEntry400JSONResponse{
			Status: "400",
			Title:  "Bad Request - invalid request body",
		}, nil
	}

	body := *request.Body
	entryID, err := strconv.ParseInt(request.EntryId, 10, 64)
	if err != nil {
		return webapi.UpdateAccountingEntry400JSONResponse{
			Status: "400",
			Title:  "Bad Request - 'entry_id' must be numeric (e.g. \"72\")",
		}, nil
	}

	ledgerAccountsId, err := strconv.ParseInt(body.LedgerAccountsId, 10, 64)
	if err != nil {
		return webapi.UpdateAccountingEntry400JSONResponse{
			Status: "400",
			Title:  "Bad Request - 'ledger_accounts_id' must be numeric (e.g. \"72\")",
		}, nil
	}

	if request.EntryId != body.Id {
		return webapi.UpdateAccountingEntry400JSONResponse{
			Status: "400",
			Title:  fmt.Sprintf("Bad Request - entry id in url ('%s') does not match entry id found in body ('%s')", request.EntryId, body.Id),
		}, nil
	}

	debit, err := decimal.NewFromString(strings.TrimSpace(strings.ReplaceAll(body.DebitAmount, ",", "")))
	if err != nil {
		return webapi.UpdateAccountingEntry400JSONResponse{
			Status: "400",
			Title:  fmt.Sprintf("Bad Request - debit_amount (%v) must be numeric (e.g. \"72.00\"): %v", body.DebitAmount, err),
		}, nil
	}

	credit, err := decimal.NewFromString(strings.TrimSpace(strings.ReplaceAll(body.CreditAmount, ",", "")))
	if err != nil {
		return webapi.UpdateAccountingEntry400JSONResponse{
			Status: "400",
			Title:  fmt.Sprintf("Bad Request - credit_amount (%v) must be numeric (e.g. \"72.00\"): %v", body.CreditAmount, err),
		}, nil
	}

	err = s.querier().UpdateEntry(ctx, dbgen.UpdateEntryParams(dbgen.UpdateEntryParams{
		ID:               entryID,
		Description:      pgtype.Text{String: body.Description, Valid: true},
		SystemNotes:      pgtype.Text{String: "", Valid: false}, // WARN: this will clear system notes - we might need to do a pre-fetch for old values
		LedgerAccountsID: ledgerAccountsId,
		DebitMicrosgd:    debit.Mul(mil).IntPart(),
		CreditMicrosgd:   credit.Mul(mil).IntPart(),
	}))
	if err != nil {
		return nil, fmt.Errorf("failed to update database 'entry' object: %w", err)
	}

	return webapi.UpdateAccountingEntry200JSONResponse{}, nil
}

// Balance and health check
// (GET /api/v1/accounting/accounts/info)
func (s *Server) GetApiV1AccountingAccountsInfo(ctx context.Context, request webapi.GetApiV1AccountingAccountsInfoRequestObject) (webapi.GetApiV1AccountingAccountsInfoResponseObject, error) {
	return webapi.GetApiV1AccountingAccountsInfo200JSONResponse{}, nil
}

// List all transactions
// (GET /api/v1/accounting/transactions)
func (s *Server) GetApiV1AccountingTransactions(ctx context.Context, request webapi.GetApiV1AccountingTransactionsRequestObject) (webapi.GetApiV1AccountingTransactionsResponseObject, error) {
	return webapi.GetApiV1AccountingTransactions200JSONResponse{}, nil
}

// Create transaction
// (POST /api/v1/accounting/transactions)
func (s *Server) PostApiV1AccountingTransactions(ctx context.Context, request webapi.PostApiV1AccountingTransactionsRequestObject) (webapi.PostApiV1AccountingTransactionsResponseObject, error) {
	return webapi.PostApiV1AccountingTransactions201Response{}, nil
}

// Update or finalize a generated transaction
// (PUT /api/v1/accounting/transactions/{transaction_id})
func (s *Server) PutApiV1AccountingTransactionsTransactionId(ctx context.Context, request webapi.PutApiV1AccountingTransactionsTransactionIdRequestObject) (webapi.PutApiV1AccountingTransactionsTransactionIdResponseObject, error) {
	return webapi.PutApiV1AccountingTransactionsTransactionId200Response{}, nil
}

// (GET /livez)
func (s *Server) GetApiLivez(ctx context.Context, request webapi.GetApiLivezRequestObject) (webapi.GetApiLivezResponseObject, error) {
	return webapi.GetApiLivez200JSONResponse{}, nil
}

// (GET /readyz)
func (s *Server) GetApiReadyz(ctx context.Context, request webapi.GetApiReadyzRequestObject) (webapi.GetApiReadyzResponseObject, error) {
	slog.InfoContext(ctx, "GetApiReadyz", slog.Any("request", request))
	return webapi.GetApiReadyz200JSONResponse{}, nil
}
