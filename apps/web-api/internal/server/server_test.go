package server_test

import (
	"apps/web-api/internal/server"
	"apps/web-api/internal/webapi"
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
	// Arrange
	testEnv := NewTestEnv(t)

	requestURL, err := url.JoinPath(testEnv.server.URL, "/api/readyz")
	require.NoError(t, err)

	// Act
	response, err := http.Get(requestURL)

	// Assert
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.StatusCode)
}

type testEnv struct {
	pool   *pgxpool.Pool
	server *httptest.Server
}

func NewTestEnv(t *testing.T) testEnv {
	ctx := t.Context()

	postgresContainer, err := postgres.Run(ctx, "postgres:18-alpine")
	t.Cleanup(func() {
		if err := testcontainers.TerminateContainer(postgresContainer); err != nil {
			t.Logf("failed to close postgres testcontainer: %v", err)
		}
	})
	require.NoError(t, err)

	dsn, err := postgresContainer.ConnectionString(ctx)
	require.NoError(t, err)

	pool, err := pgxpool.New(ctx, dsn)
	require.NoError(t, err)

	s := server.NewServer(pool)
	si := webapi.NewStrictHandler(s, nil)
	handler := webapi.HandlerWithOptions(si, webapi.ChiServerOptions{
		BaseRouter: chi.NewRouter(),
	})

	server := httptest.NewServer(handler)
	t.Cleanup(func() {
		server.Close()
	})

	return testEnv{
		pool:   pool,
		server: server,
	}
}
