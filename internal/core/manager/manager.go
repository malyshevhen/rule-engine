package manager

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/malyshevhen/rule-engine/internal/core/rule"
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

// Manager handles trigger execution
type Manager struct {
	nc             *nats.Conn
	cron           *cron.Cron
	ruleSvc        RuleService
	executor       Executor
	executingRules map[uuid.UUID]bool // To detect cycles in rule chaining
}

// NewManager creates a new trigger manager
func NewManager(nc *nats.Conn, cron *cron.Cron, ruleSvc RuleService, executor Executor) *Manager {
	return &Manager{
		nc:             nc,
		cron:           cron,
		ruleSvc:        ruleSvc,
		executor:       executor,
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

	// TODO: Load conditional triggers and evaluate conditions
	// For now, stub - assume a trigger fires
	slog.Info("Received conditional trigger event", "subject", msg.Subject, "data", eventData)

	// Dummy: execute a rule
	dummyRuleID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	m.executeRule(ctx, dummyRuleID)
}

// loadScheduledTriggers loads and schedules CRON triggers
func (m *Manager) loadScheduledTriggers(ctx context.Context) error {
	// TODO: Load all enabled CRON triggers from DB
	// For each, parse condition_script as cron expression, add job

	// Example: add a test cron job
	_, err := m.cron.AddFunc("@every 1m", func() {
		m.handleScheduledTrigger(ctx, uuid.New()) // dummy trigger ID
	})
	return err
}

// handleScheduledTrigger processes scheduled triggers
func (m *Manager) handleScheduledTrigger(ctx context.Context, triggerID uuid.UUID) {
	// Record metric
	metrics.TriggerEventsTotal.WithLabelValues("scheduled", "processed").Inc()

	slog.Info("Executing scheduled trigger", "trigger_id", triggerID)

	// TODO: Get trigger, get rule ID, execute rule
	dummyRuleID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	m.executeRule(ctx, dummyRuleID)
}

// executeRule executes a rule's logic
func (m *Manager) executeRule(ctx context.Context, ruleID uuid.UUID) {
	ctx, span := tracing.StartSpan(ctx, "manager.execute_rule")
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
	execCtx := m.executor.GetContextService().CreateContext(ruleID.String(), "")

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
			slog.Info("Executing chained rule", "action_id", action.ID, "target_rule_id", targetRuleID)
			m.executeRule(actionCtx, targetRuleID)
		default:
			actionSpan.RecordError(fmt.Errorf("unknown action type: %s", action.Type))
			slog.Error("Unknown action type", "action_id", action.ID, "type", action.Type)
		}

		actionSpan.End()
	}
}
