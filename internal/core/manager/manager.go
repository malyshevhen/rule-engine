package manager

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/google/uuid"
	"github.com/malyshevhen/rule-engine/internal/core/rule"
	"github.com/malyshevhen/rule-engine/internal/engine/executor"
	"github.com/nats-io/nats.go"
	"github.com/robfig/cron/v3"
)

// Manager handles trigger execution
type Manager struct {
	nc       *nats.Conn
	cron     *cron.Cron
	ruleSvc  *rule.Service
	executor *executor.Service
}

// NewManager creates a new trigger manager
func NewManager(nc *nats.Conn, cron *cron.Cron, ruleSvc *rule.Service, executor *executor.Service) *Manager {
	return &Manager{
		nc:       nc,
		cron:     cron,
		ruleSvc:  ruleSvc,
		executor: executor,
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
	// Parse event data
	var eventData map[string]interface{}
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
	slog.Info("Executing scheduled trigger", "trigger_id", triggerID)

	// TODO: Get trigger, get rule ID, execute rule
	dummyRuleID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	m.executeRule(ctx, dummyRuleID)
}

// executeRule executes a rule's logic
func (m *Manager) executeRule(ctx context.Context, ruleID uuid.UUID) {
	rule, err := m.ruleSvc.GetByID(ctx, ruleID)
	if err != nil {
		slog.Error("Failed to get rule", "rule_id", ruleID, "error", err)
		return
	}

	// Create execution context
	execCtx := m.executor.GetContextService().CreateContext(ruleID.String(), "")

	// Execute rule script
	result := m.executor.ExecuteScript(ctx, rule.LuaScript, execCtx)
	if result.Error != "" {
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

	if !rulePassed {
		slog.Info("Rule condition not met", "rule_id", ruleID)
		return
	}

	slog.Info("Rule executed successfully", "rule_id", ruleID)

	// Execute actions
	for _, action := range rule.Actions {
		actionResult := m.executor.ExecuteScript(ctx, action.LuaScript, execCtx)
		if actionResult.Error != "" {
			slog.Error("Action execution failed", "action_id", action.ID, "error", actionResult.Error)
		} else {
			slog.Info("Action executed", "action_id", action.ID)
		}
	}
}
