package queue

import (
	"context"
	"fmt"
	"log/slog"
	"maps"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/malyshevhen/rule-engine/internal/core/rule"
	"github.com/malyshevhen/rule-engine/internal/engine/executor"
	execCtx "github.com/malyshevhen/rule-engine/internal/engine/executor/context"
	"github.com/malyshevhen/rule-engine/internal/metrics"
	"github.com/malyshevhen/rule-engine/pkg/tracing"
	"go.opentelemetry.io/otel/attribute"
)

// RuleService interface for rule operations
type RuleService interface {
	GetByID(ctx context.Context, id uuid.UUID) (*rule.Rule, error)
}

// Executor interface for script execution
type Executor interface {
	GetContextService() *execCtx.Service
	ExecuteScript(ctx context.Context, script string, execCtx *execCtx.ExecutionContext) *executor.ExecuteResult
}

// WorkerPool manages a pool of workers that process rule execution requests
type WorkerPool struct {
	queue      Queue
	ruleSvc    RuleService
	executor   Executor
	numWorkers int
	wg         sync.WaitGroup
	stopCh     chan struct{}
	stopped    bool
	mu         sync.RWMutex
	cleanupCh  chan struct{} // for cleanup goroutine
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(queue Queue, ruleSvc RuleService, executor Executor, numWorkers int) *WorkerPool {
	if numWorkers <= 0 {
		numWorkers = 5 // default
	}

	return &WorkerPool{
		queue:      queue,
		ruleSvc:    ruleSvc,
		executor:   executor,
		numWorkers: numWorkers,
		stopCh:     make(chan struct{}),
		cleanupCh:  make(chan struct{}),
	}
}

// Start begins processing the queue with the worker pool
func (wp *WorkerPool) Start(ctx context.Context) {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if wp.stopped {
		return
	}

	slog.Info("Starting rule execution worker pool", "num_workers", wp.numWorkers)

	for i := 0; i < wp.numWorkers; i++ {
		wp.wg.Add(1)
		go wp.worker(ctx, i)
	}

	// Start cleanup goroutine for Redis queues
	wp.wg.Add(1)
	go wp.cleanupWorker(ctx)
}

// Stop gracefully stops the worker pool
func (wp *WorkerPool) Stop() {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if wp.stopped {
		return
	}

	slog.Info("Stopping rule execution worker pool")
	close(wp.stopCh)
	close(wp.cleanupCh)
	wp.stopped = true
	wp.wg.Wait()
	slog.Info("Rule execution worker pool stopped")
}

// worker processes rule execution requests from the queue
func (wp *WorkerPool) worker(ctx context.Context, workerID int) {
	defer wp.wg.Done()

	slog.Info("Rule execution worker started", "worker_id", workerID)

	for {
		select {
		case <-wp.stopCh:
			slog.Info("Rule execution worker stopping", "worker_id", workerID)
			return
		default:
			// Try to dequeue a request
			req, err := wp.queue.Dequeue(ctx)
			if err != nil {
				if err == ErrQueueEmpty {
					// Queue is empty, wait a bit before trying again
					select {
					case <-wp.stopCh:
						return
					case <-ctx.Done():
						return
					default:
						// Small delay to avoid busy waiting
						continue
					}
				} else if err == ErrQueueClosed {
					slog.Info("Queue closed, worker stopping", "worker_id", workerID)
					return
				} else {
					slog.Error("Failed to dequeue request", "worker_id", workerID, "error", err)
					continue
				}
			}

			// Process the request
			wp.processRequest(ctx, req, workerID)
		}
	}
}

// processRequest executes a rule based on the queued request
func (wp *WorkerPool) processRequest(ctx context.Context, req *ExecutionRequest, workerID int) {
	startTime := time.Now()
	ctx, span := tracing.StartSpan(ctx, "worker.process_request")
	defer span.End()

	span.SetAttributes(
		attribute.String("request.id", req.ID.String()),
		attribute.String("rule.id", req.RuleID.String()),
		attribute.String("trigger.id", req.TriggerID.String()),
		attribute.Int("worker.id", workerID),
	)

	slog.Info("Processing rule execution request",
		"request_id", req.ID,
		"rule_id", req.RuleID,
		"trigger_id", req.TriggerID,
		"worker_id", workerID,
		"queue_time", time.Since(req.QueuedAt))

	// Get the rule
	rule, err := wp.ruleSvc.GetByID(ctx, req.RuleID)
	if err != nil {
		span.RecordError(err)
		slog.Error("Failed to get rule for execution",
			"request_id", req.ID,
			"rule_id", req.RuleID,
			"error", err)
		return
	}

	span.SetAttributes(
		attribute.String("rule.name", rule.Name),
	)

	// Create execution context
	execCtx := wp.executor.GetContextService().CreateContext(req.RuleID.String(), req.TriggerID.String())

	// Add event data to context if available
	if req.EventData != nil {
		maps.Copy(execCtx.Data, req.EventData)
	}

	// Execute rule script
	ruleCtx, ruleSpan := tracing.StartSpan(ctx, "rule.script_execution")
	result := wp.executor.ExecuteScript(ruleCtx, rule.LuaScript, execCtx)
	ruleSpan.End()

	if result.Error != "" {
		span.RecordError(fmt.Errorf("rule execution failed: %s", result.Error))
		slog.Error("Rule execution failed",
			"request_id", req.ID,
			"rule_id", req.RuleID,
			"error", result.Error)
		return
	}

	// Check if rule condition is met (assume script returns boolean)
	rulePassed := false
	if len(result.Output) > 0 {
		if b, ok := result.Output[0].(bool); ok {
			rulePassed = b
		}
	}

	span.SetAttributes(
		attribute.Bool("rule.condition_met", rulePassed),
	)

	if !rulePassed {
		slog.Info("Rule condition not met",
			"request_id", req.ID,
			"rule_id", req.RuleID)
		return
	}

	slog.Info("Rule executed successfully",
		"request_id", req.ID,
		"rule_id", req.RuleID)

	// Record execution metric
	metrics.RuleExecutionsTotal.WithLabelValues(rule.ID.String(), "success").Inc()

	// Execute actions
	for _, action := range rule.Actions {
		actionCtx, actionSpan := tracing.StartSpan(ctx, "action.execution")
		actionSpan.SetAttributes(
			attribute.String("action.id", action.ID.String()),
			attribute.String("action.type", action.Type),
		)

		switch action.Type {
		case "lua_script":
			actionResult := wp.executor.ExecuteScript(actionCtx, action.LuaScript, execCtx)
			if actionResult.Error != "" {
				actionSpan.RecordError(fmt.Errorf("lua action execution failed: %s", actionResult.Error))
				slog.Error("Lua action execution failed",
					"request_id", req.ID,
					"action_id", action.ID,
					"error", actionResult.Error)
			} else {
				slog.Info("Lua action executed",
					"request_id", req.ID,
					"action_id", action.ID)
			}
		case "execute_rule":
			// For now, skip rule chaining in queued executions to avoid complexity
			// This can be implemented later with proper cycle detection
			slog.Info("Skipping rule chaining in queued execution",
				"request_id", req.ID,
				"action_id", action.ID)
		default:
			actionSpan.RecordError(fmt.Errorf("unknown action type: %s", action.Type))
			slog.Error("Unknown action type",
				"request_id", req.ID,
				"action_id", action.ID,
				"type", action.Type)
		}

		actionSpan.End()
	}

	// Record processing duration metric
	processingDuration := time.Since(startTime)
	metrics.QueueProcessingDuration.Observe(processingDuration.Seconds())

	// Release the distributed lock if using Redis queue
	if redisQueue, ok := wp.queue.(*RedisQueue); ok {
		if err := redisQueue.ReleaseLock(ctx, req.ID.String()); err != nil {
			slog.Error("Failed to release distributed lock", "request_id", req.ID, "error", err)
		}
	}
}

// cleanupWorker periodically cleans up expired items from Redis queues and sends heartbeats
func (wp *WorkerPool) cleanupWorker(ctx context.Context) {
	defer wp.wg.Done()

	cleanupTicker := time.NewTicker(5 * time.Minute)    // cleanup every 5 minutes
	heartbeatTicker := time.NewTicker(30 * time.Second) // heartbeat every 30 seconds
	defer cleanupTicker.Stop()
	defer heartbeatTicker.Stop()

	for {
		select {
		case <-wp.cleanupCh:
			slog.Info("Queue cleanup worker stopping")
			return
		case <-cleanupTicker.C:
			wp.performCleanup(ctx)
		case <-heartbeatTicker.C:
			wp.sendHeartbeat(ctx)
		}
	}
}

// performCleanup cleans up expired queue items
func (wp *WorkerPool) performCleanup(ctx context.Context) {
	// Check if queue supports cleanup (Redis queue does, in-memory doesn't)
	if redisQueue, ok := wp.queue.(*RedisQueue); ok {
		if err := redisQueue.CleanupExpired(ctx, 24*time.Hour); err != nil {
			slog.Error("Failed to cleanup expired queue items", "error", err)
		}
		// Also cleanup stale locks from dead instances
		if err := redisQueue.CleanupStaleLocks(ctx); err != nil {
			slog.Error("Failed to cleanup stale locks", "error", err)
		}
	}
}

// sendHeartbeat sends a heartbeat for Redis queue instances
func (wp *WorkerPool) sendHeartbeat(ctx context.Context) {
	if redisQueue, ok := wp.queue.(*RedisQueue); ok {
		if err := redisQueue.SendHeartbeat(ctx); err != nil {
			slog.Error("Failed to send heartbeat", "error", err)
		}
	}
}
