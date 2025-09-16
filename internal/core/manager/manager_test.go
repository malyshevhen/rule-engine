package manager

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/malyshevhen/rule-engine/internal/core/action"
	"github.com/malyshevhen/rule-engine/internal/core/rule"
	"github.com/malyshevhen/rule-engine/internal/core/trigger"
	execPkg "github.com/malyshevhen/rule-engine/internal/engine/executor"
	ctxPkg "github.com/malyshevhen/rule-engine/internal/engine/executor/context"
	"github.com/stretchr/testify/mock"
)

// mockRuleService is a mock implementation of RuleService
type mockRuleService struct {
	mock.Mock
}

func (m *mockRuleService) GetByID(ctx context.Context, id uuid.UUID) (*rule.Rule, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*rule.Rule), args.Error(1)
}

// mockExecutor is a mock implementation of Executor
type mockExecutor struct {
	mock.Mock
}

func (m *mockExecutor) GetContextService() *ctxPkg.Service {
	args := m.Called()
	return args.Get(0).(*ctxPkg.Service)
}

func (m *mockExecutor) ExecuteScript(ctx context.Context, script string, executionCtx *ctxPkg.ExecutionContext) *execPkg.ExecuteResult {
	args := m.Called(ctx, script, executionCtx)
	return args.Get(0).(*execPkg.ExecuteResult)
}

// mockTriggerService is a mock implementation of TriggerService
type mockTriggerService struct {
	mock.Mock
}

func (m *mockTriggerService) GetByID(ctx context.Context, id uuid.UUID) (*trigger.Trigger, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*trigger.Trigger), args.Error(1)
}

func (m *mockTriggerService) GetEnabledConditionalTriggers(ctx context.Context) ([]*trigger.Trigger, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*trigger.Trigger), args.Error(1)
}

func (m *mockTriggerService) GetEnabledScheduledTriggers(ctx context.Context) ([]*trigger.Trigger, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*trigger.Trigger), args.Error(1)
}

func TestManager_executeRule(t *testing.T) {
	mockRuleSvc := &mockRuleService{}
	mockExec := &mockExecutor{}

	mgr := &Manager{
		ruleSvc:        mockRuleSvc,
		executor:       mockExec,
		executingRules: make(map[uuid.UUID]bool),
	}

	ruleID := uuid.New()
	expectedRule := &rule.Rule{
		ID:        ruleID,
		Name:      "Test Rule",
		LuaScript: "return true",
		Enabled:   true,
		Actions: []action.Action{
			{
				ID:        uuid.New(),
				Type:      "lua_script",
				Params:    "print('action')",
				LuaScript: "print('action')",
				Enabled:   true,
			},
		},
	}

	mockRuleSvc.On("GetByID", mock.Anything, ruleID).Return(expectedRule, nil)

	mockExec.On("GetContextService").Return(ctxPkg.NewService())

	// Mock rule execution - returns true
	ruleResult := &execPkg.ExecuteResult{
		Success: true,
		Output:  []any{true},
		Error:   "",
	}
	mockExec.On("ExecuteScript", mock.Anything, expectedRule.LuaScript, mock.Anything).Return(ruleResult)

	// Mock action execution
	actionResult := &execPkg.ExecuteResult{
		Success: true,
		Output:  []any{},
		Error:   "",
	}
	mockExec.On("ExecuteScript", mock.Anything, expectedRule.Actions[0].LuaScript, mock.Anything).Return(actionResult)

	mgr.executeRule(context.Background(), ruleID)

	mockRuleSvc.AssertExpectations(t)
	mockExec.AssertExpectations(t)
}

func TestManager_executeRule_FailedCondition(t *testing.T) {
	mockRuleSvc := &mockRuleService{}
	mockExec := &mockExecutor{}

	mgr := &Manager{
		ruleSvc:        mockRuleSvc,
		executor:       mockExec,
		executingRules: make(map[uuid.UUID]bool),
	}

	ruleID := uuid.New()
	expectedRule := &rule.Rule{
		ID:        ruleID,
		Name:      "Test Rule",
		LuaScript: "return false",
		Enabled:   true,
		Actions:   []action.Action{}, // No actions
	}

	mockRuleSvc.On("GetByID", mock.Anything, ruleID).Return(expectedRule, nil)

	mockExec.On("GetContextService").Return(ctxPkg.NewService())

	// Mock rule execution - returns false
	ruleResult := &execPkg.ExecuteResult{
		Success: true,
		Output:  []any{false},
		Error:   "",
	}
	mockExec.On("ExecuteScript", mock.Anything, expectedRule.LuaScript, mock.Anything).Return(ruleResult)

	mgr.executeRule(context.Background(), ruleID)

	mockRuleSvc.AssertExpectations(t)
	mockExec.AssertExpectations(t)
	// Ensure only rule script was executed, no actions
	mockExec.AssertNumberOfCalls(t, "ExecuteScript", 1)
}

func TestManager_executeRule_RuleNotFound(t *testing.T) {
	mockRuleSvc := &mockRuleService{}
	mockExec := &mockExecutor{}

	mgr := &Manager{
		ruleSvc:        mockRuleSvc,
		executor:       mockExec,
		executingRules: make(map[uuid.UUID]bool),
	}

	ruleID := uuid.New()

	mockRuleSvc.On("GetByID", mock.Anything, ruleID).Return((*rule.Rule)(nil), errors.New("rule not found"))

	mgr.executeRule(context.Background(), ruleID)

	mockRuleSvc.AssertExpectations(t)
	mockExec.AssertNotCalled(t, "ExecuteScript")
}

func TestManager_executeRule_ExecutionError(t *testing.T) {
	mockRuleSvc := &mockRuleService{}
	mockExec := &mockExecutor{}

	mgr := &Manager{
		ruleSvc:        mockRuleSvc,
		executor:       mockExec,
		executingRules: make(map[uuid.UUID]bool),
	}

	ruleID := uuid.New()
	expectedRule := &rule.Rule{
		ID:        ruleID,
		Name:      "Test Rule",
		LuaScript: "error('test error')",
		Enabled:   true,
		Actions:   []action.Action{},
	}

	mockRuleSvc.On("GetByID", mock.Anything, ruleID).Return(expectedRule, nil)

	mockExec.On("GetContextService").Return(ctxPkg.NewService())

	// Mock rule execution - returns error
	ruleResult := &execPkg.ExecuteResult{
		Success: false,
		Output:  []any{},
		Error:   "test error",
	}
	mockExec.On("ExecuteScript", mock.Anything, expectedRule.LuaScript, mock.Anything).Return(ruleResult)

	mgr.executeRule(context.Background(), ruleID)

	mockRuleSvc.AssertExpectations(t)
	mockExec.AssertExpectations(t)
	mockExec.AssertNumberOfCalls(t, "ExecuteScript", 1)
}

func TestManager_executeRule_MultipleActions(t *testing.T) {
	mockRuleSvc := &mockRuleService{}
	mockExec := &mockExecutor{}

	mgr := &Manager{
		ruleSvc:        mockRuleSvc,
		executor:       mockExec,
		executingRules: make(map[uuid.UUID]bool),
	}

	ruleID := uuid.New()
	expectedRule := &rule.Rule{
		ID:        ruleID,
		Name:      "Test Rule",
		LuaScript: "return true",
		Enabled:   true,
		Actions: []action.Action{
			{
				ID:        uuid.New(),
				Type:      "lua_script",
				Params:    "print('action1')",
				LuaScript: "print('action1')",
				Enabled:   true,
			},
			{
				ID:        uuid.New(),
				Type:      "lua_script",
				Params:    "print('action2')",
				LuaScript: "print('action2')",
				Enabled:   true,
			},
		},
	}

	mockRuleSvc.On("GetByID", mock.Anything, ruleID).Return(expectedRule, nil)

	mockExec.On("GetContextService").Return(ctxPkg.NewService())

	// Mock rule execution - returns true
	ruleResult := &execPkg.ExecuteResult{
		Success: true,
		Output:  []any{true},
		Error:   "",
	}
	mockExec.On("ExecuteScript", mock.Anything, expectedRule.LuaScript, mock.Anything).Return(ruleResult)

	// Mock action executions
	actionResult := &execPkg.ExecuteResult{
		Success: true,
		Output:  []any{},
		Error:   "",
	}
	mockExec.On("ExecuteScript", mock.Anything, expectedRule.Actions[0].LuaScript, mock.Anything).Return(actionResult)
	mockExec.On("ExecuteScript", mock.Anything, expectedRule.Actions[1].LuaScript, mock.Anything).Return(actionResult)

	mgr.executeRule(context.Background(), ruleID)

	mockRuleSvc.AssertExpectations(t)
	mockExec.AssertExpectations(t)
	mockExec.AssertNumberOfCalls(t, "ExecuteScript", 3) // rule + 2 actions
}

func TestManager_handleScheduledTrigger(t *testing.T) {
	mockRuleSvc := &mockRuleService{}
	mockTriggerSvc := &mockTriggerService{}
	mockExec := &mockExecutor{}

	mgr := &Manager{
		ruleSvc:        mockRuleSvc,
		triggerSvc:     mockTriggerSvc,
		executor:       mockExec,
		executingRules: make(map[uuid.UUID]bool),
	}

	triggerID := uuid.New()
	ruleID := uuid.New()

	// Mock trigger
	expectedTrigger := &trigger.Trigger{
		ID:              triggerID,
		RuleID:          ruleID,
		Type:            trigger.Cron,
		ConditionScript: "@every 1m",
		Enabled:         true,
	}

	// Mock rule
	expectedRule := &rule.Rule{
		ID:        ruleID,
		Name:      "Scheduled Rule",
		LuaScript: "return true",
		Enabled:   true,
		Actions:   []action.Action{},
	}

	mockTriggerSvc.On("GetByID", mock.Anything, triggerID).Return(expectedTrigger, nil)
	mockRuleSvc.On("GetByID", mock.Anything, ruleID).Return(expectedRule, nil)

	mockExec.On("GetContextService").Return(ctxPkg.NewService())

	// Mock rule execution - returns true
	ruleResult := &execPkg.ExecuteResult{
		Success: true,
		Output:  []any{true},
		Error:   "",
	}
	mockExec.On("ExecuteScript", mock.Anything, expectedRule.LuaScript, mock.Anything).Return(ruleResult)

	mgr.handleScheduledTrigger(context.Background(), triggerID)

	mockTriggerSvc.AssertExpectations(t)
	mockRuleSvc.AssertExpectations(t)
	mockExec.AssertExpectations(t)
}

func TestManager_handleScheduledTrigger_TriggerNotFound(t *testing.T) {
	mockTriggerSvc := &mockTriggerService{}

	mgr := &Manager{
		triggerSvc: mockTriggerSvc,
	}

	triggerID := uuid.New()

	mockTriggerSvc.On("GetByID", mock.Anything, triggerID).Return((*trigger.Trigger)(nil), errors.New("trigger not found"))

	mgr.handleScheduledTrigger(context.Background(), triggerID)

	mockTriggerSvc.AssertExpectations(t)
}

func TestManager_handleScheduledTrigger_DisabledTrigger(t *testing.T) {
	mockTriggerSvc := &mockTriggerService{}

	mgr := &Manager{
		triggerSvc: mockTriggerSvc,
	}

	triggerID := uuid.New()

	// Mock disabled trigger
	disabledTrigger := &trigger.Trigger{
		ID:              triggerID,
		RuleID:          uuid.New(),
		Type:            trigger.Cron,
		ConditionScript: "@every 1m",
		Enabled:         false, // Disabled
	}

	mockTriggerSvc.On("GetByID", mock.Anything, triggerID).Return(disabledTrigger, nil)

	mgr.handleScheduledTrigger(context.Background(), triggerID)

	mockTriggerSvc.AssertExpectations(t)
	// Should not attempt to get or execute the rule
}
