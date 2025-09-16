//go:build performance

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/malyshevhen/rule-engine/internal/analytics"
	"github.com/malyshevhen/rule-engine/internal/core/action"
	"github.com/malyshevhen/rule-engine/internal/core/rule"
	"github.com/malyshevhen/rule-engine/internal/core/trigger"
	actionStorage "github.com/malyshevhen/rule-engine/internal/storage/action"
	"github.com/malyshevhen/rule-engine/internal/storage/db"
	ruleStorage "github.com/malyshevhen/rule-engine/internal/storage/rule"
	triggerStorage "github.com/malyshevhen/rule-engine/internal/storage/trigger"
	"github.com/stretchr/testify/require"
)

// PerformanceReport holds comprehensive performance metrics
type PerformanceReport struct {
	TestName        string
	TotalRequests   int
	SuccessfulReqs  int
	FailedReqs      int
	TotalDuration   time.Duration
	AvgResponseTime time.Duration
	MinResponseTime time.Duration
	MaxResponseTime time.Duration
	P50ResponseTime time.Duration
	P95ResponseTime time.Duration
	P99ResponseTime time.Duration
	RequestsPerSec  float64
	ErrorRate       float64
}

type PerformanceMetrics struct {
	TotalRequests   int
	SuccessfulReqs  int
	FailedReqs      int
	TotalDuration   time.Duration
	AvgResponseTime time.Duration
	MinResponseTime time.Duration
	MaxResponseTime time.Duration
	RequestsPerSec  float64
}

func setupPerformanceTest(t *testing.T) (*Server, func()) {
	// Skip if not running performance tests
	if testing.Short() {
		t.Skip("Skipping performance test")
	}

	// Reset middleware state for test isolation
	ResetMiddlewareForTesting()

	// Set test API key for authentication
	testAPIKey := "test-api-key-performance"
	originalAPIKey := os.Getenv("API_KEY")
	os.Setenv("API_KEY", testAPIKey)

	// Setup test containers
	ctx := context.Background()
	tc, cleanupContainers := SetupTestContainers(ctx, t)

	// Wait for services to be ready
	tc.WaitForServices(ctx, t)

	// Setup database pool
	pool := tc.SetupDatabasePool(ctx, t)

	// Run migrations
	err := db.RunMigrations(pool)
	require.NoError(t, err)

	// Setup Redis client
	redisClient := tc.GetRedisClient(ctx, t)

	// Create repositories
	ruleRepo := ruleStorage.NewRepository(pool)
	triggerRepo := triggerStorage.NewRepository(pool)
	actionRepo := actionStorage.NewRepository(pool)

	// Create services
	ruleSvc := rule.NewService(ruleRepo, triggerRepo, actionRepo, redisClient)
	triggerSvc := trigger.NewService(triggerRepo, redisClient)
	actionSvc := action.NewService(actionRepo)
	analyticsSvc := analytics.NewService()

	// Disable rate limiting for performance testing
	DisableRateLimiting()

	// Create server
	config := &ServerConfig{Port: "8080"}
	server := NewServer(config, ruleSvc, triggerSvc, actionSvc, analyticsSvc)

	// Return cleanup function
	cleanup := func() {
		EnableRateLimiting() // Re-enable rate limiting
		redisClient.Close()
		pool.Close()
		cleanupContainers()
		if originalAPIKey != "" {
			os.Setenv("API_KEY", originalAPIKey)
		} else {
			os.Unsetenv("API_KEY")
		}
	}

	return server, cleanup
}

func runLoadTest(t *testing.T, server *Server, endpoint string, method string, bodyFunc func(int) []byte, numRequests int, concurrency int) PerformanceReport {
	var wg sync.WaitGroup
	results := make(chan time.Duration, numRequests)
	errors := make(chan error, numRequests)

	start := time.Now()

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < numRequests/concurrency; j++ {
				reqID := workerID*(numRequests/concurrency) + j

				var body []byte
				if bodyFunc != nil {
					body = bodyFunc(reqID)
				}

				req := httptest.NewRequest(method, endpoint, bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "ApiKey test-api-key-performance")

				w := httptest.NewRecorder()
				reqStart := time.Now()
				server.router.ServeHTTP(w, req)
				reqEnd := time.Now()

				duration := reqEnd.Sub(reqStart)
				results <- duration

				if w.Code >= 400 {
					t.Logf("Request %d failed with status %d: %s", reqID, w.Code, w.Body.String())
					errors <- fmt.Errorf("request %d failed with status %d: %s", reqID, w.Code, w.Body.String())
				}
			}
		}(i)
	}

	wg.Wait()
	close(results)
	close(errors)

	totalDuration := time.Since(start)

	var durations []time.Duration
	errorCount := 0

	for d := range results {
		durations = append(durations, d)
	}

	for range errors {
		errorCount++
	}

	successCount := len(durations) - errorCount

	if len(durations) == 0 {
		return PerformanceReport{}
	}

	// Sort durations for percentile calculations
	sort.Slice(durations, func(i, j int) bool {
		return durations[i] < durations[j]
	})

	min := durations[0]
	max := durations[len(durations)-1]
	sum := time.Duration(0)

	for _, d := range durations {
		sum += d
	}

	avg := sum / time.Duration(len(durations))
	reqsPerSec := float64(len(durations)) / totalDuration.Seconds()
	errorRate := float64(errorCount) / float64(len(durations)+errorCount) * 100

	// Calculate percentiles
	p50 := durations[len(durations)/2]
	p95 := durations[int(float64(len(durations))*0.95)]
	p99 := durations[int(float64(len(durations))*0.99)]

	return PerformanceReport{
		TestName:        fmt.Sprintf("%s_%s", method, endpoint),
		TotalRequests:   len(durations),
		SuccessfulReqs:  successCount,
		FailedReqs:      errorCount,
		TotalDuration:   totalDuration,
		AvgResponseTime: avg,
		MinResponseTime: min,
		MaxResponseTime: max,
		P50ResponseTime: p50,
		P95ResponseTime: p95,
		P99ResponseTime: p99,
		RequestsPerSec:  reqsPerSec,
		ErrorRate:       errorRate,
	}
}

func TestPerformance_RulesAPI_LoadTest(t *testing.T) {
	server, cleanup := setupPerformanceTest(t)
	defer cleanup()

	numRequests := 100
	concurrency := 10

	t.Run("CreateRule_LoadTest", func(t *testing.T) {
		bodyFunc := func(reqID int) []byte {
			req := CreateRuleRequest{
				Name:      fmt.Sprintf("Load Test Rule %d", reqID),
				LuaScript: "return event.temperature > 25",
				Priority:  &[]int{0}[0],
				Enabled:   &[]bool{true}[0],
			}
			body, _ := json.Marshal(req)
			return body
		}

		report := runLoadTest(t, server, "/api/v1/rules", "POST", bodyFunc, numRequests, concurrency)

		t.Logf("Create Rule Load Test Results:")
		t.Logf("Test Name: %s", report.TestName)
		t.Logf("Total Requests: %d", report.TotalRequests)
		t.Logf("Successful Requests: %d", report.SuccessfulReqs)
		t.Logf("Failed Requests: %d", report.FailedReqs)
		t.Logf("Total Duration: %v", report.TotalDuration)
		t.Logf("Avg Response Time: %v", report.AvgResponseTime)
		t.Logf("Min Response Time: %v", report.MinResponseTime)
		t.Logf("Max Response Time: %v", report.MaxResponseTime)
		t.Logf("P50 Response Time: %v", report.P50ResponseTime)
		t.Logf("P95 Response Time: %v", report.P95ResponseTime)
		t.Logf("P99 Response Time: %v", report.P99ResponseTime)
		t.Logf("Requests/sec: %.2f", report.RequestsPerSec)
		t.Logf("Error Rate: %.2f%%", report.ErrorRate)

		// Without rate limiting, all requests should succeed
		require.Equal(t, numRequests, report.SuccessfulReqs, "All requests should succeed")
		require.Less(t, report.AvgResponseTime, 100*time.Millisecond, "Average response time should be under 100ms")
	})

	t.Run("GetRules_LoadTest", func(t *testing.T) {
		report := runLoadTest(t, server, "/api/v1/rules", "GET", nil, numRequests, concurrency)

		t.Logf("Get Rules Load Test Results:")
		t.Logf("Test Name: %s", report.TestName)
		t.Logf("Total Requests: %d", report.TotalRequests)
		t.Logf("Successful Requests: %d", report.SuccessfulReqs)
		t.Logf("Failed Requests: %d", report.FailedReqs)
		t.Logf("Total Duration: %v", report.TotalDuration)
		t.Logf("Avg Response Time: %v", report.AvgResponseTime)
		t.Logf("Min Response Time: %v", report.MinResponseTime)
		t.Logf("Max Response Time: %v", report.MaxResponseTime)
		t.Logf("P50 Response Time: %v", report.P50ResponseTime)
		t.Logf("P95 Response Time: %v", report.P95ResponseTime)
		t.Logf("P99 Response Time: %v", report.P99ResponseTime)
		t.Logf("Requests/sec: %.2f", report.RequestsPerSec)
		t.Logf("Error Rate: %.2f%%", report.ErrorRate)

		// Without rate limiting, all requests should succeed
		require.Equal(t, numRequests, report.SuccessfulReqs, "All requests should succeed")
		require.Less(t, report.AvgResponseTime, 50*time.Millisecond, "Average response time should be under 50ms")
	})
}

func TestPerformance_TriggersAPI_LoadTest(t *testing.T) {
	server, cleanup := setupPerformanceTest(t)
	defer cleanup()

	// First create a rule to use for triggers
	ruleReq := CreateRuleRequest{
		Name:      "Performance Test Rule",
		LuaScript: "return event.temperature > 25",
		Priority:  &[]int{0}[0],
		Enabled:   &[]bool{true}[0],
	}
	ruleBody, _ := json.Marshal(ruleReq)
	ruleHttpReq := httptest.NewRequest("POST", "/api/v1/rules", bytes.NewReader(ruleBody))
	ruleHttpReq.Header.Set("Content-Type", "application/json")
	ruleHttpReq.Header.Set("Authorization", "ApiKey test-api-key-performance")
	ruleW := httptest.NewRecorder()
	server.router.ServeHTTP(ruleW, ruleHttpReq)
	require.Equal(t, http.StatusOK, ruleW.Code)

	var ruleResp map[string]interface{}
	json.Unmarshal(ruleW.Body.Bytes(), &ruleResp)
	ruleIDStr := ruleResp["id"].(string)
	ruleID, _ := uuid.Parse(ruleIDStr)

	numRequests := 100
	concurrency := 10

	t.Run("CreateTrigger_LoadTest", func(t *testing.T) {
		bodyFunc := func(reqID int) []byte {
			req := CreateTriggerRequest{
				RuleID:          ruleID,
				Type:            "CONDITIONAL",
				ConditionScript: "return event.temperature > 25",
				Enabled:         &[]bool{true}[0],
			}
			body, _ := json.Marshal(req)
			return body
		}

		report := runLoadTest(t, server, "/api/v1/triggers", "POST", bodyFunc, numRequests, concurrency)

		t.Logf("Create Trigger Load Test Results:")
		t.Logf("Test Name: %s", report.TestName)
		t.Logf("Total Requests: %d", report.TotalRequests)
		t.Logf("Successful Requests: %d", report.SuccessfulReqs)
		t.Logf("Failed Requests: %d", report.FailedReqs)
		t.Logf("Total Duration: %v", report.TotalDuration)
		t.Logf("Avg Response Time: %v", report.AvgResponseTime)
		t.Logf("Min Response Time: %v", report.MinResponseTime)
		t.Logf("Max Response Time: %v", report.MaxResponseTime)
		t.Logf("P50 Response Time: %v", report.P50ResponseTime)
		t.Logf("P95 Response Time: %v", report.P95ResponseTime)
		t.Logf("P99 Response Time: %v", report.P99ResponseTime)
		t.Logf("Requests/sec: %.2f", report.RequestsPerSec)
		t.Logf("Error Rate: %.2f%%", report.ErrorRate)

		// Without rate limiting, all requests should succeed
		require.Equal(t, numRequests, report.SuccessfulReqs, "All requests should succeed")
		require.Less(t, report.AvgResponseTime, 100*time.Millisecond, "Average response time should be under 100ms")
	})
}

func TestPerformance_ActionsAPI_LoadTest(t *testing.T) {
	server, cleanup := setupPerformanceTest(t)
	defer cleanup()

	numRequests := 100
	concurrency := 10

	t.Run("CreateAction_LoadTest", func(t *testing.T) {
		bodyFunc := func(reqID int) []byte {
			req := CreateActionRequest{
				LuaScript: "platform.send_command('device_1', 'turn_on')",
				Enabled:   &[]bool{true}[0],
			}
			body, _ := json.Marshal(req)
			return body
		}

		report := runLoadTest(t, server, "/api/v1/actions", "POST", bodyFunc, numRequests, concurrency)

		t.Logf("Create Action Load Test Results:")
		t.Logf("Test Name: %s", report.TestName)
		t.Logf("Total Requests: %d", report.TotalRequests)
		t.Logf("Successful Requests: %d", report.SuccessfulReqs)
		t.Logf("Failed Requests: %d", report.FailedReqs)
		t.Logf("Total Duration: %v", report.TotalDuration)
		t.Logf("Avg Response Time: %v", report.AvgResponseTime)
		t.Logf("Min Response Time: %v", report.MinResponseTime)
		t.Logf("Max Response Time: %v", report.MaxResponseTime)
		t.Logf("P50 Response Time: %v", report.P50ResponseTime)
		t.Logf("P95 Response Time: %v", report.P95ResponseTime)
		t.Logf("P99 Response Time: %v", report.P99ResponseTime)
		t.Logf("Requests/sec: %.2f", report.RequestsPerSec)
		t.Logf("Error Rate: %.2f%%", report.ErrorRate)

		// Without rate limiting, all requests should succeed
		require.Equal(t, numRequests, report.SuccessfulReqs, "All requests should succeed")
		require.Less(t, report.AvgResponseTime, 100*time.Millisecond, "Average response time should be under 100ms")
	})
}

// BenchmarkDatabaseQueries benchmarks database query performance
func BenchmarkDatabaseQueries(b *testing.B) {
	// Setup test database

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:password@localhost:5433/rule_engine?sslmode=disable"
	}

	ctx := context.Background()
	pool, err := db.NewPostgresPool(ctx, dbURL)
	if err != nil {
		b.Fatal(err)
	}
	defer pool.Close()

	err = db.RunMigrations(pool)
	if err != nil {
		b.Fatal(err)
	}

	repo := ruleStorage.NewRepository(pool)

	// Pre-populate with test data
	for i := 0; i < 100; i++ {
		rule := &ruleStorage.Rule{
			Name:      fmt.Sprintf("Benchmark Rule %d", i),
			LuaScript: "return event.temperature > 25",
			Enabled:   true,
		}
		err := repo.Create(ctx, rule)
		if err != nil {
			b.Fatal(err)
		}
	}

	b.ResetTimer()

	b.Run("GetByID", func(b *testing.B) {
		ruleID := uuid.New() // This will fail, but we want to measure the query time
		for i := 0; i < b.N; i++ {
			_, _ = repo.GetByID(ctx, ruleID)
		}
	})

	b.Run("List", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = repo.List(ctx, 50, 0)
		}
	})

	b.Run("GetByIDWithAssociations", func(b *testing.B) {
		ruleID := uuid.New() // This will fail, but we want to measure the query time
		for i := 0; i < b.N; i++ {
			_, _, _, _ = repo.GetByIDWithAssociations(ctx, ruleID)
		}
	})
}
