package manager

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/malyshevhen/rule-engine/internal/core/queue"
	"github.com/malyshevhen/rule-engine/internal/core/rule"
	"github.com/malyshevhen/rule-engine/internal/core/trigger"
	"github.com/malyshevhen/rule-engine/internal/engine/executor"
	execCtx "github.com/malyshevhen/rule-engine/internal/engine/executor/context"
	"github.com/malyshevhen/rule-engine/internal/metrics"
	"github.com/malyshevhen/rule-engine/pkg/tracing"
	"github.com/nats-io/nats.go"
	"github.com/robfig/cron/v3"
	"go.opentelemetry.io/otel/attribute"
)

// RuleService interface
type RuleService interface {
	GetByID(ctx context.Context, id uuid.UUID) (*rule.Rule, error)
}

// Executor interface
type Executor interface {
	GetContextService() *execCtx.Service
	ExecuteScript(ctx context.Context, script string, execCtx *execCtx.ExecutionContext) *executor.ExecuteResult
}

// TriggerService interface
type TriggerService interface {
	GetByID(ctx context.Context, id uuid.UUID) (*trigger.Trigger, error)
	GetEnabledConditionalTriggers(ctx context.Context) ([]*trigger.Trigger, error)
	GetEnabledScheduledTriggers(ctx context.Context) ([]*trigger.Trigger, error)
}

// TriggerEvaluator interface
type TriggerEvaluator interface {
	EvaluateTriggers(ctx context.Context, triggers []*trigger.Trigger, eventData map[string]any) []*trigger.EvaluationResult
}

// Manager handles trigger execution
type Manager struct {
	nc             *nats.Conn
	cron           *cron.Cron
	ruleSvc        RuleService
	triggerSvc     TriggerService
	triggerEval    TriggerEvaluator
	executor       Executor
	queue          queue.Queue
	executingRules map[uuid.UUID]bool // To detect cycles in rule chaining
}

// NewManager creates a new trigger manager
func NewManager(nc *nats.Conn, cron *cron.Cron, ruleSvc RuleService, triggerSvc TriggerService, triggerEval TriggerEvaluator, executor Executor, queue queue.Queue) *Manager {
	return &Manager{
		nc:             nc,
		cron:           cron,
		ruleSvc:        ruleSvc,
		triggerSvc:     triggerSvc,
		triggerEval:    triggerEval,
		executor:       executor,
		queue:          queue,
		executingRules: make(map[uuid.UUID]bool),
	}
}

// Start begins listening for triggers
func (m *Manager) Start(ctx context.Context) error {
	// Subscribe to events for conditional triggers
	_, err := m.nc.Subscribe("events.>", func(msg *nats.Msg) {
		m.handleConditionalTrigger(ctx, msg)
	})
	if err != nil {
		return err
	}

	// Load and schedule cron triggers
	return m.loadScheduledTriggers(ctx)
}

// Stop stops the trigger manager
func (m *Manager) Stop() {
	m.cron.Stop()
	m.nc.Close()
}

// handleConditionalTrigger processes incoming events for conditional triggers
func (m *Manager) handleConditionalTrigger(ctx context.Context, msg *nats.Msg) {
	// Record metric
	metrics.TriggerEventsTotal.WithLabelValues("conditional", "processed").Inc()

	// Parse event data
	var eventData map[string]any
	if err := json.Unmarshal(msg.Data, &eventData); err != nil {
		slog.Error("Failed to parse event data", "error", err)
		return
	}

	slog.Info("Received conditional trigger event", "subject", msg.Subject, "data", eventData)

	// Load enabled conditional triggers
	conditionalTriggers, err := m.triggerSvc.GetEnabledConditionalTriggers(ctx)
	if err != nil {
		slog.Error("Failed to load conditional triggers", "error", err)
		return
	}

	if len(conditionalTriggers) == 0 {
		slog.Debug("No enabled conditional triggers found")
		return
	}

	// Evaluate all conditional triggers against the event
	results := m.triggerEval.EvaluateTriggers(ctx, conditionalTriggers, eventData)

	// Execute rules for triggers that matched
	for _, result := range results {
		if result.Matched {
			slog.Info("Trigger condition matched, executing rule",
				"trigger_id", result.TriggerID,
				"rule_id", result.RuleID)

			// Record that trigger fired
			metrics.TriggerEventsTotal.WithLabelValues("conditional", "fired").Inc()

			// Execute the associated rule
			m.executeRuleInternal(ctx, result.RuleID, eventData, result.TriggerID, true)
		} else if result.Error != "" {
			slog.Error("Trigger evaluation failed",
				"trigger_id", result.TriggerID,
				"rule_id", result.RuleID,
				"error", result.Error)
		}
	}
}

// loadScheduledTriggers loads and schedules CRON triggers
func (m *Manager) loadScheduledTriggers(ctx context.Context) error {
	// Load all enabled scheduled triggers
	scheduledTriggers, err := m.triggerSvc.GetEnabledScheduledTriggers(ctx)
	if err != nil {
		return fmt.Errorf("failed to load scheduled triggers: %w", err)
	}

	slog.Info("Loading scheduled triggers", "count", len(scheduledTriggers))

	// Schedule each trigger
	for _, trigger := range scheduledTriggers {
		// The condition_script contains the CRON expression for scheduled triggers
		cronExpr := trigger.ConditionScript
		if cronExpr == "" {
			slog.Warn("Scheduled trigger has empty CRON expression, skipping", "trigger_id", trigger.ID)
			continue
		}

		// Add the CRON job
		_, err := m.cron.AddFunc(cronExpr, func() {
			m.handleScheduledTrigger(ctx, trigger.ID)
		})
		if err != nil {
			slog.Error("Failed to schedule CRON trigger",
				"trigger_id", trigger.ID,
				"cron_expr", cronExpr,
				"error", err)
			continue
		}

		slog.Info("Scheduled CRON trigger",
			"trigger_id", trigger.ID,
			"rule_id", trigger.RuleID,
			"cron_expr", cronExpr)
	}

	return nil
}

// handleScheduledTrigger processes scheduled triggers
func (m *Manager) handleScheduledTrigger(ctx context.Context, triggerID uuid.UUID) {
	// Record metric
	metrics.TriggerEventsTotal.WithLabelValues("scheduled", "processed").Inc()

	slog.Info("Executing scheduled trigger", "trigger_id", triggerID)

	// Get the trigger to find the associated rule
	trigger, err := m.triggerSvc.GetByID(ctx, triggerID)
	if err != nil {
		slog.Error("Failed to get scheduled trigger", "trigger_id", triggerID, "error", err)
		return
	}

	if !trigger.Enabled {
		slog.Warn("Scheduled trigger is disabled, skipping", "trigger_id", triggerID)
		return
	}

	// Record that trigger fired
	metrics.TriggerEventsTotal.WithLabelValues("scheduled", "fired").Inc()

	// Execute the associated rule
	m.executeRuleInternal(ctx, trigger.RuleID, nil, triggerID, true)
}

// executeRule executes a rule's logic (queues by default)
func (m *Manager) executeRule(ctx context.Context, ruleID uuid.UUID) {
	m.executeRuleInternal(ctx, ruleID, nil, uuid.Nil, true)
}

// executeRuleSync executes a rule synchronously (for rule chaining)
func (m *Manager) executeRuleSync(ctx context.Context, ruleID uuid.UUID) {
	m.executeRuleInternal(ctx, ruleID, nil, uuid.Nil, false)
}

// executeRuleInternal executes a rule's logic with queuing option
func (m *Manager) executeRuleInternal(ctx context.Context, ruleID uuid.UUID, eventData map[string]interface{}, triggerID uuid.UUID, allowQueue bool) {
	// If queuing is allowed and we have a queue, enqueue the request
	if allowQueue && m.queue != nil {
		req := &queue.ExecutionRequest{
			RuleID:    ruleID,
			TriggerID: triggerID,
			EventData: eventData,
		}

		if err := m.queue.Enqueue(ctx, req); err != nil {
			slog.Error("Failed to enqueue rule execution", "rule_id", ruleID, "error", err)
			// Fall back to synchronous execution
			m.executeRuleSynchronous(ctx, ruleID, eventData, triggerID)
		} else {
			slog.Info("Rule execution enqueued", "rule_id", ruleID, "trigger_id", triggerID)
		}
		return
	}

	// Execute synchronously
	m.executeRuleSynchronous(ctx, ruleID, eventData, triggerID)
}

// executeRuleSynchronous executes a rule synchronously
func (m *Manager) executeRuleSynchronous(ctx context.Context, ruleID uuid.UUID, eventData map[string]interface{}, triggerID uuid.UUID) {
	ctx, span := tracing.StartSpan(ctx, "manager.execute_rule_sync")
	defer span.End()

	span.SetAttributes(
		attribute.String("rule.id", ruleID.String()),
	)

	// Check for cycles in rule chaining
	if m.executingRules[ruleID] {
		slog.Warn("Cycle detected in rule execution, skipping", "rule_id", ruleID)
		return
	}
	m.executingRules[ruleID] = true
	defer func() { delete(m.executingRules, ruleID) }()

	rule, err := m.ruleSvc.GetByID(ctx, ruleID)
	if err != nil {
		span.RecordError(err)
		slog.Error("Failed to get rule", "rule_id", ruleID, "error", err)
		return
	}

	span.SetAttributes(
		attribute.String("rule.name", rule.Name),
	)

	// Create execution context
	execCtx := m.executor.GetContextService().CreateContext(ruleID.String(), triggerID.String())

	// Add event data to context if available
	if eventData != nil {
		for key, value := range eventData {
			execCtx.Data[key] = value
		}
	}

	// Execute rule script
	ruleCtx, ruleSpan := tracing.StartSpan(ctx, "rule.script_execution")
	result := m.executor.ExecuteScript(ruleCtx, rule.LuaScript, execCtx)
	ruleSpan.End()

	if result.Error != "" {
		span.RecordError(fmt.Errorf("rule execution failed: %s", result.Error))
		slog.Error("Rule execution failed", "rule_id", ruleID, "error", result.Error)
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
		slog.Info("Rule condition not met", "rule_id", ruleID)
		return
	}

	slog.Info("Rule executed successfully", "rule_id", ruleID)

	// Execute actions
	for _, action := range rule.Actions {
		actionCtx, actionSpan := tracing.StartSpan(ctx, "action.execution")
		actionSpan.SetAttributes(
			attribute.String("action.id", action.ID.String()),
			attribute.String("action.type", action.Type),
		)

		switch action.Type {
		case "lua_script":
			actionResult := m.executor.ExecuteScript(actionCtx, action.LuaScript, execCtx)
			if actionResult.Error != "" {
				actionSpan.RecordError(fmt.Errorf("lua action execution failed: %s", actionResult.Error))
				slog.Error("Lua action execution failed", "action_id", action.ID, "error", actionResult.Error)
			} else {
				slog.Info("Lua action executed", "action_id", action.ID)
			}
		case "execute_rule":
			var params map[string]interface{}
			if err := json.Unmarshal([]byte(action.Params), &params); err != nil {
				actionSpan.RecordError(fmt.Errorf("failed to parse execute_rule params: %w", err))
				slog.Error("Failed to parse execute_rule params", "action_id", action.ID, "error", err)
				continue
			}
			ruleIDStr, ok := params["rule_id"].(string)
			if !ok {
				actionSpan.RecordError(fmt.Errorf("invalid rule_id in execute_rule params"))
				slog.Error("Invalid rule_id in execute_rule params", "action_id", action.ID)
				continue
			}
			targetRuleID, err := uuid.Parse(ruleIDStr)
			if err != nil {
				actionSpan.RecordError(fmt.Errorf("invalid rule_id format: %w", err))
				slog.Error("Invalid rule_id format", "action_id", action.ID, "rule_id_str", ruleIDStr, "error", err)
				continue
			}
			slog.Info("Executing chained rule synchronously", "action_id", action.ID, "target_rule_id", targetRuleID)
			m.executeRuleSync(actionCtx, targetRuleID)
		default:
			actionSpan.RecordError(fmt.Errorf("unknown action type: %s", action.Type))
			slog.Error("Unknown action type", "action_id", action.ID, "type", action.Type)
		}

		actionSpan.End()
	}
}
