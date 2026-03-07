package handlers

import (
	"context"
	"libs/ledger/application/commands"
	"libs/ledger/application/services"
)

type ImportCSVHandler struct {
	service services.ImportTransactionService
}

func NewImportTransactionsHandler(service services.ImportTransactionService) *ImportCSVHandler {
	return &ImportCSVHandler{
		service: service,
	}
}

func (h *ImportCSVHandler) Handle(ctx context.Context, cmd commands.ImportTransactionsCommand) error {
	return h.service.Import(ctx, cmd)
}
