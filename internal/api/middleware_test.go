package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/malyshevhen/rule-engine/internal/storage/redis"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestLoggingMiddleware(t *testing.T) {
	// Create a test handler that does nothing
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Wrap with logging middleware
	loggingHandler := loggingMiddleware(testHandler)

	// Create a test request
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	// Execute the request
	loggingHandler.ServeHTTP(w, req)

	// Check that the request was handled (status should be 200)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRateLimitMiddleware(t *testing.T) {
	ctx := context.Background()

	// Start a Redis container
	req := testcontainers.ContainerRequest{
		Image:        "redis:7-alpine",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections").WithStartupTimeout(30 * time.Second),
	}
	redisContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer redisContainer.Terminate(ctx)

	// Get the Redis address
	host, err := redisContainer.Host(ctx)
	if err != nil {
		t.Fatal(err)
	}
	port, err := redisContainer.MappedPort(ctx, "6379")
	if err != nil {
		t.Fatal(err)
	}
	addr := host + ":" + port.Port()

	// Create the internal Redis client
	config := &redis.Config{
		Addr: addr,
	}
	rdb := redis.NewClient(config)

	// Initialize the rate limiter
	InitRedisRateLimiter(rdb)

	// Reset for testing
	ResetMiddlewareForTesting()

	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	rateLimitHandler := rateLimitMiddleware(testHandler)

	t.Run("allow requests within limit", func(t *testing.T) {
		for range 100 {
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			req.RemoteAddr = "127.0.0.1:12345"
			w := httptest.NewRecorder()

			rateLimitHandler.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
		}
	})

	t.Run("block requests over limit", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.RemoteAddr = "127.0.0.1:12345"
		w := httptest.NewRecorder()

		rateLimitHandler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusTooManyRequests, w.Code)
	})
}

func TestAPIKeyMiddleware(t *testing.T) {
	// Set up environment
	os.Setenv("API_KEY", "test-api-key")
	defer os.Unsetenv("API_KEY")

	os.Setenv("JWT_SECRET", "test-jwt-secret")
	defer os.Unsetenv("JWT_SECRET")

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	apiKeyHandler := apiKeyMiddleware(testHandler)

	t.Run("valid API key", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "ApiKey test-api-key")
		w := httptest.NewRecorder()

		apiKeyHandler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("invalid API key", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "ApiKey invalid-key")
		w := httptest.NewRecorder()

		apiKeyHandler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("missing authorization header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		apiKeyHandler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("invalid authorization format", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "InvalidFormat")
		w := httptest.NewRecorder()

		apiKeyHandler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestAuthMiddleware(t *testing.T) {
	// AuthMiddleware is just an alias for APIKeyMiddleware, so test it the same way
	os.Setenv("API_KEY", "test-api-key")
	defer os.Unsetenv("API_KEY")

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	authHandler := AuthMiddleware(testHandler)

	t.Run("valid API key through AuthMiddleware", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "ApiKey test-api-key")
		w := httptest.NewRecorder()

		authHandler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
