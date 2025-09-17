package analytics

import (
	"context"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	service := NewService()
	assert.NotNil(t, service)
}

func TestService_GetDashboardData(t *testing.T) {
	service := NewService()

	tests := []struct {
		name      string
		timeRange string
		expected  string
	}{
		{"1h", "1h", "1h"},
		{"24h", "24h", "24h"},
		{"1d", "1d", "24h"}, // should normalize to 24h
		{"7d", "7d", "7d"},
		{"30d", "30d", "30d"},
		{"invalid", "invalid", "24h"}, // should default to 24h
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := service.GetDashboardData(context.Background(), tt.timeRange)
			assert.NoError(t, err)
			assert.NotNil(t, data)
			assert.Equal(t, tt.expected, data.TimeRange)
			assert.NotNil(t, data.OverallStats)
			assert.NotNil(t, data.RuleStats)
			assert.NotNil(t, data.ExecutionTrend)
			assert.NotNil(t, data.SuccessRateTrend)
			assert.NotNil(t, data.LatencyTrend)
		})
	}
}

func TestService_getOverallStats(t *testing.T) {
	// Create test metrics
	executionsCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rule_engine_rule_executions_total",
			Help: "Total number of rule executions",
		},
		[]string{"rule_id", "result"},
	)

	latencyHistogram := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "rule_engine_rule_execution_duration_seconds",
			Help:    "Duration of rule executions in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"rule_id"},
	)

	// Register metrics
	prometheus.MustRegister(executionsCounter)
	prometheus.MustRegister(latencyHistogram)
	defer prometheus.Unregister(executionsCounter)
	defer prometheus.Unregister(latencyHistogram)

	// Add test data
	executionsCounter.WithLabelValues("rule1", "success").Add(10)
	executionsCounter.WithLabelValues("rule1", "failure").Add(2)
	executionsCounter.WithLabelValues("rule2", "success").Add(5)

	latencyHistogram.WithLabelValues("rule1").Observe(0.1)  // 100ms
	latencyHistogram.WithLabelValues("rule1").Observe(0.2)  // 200ms
	latencyHistogram.WithLabelValues("rule2").Observe(0.05) // 50ms

	service := NewService()
	stats := service.getOverallStats()

	assert.Equal(t, int64(17), stats.TotalExecutions)
	assert.Equal(t, int64(15), stats.SuccessfulExecutions)
	assert.Equal(t, int64(2), stats.FailedExecutions)
	assert.InDelta(t, 88.235, stats.SuccessRate, 0.001)
	assert.InDelta(t, 116.666, stats.AverageLatency, 0.001) // (0.1+0.2+0.05)*1000/3
}

func TestService_getRuleStats(t *testing.T) {
	// Create test metrics
	executionsCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rule_engine_rule_executions_total",
			Help: "Total number of rule executions",
		},
		[]string{"rule_id", "result"},
	)

	// Register metrics
	prometheus.MustRegister(executionsCounter)
	defer prometheus.Unregister(executionsCounter)

	// Add test data
	executionsCounter.WithLabelValues("rule1", "success").Add(8)
	executionsCounter.WithLabelValues("rule1", "failure").Add(2)
	executionsCounter.WithLabelValues("rule2", "success").Add(5)

	service := NewService()
	ruleStats := service.getRuleStats()

	assert.Len(t, ruleStats, 2)

	// Find rule1 stats
	var rule1Stats, rule2Stats *RuleStats
	for i := range ruleStats {
		if ruleStats[i].RuleID == "rule1" {
			rule1Stats = &ruleStats[i]
		} else if ruleStats[i].RuleID == "rule2" {
			rule2Stats = &ruleStats[i]
		}
	}

	assert.NotNil(t, rule1Stats)
	assert.Equal(t, "rule1", rule1Stats.RuleID)
	assert.Equal(t, int64(10), rule1Stats.TotalExecutions)
	assert.Equal(t, int64(8), rule1Stats.SuccessfulExecutions)
	assert.Equal(t, int64(2), rule1Stats.FailedExecutions)
	assert.InDelta(t, 80.0, rule1Stats.SuccessRate, 0.001)

	assert.NotNil(t, rule2Stats)
	assert.Equal(t, "rule2", rule2Stats.RuleID)
	assert.Equal(t, int64(5), rule2Stats.TotalExecutions)
	assert.Equal(t, int64(5), rule2Stats.SuccessfulExecutions)
	assert.Equal(t, int64(0), rule2Stats.FailedExecutions)
	assert.InDelta(t, 100.0, rule2Stats.SuccessRate, 0.001)
}

func TestService_getExecutionTrend(t *testing.T) {
	service := NewService()
	startTime := time.Now().Add(-2 * time.Hour)

	trend := service.getExecutionTrend(startTime)

	assert.Equal(t, "executions_per_hour", trend.Metric)
	assert.True(t, len(trend.Data) > 0)

	for _, point := range trend.Data {
		assert.True(t, point.Timestamp.After(startTime) || point.Timestamp.Equal(startTime))
		assert.True(t, point.Value >= 10) // Mock data starts from 10
	}
}

func TestService_getSuccessRateTrend(t *testing.T) {
	service := NewService()
	startTime := time.Now().Add(-2 * time.Hour)

	trend := service.getSuccessRateTrend(startTime)

	assert.Equal(t, "success_rate_percent", trend.Metric)
	assert.True(t, len(trend.Data) > 0)

	for _, point := range trend.Data {
		assert.True(t, point.Timestamp.After(startTime) || point.Timestamp.Equal(startTime))
		assert.True(t, point.Value >= 85) // Mock data starts from 85
	}
}

func TestService_getLatencyTrend(t *testing.T) {
	service := NewService()
	startTime := time.Now().Add(-2 * time.Hour)

	trend := service.getLatencyTrend(startTime)

	assert.Equal(t, "average_latency_ms", trend.Metric)
	assert.True(t, len(trend.Data) > 0)

	for _, point := range trend.Data {
		assert.True(t, point.Timestamp.After(startTime) || point.Timestamp.Equal(startTime))
		assert.True(t, point.Value >= 50) // Mock data starts from 50
	}
}

func BenchmarkService_GetDashboardData(b *testing.B) {
	service := NewService()

	for b.Loop() {
		_, err := service.GetDashboardData(context.Background(), "24h")
		if err != nil {
			b.Fatal(err)
		}
	}
}
