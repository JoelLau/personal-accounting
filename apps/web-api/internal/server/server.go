package server

import (
	"apps/web-api/internal/db/dbgen"
	"apps/web-api/internal/webapi"
	"context"
	"log/slog"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oapi-codegen/runtime/types"
	"github.com/shopspring/decimal"
)

// TODO: refactor to move logic to repo layer
type Server struct {
	pool    *pgxpool.Pool
	queries *dbgen.Queries
}

var _ webapi.StrictServerInterface = &Server{}

func NewServer(pool *pgxpool.Pool, queries *dbgen.Queries) *Server {
	return &Server{
		pool:    pool,
		queries: queries,
	}
}

// GetApiV1AccountingEntries implements [webapi.StrictServerInterface].
func (s *Server) GetApiV1AccountingEntries(ctx context.Context, request webapi.GetApiV1AccountingEntriesRequestObject) (webapi.GetApiV1AccountingEntriesResponseObject, error) {
	accounts, err := s.queries.ListEntries(ctx)
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

// GetApiV1AccountingLedgerAccounts implements [webapi.StrictServerInterface].
func (s *Server) GetApiV1AccountingLedgerAccounts(ctx context.Context, request webapi.GetApiV1AccountingLedgerAccountsRequestObject) (webapi.GetApiV1AccountingLedgerAccountsResponseObject, error) {
	accounts, err := s.queries.ListLedgerAccounts(ctx)
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

// GetApiV1AccountingPostings implements [webapi.StrictServerInterface].
func (s *Server) GetApiV1AccountingPostings(ctx context.Context, request webapi.GetApiV1AccountingPostingsRequestObject) (webapi.GetApiV1AccountingPostingsResponseObject, error) {
	postings, err := s.queries.ListPostings(ctx)
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
