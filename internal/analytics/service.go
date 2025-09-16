package analytics

import (
	"context"
	"sort"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Service provides analytics and reporting functionality
type Service struct{}

// NewService creates a new analytics service
func NewService() *Service {
	return &Service{}
}

// ExecutionStats represents execution statistics for a time period
type ExecutionStats struct {
	TotalExecutions      int64   `json:"total_executions"`
	SuccessfulExecutions int64   `json:"successful_executions"`
	FailedExecutions     int64   `json:"failed_executions"`
	SuccessRate          float64 `json:"success_rate"`
	AverageLatency       float64 `json:"average_latency_ms"`
}

// RuleStats represents statistics for a specific rule
type RuleStats struct {
	RuleID               string     `json:"rule_id"`
	RuleName             string     `json:"rule_name"`
	TotalExecutions      int64      `json:"total_executions"`
	SuccessfulExecutions int64      `json:"successful_executions"`
	FailedExecutions     int64      `json:"failed_executions"`
	SuccessRate          float64    `json:"success_rate"`
	AverageLatency       float64    `json:"average_latency_ms"`
	LastExecuted         *time.Time `json:"last_executed,omitempty"`
}

// TimeSeriesPoint represents a data point in a time series
type TimeSeriesPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
}

// TimeSeriesData represents time series data for a metric
type TimeSeriesData struct {
	Metric string            `json:"metric"`
	Data   []TimeSeriesPoint `json:"data"`
}

// DashboardData represents all data needed for the dashboard
type DashboardData struct {
	OverallStats     ExecutionStats `json:"overall_stats"`
	RuleStats        []RuleStats    `json:"rule_stats"`
	ExecutionTrend   TimeSeriesData `json:"execution_trend"`
	SuccessRateTrend TimeSeriesData `json:"success_rate_trend"`
	LatencyTrend     TimeSeriesData `json:"latency_trend"`
	TimeRange        string         `json:"time_range"`
}

// GetDashboardData returns aggregated data for the analytics dashboard
func (s *Service) GetDashboardData(ctx context.Context, timeRange string) (*DashboardData, error) {
	// Parse time range
	var startTime time.Time
	switch timeRange {
	case "1h":
		startTime = time.Now().Add(-time.Hour)
		timeRange = "1h"
	case "24h", "1d":
		startTime = time.Now().Add(-24 * time.Hour)
		timeRange = "24h"
	case "7d":
		startTime = time.Now().Add(-7 * 24 * time.Hour)
		timeRange = "7d"
	case "30d":
		startTime = time.Now().Add(-30 * 24 * time.Hour)
		timeRange = "30d"
	default:
		startTime = time.Now().Add(-24 * time.Hour) // default to 24h
		timeRange = "24h"
	}

	// Get overall stats from local metrics
	overallStats := s.getOverallStats()

	// Get rule-specific stats
	ruleStats := s.getRuleStats()

	// Get time series data (simplified - in real implementation would query Prometheus)
	executionTrend := s.getExecutionTrend(startTime)
	successRateTrend := s.getSuccessRateTrend(startTime)
	latencyTrend := s.getLatencyTrend(startTime)

	return &DashboardData{
		OverallStats:     overallStats,
		RuleStats:        ruleStats,
		ExecutionTrend:   executionTrend,
		SuccessRateTrend: successRateTrend,
		LatencyTrend:     latencyTrend,
		TimeRange:        timeRange,
	}, nil
}

// getOverallStats aggregates overall execution statistics
func (s *Service) getOverallStats() ExecutionStats {
	// Get metrics from Prometheus registry
	metricsFamilies, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		return ExecutionStats{}
	}

	var totalExecutions, successfulExecutions int64
	var totalLatency float64
	var latencyCount int64

	for _, mf := range metricsFamilies {
		switch *mf.Name {
		case "rule_engine_rule_executions_total":
			for _, m := range mf.Metric {
				value := int64(*m.Counter.Value)
				totalExecutions += value

				// Check result label
				result := "unknown"
				for _, label := range m.Label {
					if *label.Name == "result" {
						result = *label.Value
						break
					}
				}

				if result == "success" {
					successfulExecutions += value
				}
			}
		case "rule_engine_rule_execution_duration_seconds":
			for _, m := range mf.Metric {
				totalLatency += *m.Histogram.SampleSum
				latencyCount += int64(*m.Histogram.SampleCount)
			}
		}
	}

	failedExecutions := totalExecutions - successfulExecutions
	successRate := float64(0)
	if totalExecutions > 0 {
		successRate = float64(successfulExecutions) / float64(totalExecutions) * 100
	}

	averageLatency := float64(0)
	if latencyCount > 0 {
		averageLatency = totalLatency / float64(latencyCount) * 1000 // Convert to milliseconds
	}

	return ExecutionStats{
		TotalExecutions:      totalExecutions,
		SuccessfulExecutions: successfulExecutions,
		FailedExecutions:     failedExecutions,
		SuccessRate:          successRate,
		AverageLatency:       averageLatency,
	}
}

// getRuleStats returns statistics for each rule
func (s *Service) getRuleStats() []RuleStats {
	metricsFamilies, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		return []RuleStats{}
	}

	ruleMap := make(map[string]*RuleStats)

	for _, mf := range metricsFamilies {
		switch *mf.Name {
		case "rule_engine_rule_executions_total":
			for _, m := range mf.Metric {
				var ruleID, result string
				value := int64(*m.Counter.Value)

				for _, label := range m.Label {
					switch *label.Name {
					case "rule_id":
						ruleID = *label.Value
					case "result":
						result = *label.Value
					}
				}

				if ruleID == "" {
					continue
				}

				if ruleMap[ruleID] == nil {
					ruleMap[ruleID] = &RuleStats{RuleID: ruleID}
				}

				ruleMap[ruleID].TotalExecutions += value
				if result == "success" {
					ruleMap[ruleID].SuccessfulExecutions += value
				}
			}
		}
	}

	// Calculate derived metrics and convert to slice
	result := make([]RuleStats, 0)
	for _, stats := range ruleMap {
		stats.FailedExecutions = stats.TotalExecutions - stats.SuccessfulExecutions
		if stats.TotalExecutions > 0 {
			stats.SuccessRate = float64(stats.SuccessfulExecutions) / float64(stats.TotalExecutions) * 100
		}
		result = append(result, *stats)
	}

	// Sort by total executions (descending)
	sort.Slice(result, func(i, j int) bool {
		return result[i].TotalExecutions > result[j].TotalExecutions
	})

	return result
}

// getExecutionTrend returns time series data for execution counts
func (s *Service) getExecutionTrend(startTime time.Time) TimeSeriesData {
	// For now, return mock data - in a real implementation, this would query Prometheus
	// for historical data or use a time-series database
	now := time.Now()
	data := []TimeSeriesPoint{}

	// Generate hourly data points for the last 24 hours
	for t := startTime; t.Before(now); t = t.Add(time.Hour) {
		data = append(data, TimeSeriesPoint{
			Timestamp: t,
			Value:     float64(10 + (t.Hour() % 5)), // Mock data
		})
	}

	return TimeSeriesData{
		Metric: "executions_per_hour",
		Data:   data,
	}
}

// getSuccessRateTrend returns time series data for success rates
func (s *Service) getSuccessRateTrend(startTime time.Time) TimeSeriesData {
	now := time.Now()
	data := []TimeSeriesPoint{}

	for t := startTime; t.Before(now); t = t.Add(time.Hour) {
		data = append(data, TimeSeriesPoint{
			Timestamp: t,
			Value:     85 + float64(t.Hour()%10), // Mock success rate between 85-95%
		})
	}

	return TimeSeriesData{
		Metric: "success_rate_percent",
		Data:   data,
	}
}

// getLatencyTrend returns time series data for latency
func (s *Service) getLatencyTrend(startTime time.Time) TimeSeriesData {
	now := time.Now()
	data := []TimeSeriesPoint{}

	for t := startTime; t.Before(now); t = t.Add(time.Hour) {
		data = append(data, TimeSeriesPoint{
			Timestamp: t,
			Value:     50 + float64(t.Hour()%20), // Mock latency in milliseconds
		})
	}

	return TimeSeriesData{
		Metric: "average_latency_ms",
		Data:   data,
	}
}
