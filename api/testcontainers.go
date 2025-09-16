package api

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	redisPkg "github.com/malyshevhen/rule-engine/internal/storage/redis"

	// Import the postgres driver for testcontainers
	_ "github.com/lib/pq"
)

// TestContainers holds the test container instances
type TestContainers struct {
	PostgresContainer testcontainers.Container
	RedisContainer    testcontainers.Container
	NATSContainer     testcontainers.Container
}

// SetupTestContainers creates and starts all required test containers
func SetupTestContainers(ctx context.Context, t *testing.T) (*TestContainers, func()) {
	t.Helper()

	// Start PostgreSQL container
	pgContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:15-alpine",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_DB":       "rule_engine_test",
				"POSTGRES_USER":     "postgres",
				"POSTGRES_PASSWORD": "password",
			},
			WaitingFor: wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30 * time.Second),
		},
		Started: true,
	})
	require.NoError(t, err)

	// Start Redis container
	redisContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "redis:7-alpine",
			ExposedPorts: []string{"6379/tcp"},
			Cmd:          []string{"redis-server", "--appendonly", "yes"},
			WaitingFor:   wait.ForLog("Ready to accept connections").WithStartupTimeout(30 * time.Second),
		},
		Started: true,
	})
	require.NoError(t, err)

	// Start NATS container
	natsContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "nats:2.10-alpine",
			ExposedPorts: []string{"4222/tcp"},
			WaitingFor:   wait.ForLog("Server is ready").WithStartupTimeout(30 * time.Second),
		},
		Started: true,
	})
	require.NoError(t, err)

	tc := &TestContainers{
		PostgresContainer: pgContainer,
		RedisContainer:    redisContainer,
		NATSContainer:     natsContainer,
	}

	// Cleanup function
	cleanup := func() {
		require.NoError(t, natsContainer.Terminate(ctx))
		require.NoError(t, redisContainer.Terminate(ctx))
		require.NoError(t, pgContainer.Terminate(ctx))
	}

	return tc, cleanup
}

// GetPostgresConnectionString returns the PostgreSQL connection string
func (tc *TestContainers) GetPostgresConnectionString(ctx context.Context, t *testing.T) string {
	t.Helper()
	host, err := tc.PostgresContainer.Host(ctx)
	require.NoError(t, err)
	port, err := tc.PostgresContainer.MappedPort(ctx, "5432")
	require.NoError(t, err)
	return fmt.Sprintf("postgres://postgres:password@%s:%s/rule_engine_test?sslmode=disable", host, port.Port())
}

// GetRedisClient returns a Redis client connected to the test container
func (tc *TestContainers) GetRedisClient(ctx context.Context, t *testing.T) *redisPkg.Client {
	t.Helper()
	host, err := tc.RedisContainer.Host(ctx)
	require.NoError(t, err)
	port, err := tc.RedisContainer.MappedPort(ctx, "6379")
	require.NoError(t, err)

	config := &redisPkg.Config{
		Addr: fmt.Sprintf("%s:%s", host, port.Port()),
	}
	return redisPkg.NewClient(config)
}

// GetNATSURL returns the NATS connection URL
func (tc *TestContainers) GetNATSURL(ctx context.Context, t *testing.T) string {
	t.Helper()
	host, err := tc.NATSContainer.Host(ctx)
	require.NoError(t, err)
	port, err := tc.NATSContainer.MappedPort(ctx, "4222")
	require.NoError(t, err)
	return fmt.Sprintf("nats://%s:%s", host, port.Port())
}

// SetupDatabasePool creates a PostgreSQL connection pool for testing
func (tc *TestContainers) SetupDatabasePool(ctx context.Context, t *testing.T) *pgxpool.Pool {
	t.Helper()

	connStr := tc.GetPostgresConnectionString(ctx, t)

	// Configure connection pool for testing
	config, err := pgxpool.ParseConfig(connStr)
	require.NoError(t, err)

	config.MaxConns = 10
	config.MinConns = 1

	pool, err := pgxpool.NewWithConfig(ctx, config)
	require.NoError(t, err)

	// Test the connection
	require.NoError(t, pool.Ping(ctx))

	return pool
}

// WaitForServices waits for all services to be ready
func (tc *TestContainers) WaitForServices(ctx context.Context, t *testing.T) {
	t.Helper()

	// Test PostgreSQL connection
	connStr := tc.GetPostgresConnectionString(ctx, t)
	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)
	defer db.Close()

	require.NoError(t, db.Ping())

	// Test Redis connection
	rdb := tc.GetRedisClient(ctx, t)
	defer rdb.Close()

	require.NoError(t, rdb.Ping(ctx))
}
