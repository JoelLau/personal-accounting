package main

import (
	"apps/web-api/internal/config"
	"apps/web-api/internal/server"
	"apps/web-api/internal/webapi"
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	ctx := context.Background()
	cfg := config.Init(ctx)

	myStrictServer := server.NewServer()
	strictHandler := webapi.NewStrictHandler(myStrictServer, nil)
	handler := webapi.HandlerWithOptions(strictHandler, webapi.ChiServerOptions{
		BaseRouter: chi.NewRouter(),
		Middlewares: []webapi.MiddlewareFunc{
			middleware.Logger,
			middleware.Recoverer,
		},
	})

	srv := &http.Server{
		Handler: handler,
		Addr:    cfg.Address,
	}

	err := srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.ErrorContext(ctx, "unexpected server shutdown", slog.Any("error", err))
		return
	}
}
