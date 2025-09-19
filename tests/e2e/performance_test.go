package e2e

import (
	"context"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// LoadTestResult holds the results of a load test
type LoadTestResult struct {
	TotalRequests   int64
	SuccessfulReqs  int64
	FailedReqs      int64
	TotalDuration   time.Duration
	AvgResponseTime time.Duration
	MinResponseTime time.Duration
	MaxResponseTime time.Duration
	RequestsPerSec  float64
	ErrorRate       float64
}

// runLoadTest performs a load test on the given URL
func runLoadTest(ctx context.Context, t *testing.T, baseURL, endpoint string, numRequests, concurrency int) *LoadTestResult {
	t.Helper()

	url := baseURL + endpoint
	var wg sync.WaitGroup
	var mu sync.Mutex
	var results []time.Duration
	var successCount int64
	var failCount int64

	start := time.Now()

	sem := make(chan struct{}, concurrency)

	for range numRequests {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			reqStart := time.Now()
			req, err := MakeAuthenticatedRequest("GET", url, "")
			if err != nil {
				atomic.AddInt64(&failCount, 1)
				return
			}

			client := &http.Client{Timeout: 10 * time.Second}
			resp, err := client.Do(req)
			reqEnd := time.Now()
			duration := reqEnd.Sub(reqStart)

			mu.Lock()
			results = append(results, duration)
			mu.Unlock()

			if err != nil || resp.StatusCode >= 400 {
				atomic.AddInt64(&failCount, 1)
			} else {
				atomic.AddInt64(&successCount, 1)
			}
			if resp != nil {
				resp.Body.Close()
			}
		}()
	}

	wg.Wait()
	totalDuration := time.Since(start)

	// Calculate metrics
	var totalResponseTime time.Duration
	var minTime = time.Hour
	var maxTime time.Duration

	for _, d := range results {
		totalResponseTime += d
		if d < minTime {
			minTime = d
		}
		if d > maxTime {
			maxTime = d
		}
	}

	avgResponseTime := totalResponseTime / time.Duration(len(results))
	requestsPerSec := float64(numRequests) / totalDuration.Seconds()
	errorRate := float64(failCount) / float64(numRequests) * 100

	return &LoadTestResult{
		TotalRequests:   int64(numRequests),
		SuccessfulReqs:  successCount,
		FailedReqs:      failCount,
		TotalDuration:   totalDuration,
		AvgResponseTime: avgResponseTime,
		MinResponseTime: minTime,
		MaxResponseTime: maxTime,
		RequestsPerSec:  requestsPerSec,
		ErrorRate:       errorRate,
	}
}

func TestPerformance_RulesAPI_LoadTest(t *testing.T) {
	ctx := context.Background()
	env, cleanup := SetupTestEnvironment(ctx, t)
	defer cleanup()

	env.WaitForServices(ctx, t)

	baseURL := env.GetRuleEngineURL(ctx, t) + "/api/v1"

	// Run load test on /rules endpoint
	result := runLoadTest(ctx, t, baseURL, "/rules", 100, 10)

	// Log results
	t.Logf("Rules API Load Test Results:")
	t.Logf("Total Requests: %d", result.TotalRequests)
	t.Logf("Successful Requests: %d", result.SuccessfulReqs)
	t.Logf("Failed Requests: %d", result.FailedReqs)
	t.Logf("Total Duration: %v", result.TotalDuration)
	t.Logf("Avg Response Time: %v", result.AvgResponseTime)
	t.Logf("Min Response Time: %v", result.MinResponseTime)
	t.Logf("Max Response Time: %v", result.MaxResponseTime)
	t.Logf("Requests/sec: %.2f", result.RequestsPerSec)
	t.Logf("Error Rate: %.2f%%", result.ErrorRate)

	// Assertions
	require.Greater(t, result.SuccessfulReqs, int64(90), "Should have at least 90 successful requests")
	require.Less(t, result.ErrorRate, 10.0, "Error rate should be less than 10%")
	require.Less(t, result.AvgResponseTime, 500*time.Millisecond, "Average response time should be less than 500ms")
}

func TestPerformance_TriggersAPI_LoadTest(t *testing.T) {
	ctx := context.Background()
	env, cleanup := SetupTestEnvironment(ctx, t)
	defer cleanup()

	env.WaitForServices(ctx, t)

	baseURL := env.GetRuleEngineURL(ctx, t) + "/api/v1"

	// Run load test on /triggers endpoint
	result := runLoadTest(ctx, t, baseURL, "/triggers", 100, 10)

	// Log results
	t.Logf("Triggers API Load Test Results:")
	t.Logf("Total Requests: %d", result.TotalRequests)
	t.Logf("Successful Requests: %d", result.SuccessfulReqs)
	t.Logf("Failed Requests: %d", result.FailedReqs)
	t.Logf("Total Duration: %v", result.TotalDuration)
	t.Logf("Avg Response Time: %v", result.AvgResponseTime)
	t.Logf("Min Response Time: %v", result.MinResponseTime)
	t.Logf("Max Response Time: %v", result.MaxResponseTime)
	t.Logf("Requests/sec: %.2f", result.RequestsPerSec)
	t.Logf("Error Rate: %.2f%%", result.ErrorRate)

	// Assertions
	require.Greater(t, result.SuccessfulReqs, int64(90), "Should have at least 90 successful requests")
	require.Less(t, result.ErrorRate, 10.0, "Error rate should be less than 10%")
	require.Less(t, result.AvgResponseTime, 500*time.Millisecond, "Average response time should be less than 500ms")
}

func TestPerformance_ActionsAPI_LoadTest(t *testing.T) {
	ctx := context.Background()
	env, cleanup := SetupTestEnvironment(ctx, t)
	defer cleanup()

	env.WaitForServices(ctx, t)

	baseURL := env.GetRuleEngineURL(ctx, t) + "/api/v1"

	// Run load test on /actions endpoint
	result := runLoadTest(ctx, t, baseURL, "/actions", 100, 10)

	// Log results
	t.Logf("Actions API Load Test Results:")
	t.Logf("Total Requests: %d", result.TotalRequests)
	t.Logf("Successful Requests: %d", result.SuccessfulReqs)
	t.Logf("Failed Requests: %d", result.FailedReqs)
	t.Logf("Total Duration: %v", result.TotalDuration)
	t.Logf("Avg Response Time: %v", result.AvgResponseTime)
	t.Logf("Min Response Time: %v", result.MinResponseTime)
	t.Logf("Max Response Time: %v", result.MaxResponseTime)
	t.Logf("Requests/sec: %.2f", result.RequestsPerSec)
	t.Logf("Error Rate: %.2f%%", result.ErrorRate)

	// Assertions
	require.Greater(t, result.SuccessfulReqs, int64(90), "Should have at least 90 successful requests")
	require.Less(t, result.ErrorRate, 10.0, "Error rate should be less than 10%")
	require.Less(t, result.AvgResponseTime, 500*time.Millisecond, "Average response time should be less than 500ms")
}
