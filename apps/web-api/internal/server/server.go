package server

import (
	"apps/web-api/internal/db/dbgen"
	"apps/web-api/internal/webapi"
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

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
