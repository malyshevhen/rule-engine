package trigger

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/malyshevhen/rule-engine/internal/engine/executor"
	execCtx "github.com/malyshevhen/rule-engine/internal/engine/executor/context"
	"github.com/malyshevhen/rule-engine/internal/metrics"
)

// Executor interface for script execution
type Executor interface {
	GetContextService() *execCtx.Service
	ExecuteScript(ctx context.Context, script string, execCtx *execCtx.ExecutionContext) *executor.ExecuteResult
}

// Evaluator handles trigger condition evaluation
type Evaluator struct {
	executor Executor
}

// NewEvaluator creates a new trigger evaluator
func NewEvaluator(executor Executor) *Evaluator {
	return &Evaluator{
		executor: executor,
	}
}

// EvaluationResult represents the result of evaluating a trigger condition
type EvaluationResult struct {
	TriggerID uuid.UUID
	RuleID    uuid.UUID
	Matched   bool
	Error     string
	Duration  time.Duration
}

// EvaluateCondition evaluates a trigger condition script against an event
func (e *Evaluator) EvaluateCondition(ctx context.Context, triggerID, ruleID uuid.UUID, conditionScript string, eventData map[string]any) *EvaluationResult {
	start := time.Now()

	// Record metric
	defer func() {
		duration := time.Since(start)
		metrics.TriggerEvaluationDuration.WithLabelValues("conditional").Observe(duration.Seconds())
	}()

	// Create execution context with event data
	execContext := e.executor.GetContextService().CreateContext(triggerID.String(), "")

	// Set event data in the context
	execContext.Data["event"] = eventData

	// Execute the condition script
	result := e.executor.ExecuteScript(ctx, conditionScript, execContext)

	duration := time.Since(start)

	// Check if the script executed successfully
	if result.Error != "" {
		slog.Error("Trigger condition evaluation failed",
			"trigger_id", triggerID,
			"rule_id", ruleID,
			"error", result.Error,
			"duration", duration)

		metrics.TriggerEvaluationTotal.WithLabelValues("conditional", "error").Inc()

		return &EvaluationResult{
			TriggerID: triggerID,
			RuleID:    ruleID,
			Matched:   false,
			Error:     result.Error,
			Duration:  duration,
		}
	}

	// Evaluate the result - expect a boolean return value
	matched := false
	if len(result.Output) > 0 {
		if b, ok := result.Output[0].(bool); ok {
			matched = b
		} else {
			// If not a boolean, treat non-nil/non-false values as true
			matched = result.Output[0] != nil && result.Output[0] != false
		}
	}

	slog.Debug("Trigger condition evaluated",
		"trigger_id", triggerID,
		"rule_id", ruleID,
		"matched", matched,
		"duration", duration)

	if matched {
		metrics.TriggerEvaluationTotal.WithLabelValues("conditional", "matched").Inc()
	} else {
		metrics.TriggerEvaluationTotal.WithLabelValues("conditional", "not_matched").Inc()
	}

	return &EvaluationResult{
		TriggerID: triggerID,
		RuleID:    ruleID,
		Matched:   matched,
		Error:     "",
		Duration:  duration,
	}
}

// EvaluateTriggers evaluates multiple triggers against an event
func (e *Evaluator) EvaluateTriggers(ctx context.Context, triggers []*Trigger, eventData map[string]any) []*EvaluationResult {
	results := make([]*EvaluationResult, 0, len(triggers))

	for _, trigger := range triggers {
		if trigger.Type != Conditional || !trigger.Enabled {
			continue
		}

		result := e.EvaluateCondition(ctx, trigger.ID, trigger.RuleID, trigger.ConditionScript, eventData)
		results = append(results, result)
	}

	return results
}
