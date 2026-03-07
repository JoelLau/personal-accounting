package testutil

import (
	"database/sql"
	"embed"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func NewPostgresContainer(t *testing.T, migrations embed.FS, seed embed.FS) *postgres.PostgresContainer {
	t.Helper()

	ctx := t.Context()

	postgresContainer, err := postgres.Run(ctx, "postgres:18-alpine",
		postgres.WithDatabase(t.Name()),
		postgres.WithUsername("username"),
		postgres.WithPassword("p@ssw0rd!123"),
		postgres.BasicWaitStrategies(),
	)
	require.NoError(t, err)
	t.Cleanup(func() {
		if err := testcontainers.TerminateContainer(postgresContainer); err != nil {
			t.Logf("error terminating postgres container: %v", err)
		}
	})

	connStr, err := postgresContainer.ConnectionString(ctx)
	require.NoError(t, err)

	pool, err := pgxpool.New(ctx, connStr)
	require.NoError(t, err)

	goose.SetBaseFS(migrations)
	err = goose.SetDialect("postgres")
	require.NoError(t, err)

	db := sql.OpenDB(stdlib.GetPoolConnector(pool))
	err = goose.Up(db, "migrations")
	require.NoError(t, err)

	goose.SetBaseFS(seed)
	err = goose.SetDialect("postgres")
	require.NoError(t, err)

	db = sql.OpenDB(stdlib.GetPoolConnector(pool))
	err = goose.Up(db, "seed", goose.WithNoVersioning())
	require.NoError(t, err)

	return postgresContainer
}
