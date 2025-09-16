package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// RuleExecutionsTotal counts total rule executions by rule ID and result
	RuleExecutionsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rule_engine_rule_executions_total",
			Help: "Total number of rule executions",
		},
		[]string{"rule_id", "result"}, // result: success, failure
	)

	// TriggerEventsTotal counts trigger events processed
	TriggerEventsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rule_engine_trigger_events_total",
			Help: "Total number of trigger events processed",
		},
		[]string{"trigger_type", "action"}, // action: processed, fired
	)

	// LuaExecutionErrorsTotal counts Lua execution errors
	LuaExecutionErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rule_engine_lua_execution_errors_total",
			Help: "Total number of Lua execution errors",
		},
		[]string{"rule_id", "error_type"},
	)

	// RuleExecutionDuration measures execution duration
	RuleExecutionDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "rule_engine_rule_execution_duration_seconds",
			Help:    "Duration of rule executions in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"rule_id"},
	)

	// TriggerEvaluationTotal counts trigger condition evaluations
	TriggerEvaluationTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rule_engine_trigger_evaluation_total",
			Help: "Total number of trigger condition evaluations",
		},
		[]string{"trigger_type", "result"}, // result: matched, not_matched, error
	)

	// TriggerEvaluationDuration measures trigger evaluation duration
	TriggerEvaluationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "rule_engine_trigger_evaluation_duration_seconds",
			Help:    "Duration of trigger condition evaluations in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"trigger_type"},
	)

	// QueueSize measures the current queue size
	QueueSize = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "rule_engine_queue_size",
			Help: "Current number of items in the execution queue",
		},
	)

	// QueueEnqueueTotal counts enqueued execution requests
	QueueEnqueueTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "rule_engine_queue_enqueue_total",
			Help: "Total number of execution requests enqueued",
		},
	)

	// QueueDequeueTotal counts dequeued execution requests
	QueueDequeueTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "rule_engine_queue_dequeue_total",
			Help: "Total number of execution requests dequeued",
		},
	)

	// QueueProcessingDuration measures queue processing duration
	QueueProcessingDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "rule_engine_queue_processing_duration_seconds",
			Help:    "Duration of queue processing in seconds",
			Buckets: prometheus.DefBuckets,
		},
	)

	// AlertsTotal counts total alerts sent
	AlertsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rule_engine_alerts_total",
			Help: "Total number of alerts sent",
		},
		[]string{"alert_type", "severity"},
	)

	// AlertSendErrorsTotal counts alert sending errors
	AlertSendErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rule_engine_alert_send_errors_total",
			Help: "Total number of alert sending errors",
		},
		[]string{"alert_type"},
	)
)
