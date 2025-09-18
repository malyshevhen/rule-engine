package e2e

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/testcontainers/testcontainers-go/wait"

	redisPkg "github.com/malyshevhen/rule-engine/internal/storage/redis"
)

// TestEnvironment holds all test infrastructure components
type TestEnvironment struct {
	PostgresContainer   testcontainers.Container
	RedisContainer      testcontainers.Container
	NATSContainer       testcontainers.Container
	HoverflyContainer   testcontainers.Container
	RuleEngineContainer testcontainers.Container
}

// SetupTestEnvironment creates and starts all required test infrastructure
func SetupTestEnvironment(ctx context.Context, t *testing.T) (*TestEnvironment, func()) {
	t.Helper()

	// Start environment containers
	stack, err := compose.NewDockerCompose("containers/compose.yaml")
	require.NoError(t, err)

	cleanup := func() {
		if err = stack.Down(ctx, compose.RemoveOrphans(true), compose.RemoveImagesLocal); err != nil {
			require.NoError(t, err, "Failed to cleanup test environment: %v")
		}
	}

	err = stack.WaitForService("rule-engine", wait.ForHealthCheck()).Up(ctx)
	require.NoError(t, err)

	pgContainer, err := stack.ServiceContainer(ctx, "postgres")
	require.NoError(t, err)
	redisContainer, err := stack.ServiceContainer(ctx, "redis")
	require.NoError(t, err)
	natsContainer, err := stack.ServiceContainer(ctx, "nats")
	require.NoError(t, err)
	hoverflyContainer, err := stack.ServiceContainer(ctx, "hoverfly")
	require.NoError(t, err)
	ruleEngineContainer, err := stack.ServiceContainer(ctx, "rule-engine")
	require.NoError(t, err)

	// Start Rule Engine as host process
	env := &TestEnvironment{
		PostgresContainer:   pgContainer,
		RedisContainer:      redisContainer,
		NATSContainer:       natsContainer,
		HoverflyContainer:   hoverflyContainer,
		RuleEngineContainer: ruleEngineContainer,
	}

	return env, cleanup
}

// GetRuleEngineURL returns the Rule Engine service URL
func (env *TestEnvironment) GetRuleEngineURL(ctx context.Context, t *testing.T) string {
	t.Helper()
	return "http://localhost:8080"
}

// GetHoverflyAdminURL returns the Hoverfly admin API URL
func (env *TestEnvironment) GetHoverflyAdminURL(ctx context.Context, t *testing.T) string {
	t.Helper()
	port, err := env.HoverflyContainer.MappedPort(ctx, "8888")
	require.NoError(t, err)
	return fmt.Sprintf("http://localhost:%s", port.Port())
}

// GetPostgresConnectionString returns the PostgreSQL connection string
func (env *TestEnvironment) GetPostgresConnectionString(ctx context.Context, t *testing.T) string {
	t.Helper()
	host, err := env.PostgresContainer.Host(ctx)
	require.NoError(t, err)
	port, err := env.PostgresContainer.MappedPort(ctx, "5432")
	require.NoError(t, err)
	return fmt.Sprintf("postgres://postgres:password@%s:%s/rule_engine_test?sslmode=disable", host, port.Port())
}

// GetRedisClient returns a Redis client connected to the test container
func (env *TestEnvironment) GetRedisClient(ctx context.Context, t *testing.T) *redisPkg.Client {
	t.Helper()
	host, err := env.RedisContainer.Host(ctx)
	require.NoError(t, err)
	port, err := env.RedisContainer.MappedPort(ctx, "6379")
	require.NoError(t, err)

	config := &redisPkg.Config{
		Addr: fmt.Sprintf("%s:%s", host, port.Port()),
	}
	return redisPkg.NewClient(config)
}

// GetNATSURL returns the NATS connection URL
func (env *TestEnvironment) GetNATSURL(ctx context.Context, t *testing.T) string {
	t.Helper()
	host, err := env.NATSContainer.Host(ctx)
	require.NoError(t, err)
	port, err := env.NATSContainer.MappedPort(ctx, "4222")
	require.NoError(t, err)
	return fmt.Sprintf("nats://%s:%s", host, port.Port())
}

// SetupDatabasePool creates a PostgreSQL connection pool for testing
func (env *TestEnvironment) SetupDatabasePool(ctx context.Context, t *testing.T) *pgxpool.Pool {
	t.Helper()

	connStr := env.GetPostgresConnectionString(ctx, t)

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
func (env *TestEnvironment) WaitForServices(ctx context.Context, t *testing.T) {
	t.Helper()

	// Test Hoverfly admin API first (it's faster)
	hoverflyAdminURL := env.GetHoverflyAdminURL(ctx, t)
	resp, err := http.Get(hoverflyAdminURL + "/api/v2/hoverfly")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	// Test Rule Engine health endpoint
	ruleEngineURL := env.GetRuleEngineURL(ctx, t)
	resp, err = http.Get(ruleEngineURL + "/health")
	require.NoError(t, err)
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Health check failed with status %d: %s", resp.StatusCode, string(body))
	}
	resp.Body.Close()
}

// SetupHoverflySimulation loads a simulation into Hoverfly
func (env *TestEnvironment) SetupHoverflySimulation(ctx context.Context, t *testing.T, simulationData string) {
	t.Helper()

	hoverflyAdminURL := env.GetHoverflyAdminURL(ctx, t)
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequestWithContext(ctx, "PUT", hoverflyAdminURL+"/api/v2/simulation", strings.NewReader(simulationData))
	if err != nil {
		t.Fatal(err, "Failed to create request")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected status code %d, got %d; Response: %s", http.StatusOK, resp.StatusCode, string(body))
	}
}
