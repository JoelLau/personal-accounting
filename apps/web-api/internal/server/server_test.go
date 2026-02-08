package server_test

import (
	"apps/web-api/internal/db/dbgen"
	"apps/web-api/internal/server"
	"apps/web-api/internal/webapi"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestServer_Healthz(t *testing.T) {
	ctx := t.Context()

	postgresContainer, err := postgres.Run(ctx, "postgres:17-alpine")
	defer func() {
		err := testcontainers.TerminateContainer(postgresContainer)
		slog.WarnContext(ctx, fmt.Sprintf("failed to terminate container: %v", err))
	}()
	require.NoError(t, err)

	dsn, err := postgresContainer.ConnectionString(ctx)
	require.NoError(t, err)

	pool, err := pgxpool.New(ctx, dsn)
	require.NoError(t, err)

	queries := dbgen.New(pool)

	s := server.NewServer(pool, queries)
	si := webapi.NewStrictHandler(s, nil)
	handler := webapi.HandlerWithOptions(si, webapi.ChiServerOptions{
		BaseRouter: chi.NewRouter(),
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	requestURL, err := url.JoinPath(server.URL, "/api/readyz")
	require.NoError(t, err)

	response, err := http.Get(requestURL)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.StatusCode)
}
