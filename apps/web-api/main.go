package main

import (
	"apps/web-api/internal/config"
	"apps/web-api/internal/db/dbgen"
	"apps/web-api/internal/server"
	"apps/web-api/internal/webapi"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	ExitCode_FailedPgPoolCreation = 1001
	ExitCode_FailedDBQuerier      = 1002
)

func main() {
	ctx := context.Background()
	cfg := config.Init(ctx)

	pool, err := pgxpool.New(ctx, cfg.PostgresDSN)
	if err != nil {
		slog.ErrorContext(ctx, "could not create postgres connection pool", slog.Any("error", err))
		os.Exit(ExitCode_FailedPgPoolCreation)
		return
	}
	defer pool.Close()

	queries := dbgen.New(pool)

	serverImpl := server.NewServer(pool, queries)
	strictHandler := webapi.NewStrictHandler(serverImpl, nil)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	webapi.HandlerFromMux(strictHandler, r)

	srv := &http.Server{
		Handler: r,
		Addr:    cfg.Address,
	}

	err = srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.ErrorContext(ctx, "unexpected server shutdown", slog.Any("error", err))
		return
	}
}
