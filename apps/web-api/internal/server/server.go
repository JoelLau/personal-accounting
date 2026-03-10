package server

import (
	"apps/web-api/internal/webapi"
	"context"
	"fmt"
	dbgen "libs/ledger/infrastructure/database/gen"
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

func (s *Server) GetApiV1AccountingAccountsInfo(ctx context.Context, request webapi.GetApiV1AccountingAccountsInfoRequestObject) (webapi.GetApiV1AccountingAccountsInfoResponseObject, error) {
	panic("unimplemented")
}

func (s *Server) CreateEntry(ctx context.Context, request webapi.CreateEntryRequestObject) (webapi.CreateEntryResponseObject, error) {
	// validation
	if request.Body == nil {
		return webapi.CreateEntry400JSONResponse{
			Status: "400",
			Title:  "Bad Request - invalid request body",
		}, nil
	}

	body := *request.Body

	postingsID, err := strconv.ParseInt(body.PostingsId, 10, 64)
	if err != nil {
		return webapi.CreateEntry400JSONResponse{
			Status: "400",
			Title:  "Bad Request - 'postings_id' must be numeric (e.g. \"72\")",
		}, nil
	}

	ledgerAccountsId, err := strconv.ParseInt(body.LedgerAccountsId, 10, 64)
	if err != nil {
		return webapi.CreateEntry400JSONResponse{
			Status: "400",
			Title:  "Bad Request - 'ledger_accounts_id' must be numeric (e.g. \"72\")",
		}, nil
	}

	debit, err := decimal.NewFromString(strings.TrimSpace(strings.ReplaceAll(body.DebitAmount, ",", "")))
	if err != nil {
		return webapi.CreateEntry400JSONResponse{
			Status: "400",
			Title:  fmt.Sprintf("Bad Request - debit_amount (%v) must be numeric (e.g. \"72.00\"): %v", body.DebitAmount, err),
		}, nil
	}

	credit, err := decimal.NewFromString(strings.TrimSpace(strings.ReplaceAll(body.CreditAmount, ",", "")))
	if err != nil {
		return webapi.CreateEntry400JSONResponse{
			Status: "400",
			Title:  fmt.Sprintf("Bad Request - credit_amount (%v) must be numeric (e.g. \"72.00\"): %v", body.CreditAmount, err),
		}, nil
	}

	entry, err := s.querier().CreateEntry(ctx, dbgen.CreateEntryParams{
		Description:      pgtype.Text{String: body.Description, Valid: true},
		SystemNotes:      pgtype.Text{String: "", Valid: false},
		PostingsID:       postingsID,
		LedgerAccountsID: ledgerAccountsId,
		DebitMicrosgd:    debit.Mul(mil).IntPart(),
		CreditMicrosgd:   credit.Mul(mil).IntPart(),
	})
	if err != nil {
		err = fmt.Errorf("error persisting entry: %w", err)
		return webapi.CreateEntry500JSONResponse{
			Status: "500",
			Title:  fmt.Sprintf("Internal Server Error - %v", err),
		}, err
	}

	return webapi.CreateEntry201JSONResponse{
		Id:          strconv.FormatInt(entry.ID, 10),
		Description: entry.Description.String,
	}, nil
}

func (s *Server) DeleteEntry(ctx context.Context, request webapi.DeleteEntryRequestObject) (webapi.DeleteEntryResponseObject, error) {
	entryId, err := strconv.ParseInt(request.EntryId, 10, 64)
	if err != nil {
		return webapi.DeleteEntry400JSONResponse{
			Status: "400",
			Title:  fmt.Sprintf("Bad Request - entry_id (%v) must be numeric (e.g. \"72.00\"): %v", request.EntryId, err),
		}, nil
	}

	err = s.querier().DeleteEntry(ctx, entryId)
	if err != nil {
		err = fmt.Errorf("error persisting entry: %w", err)
		return webapi.DeleteEntry500JSONResponse{
			Status: "500",
			Title:  fmt.Sprintf("Internal Server Error - %v", err),
		}, err
	}

	return webapi.DeleteEntry200JSONResponse{}, nil
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
func (s *Server) GetEntries(ctx context.Context, request webapi.GetEntriesRequestObject) (webapi.GetEntriesResponseObject, error) {
	accounts, err := s.querier().ListEntries(ctx)
	if err != nil {
		return webapi.GetEntries500JSONResponse{
			Status: "500",
			Title:  fmt.Sprintf("Internal Server Error - %v", err),
		}, nil
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

	return webapi.GetEntries200JSONResponse{Data: &data}, nil
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
func (s *Server) UpdateEntry(ctx context.Context, request webapi.UpdateEntryRequestObject) (webapi.UpdateEntryResponseObject, error) {
	// validation
	if request.Body == nil {
		return webapi.UpdateEntry400JSONResponse{
			Status: "400",
			Title:  "Bad Request - invalid request body",
		}, nil
	}

	body := *request.Body
	entryID, err := strconv.ParseInt(request.EntryId, 10, 64)
	if err != nil {
		return webapi.UpdateEntry400JSONResponse{
			Status: "400",
			Title:  "Bad Request - 'entry_id' must be numeric (e.g. \"72\")",
		}, nil
	}

	ledgerAccountsId, err := strconv.ParseInt(body.LedgerAccountsId, 10, 64)
	if err != nil {
		return webapi.UpdateEntry400JSONResponse{
			Status: "400",
			Title:  "Bad Request - 'ledger_accounts_id' must be numeric (e.g. \"72\")",
		}, nil
	}

	if request.EntryId != body.Id {
		return webapi.UpdateEntry400JSONResponse{
			Status: "400",
			Title:  fmt.Sprintf("Bad Request - entry id in url ('%s') does not match entry id found in body ('%s')", request.EntryId, body.Id),
		}, nil
	}

	debit, err := decimal.NewFromString(strings.TrimSpace(strings.ReplaceAll(body.DebitAmount, ",", "")))
	if err != nil {
		return webapi.UpdateEntry400JSONResponse{
			Status: "400",
			Title:  fmt.Sprintf("Bad Request - debit_amount (%v) must be numeric (e.g. \"72.00\"): %v", body.DebitAmount, err),
		}, nil
	}

	credit, err := decimal.NewFromString(strings.TrimSpace(strings.ReplaceAll(body.CreditAmount, ",", "")))
	if err != nil {
		return webapi.UpdateEntry400JSONResponse{
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

	return webapi.UpdateEntry200JSONResponse{}, nil
}

// (GET /livez)
func (s *Server) GetApiLivez(ctx context.Context, request webapi.GetApiLivezRequestObject) (webapi.GetApiLivezResponseObject, error) {
	return webapi.GetApiLivez200JSONResponse{}, nil
}

// (GET /readyz)
func (s *Server) GetApiReadyz(ctx context.Context, request webapi.GetApiReadyzRequestObject) (webapi.GetApiReadyzResponseObject, error) {
	return webapi.GetApiReadyz200JSONResponse{}, nil
}
