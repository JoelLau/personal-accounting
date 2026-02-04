package server

import (
	"apps/web-api/internal/webapi"
	"context"
	"log/slog"
)

type Server struct {
	webapi.Unimplemented
}

var _ webapi.StrictServerInterface = &Server{}

func NewServer() *Server {
	return &Server{}
}

// Balance and health check
// (GET /api/v1/accounting/accounts/info)
func (s *Server) GetApiV1AccountingAccountsInfo(ctx context.Context, request webapi.GetApiV1AccountingAccountsInfoRequestObject) (webapi.GetApiV1AccountingAccountsInfoResponseObject, error) {
	return nil, nil
}

// List all transactions
// (GET /api/v1/accounting/transactions)
func (s *Server) GetApiV1AccountingTransactions(ctx context.Context, request webapi.GetApiV1AccountingTransactionsRequestObject) (webapi.GetApiV1AccountingTransactionsResponseObject, error) {
	return nil, nil
}

// Create transaction
// (POST /api/v1/accounting/transactions)
func (s *Server) PostApiV1AccountingTransactions(ctx context.Context, request webapi.PostApiV1AccountingTransactionsRequestObject) (webapi.PostApiV1AccountingTransactionsResponseObject, error) {
	return nil, nil
}

// Update or finalize a generated transaction
// (PUT /api/v1/accounting/transactions/{transaction_id})
func (s *Server) PutApiV1AccountingTransactionsTransactionId(ctx context.Context, request webapi.PutApiV1AccountingTransactionsTransactionIdRequestObject) (webapi.PutApiV1AccountingTransactionsTransactionIdResponseObject, error) {
	return nil, nil
}

// (GET /livez)
func (s *Server) GetLivez(ctx context.Context, request webapi.GetLivezRequestObject) (webapi.GetLivezResponseObject, error) {
	return webapi.GetLivez200JSONResponse{}, nil
}

// (GET /readyz)
func (s *Server) GetReadyz(ctx context.Context, request webapi.GetReadyzRequestObject) (webapi.GetReadyzResponseObject, error) {
	slog.InfoContext(ctx, "GetReadyz", slog.Any("request", request))
	return webapi.GetReadyz200JSONResponse{}, nil
}
