package manager

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/malyshevhen/rule-engine/internal/core/action"
	"github.com/malyshevhen/rule-engine/internal/core/rule"
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

func TestManager_executeRule(t *testing.T) {
	mockRuleSvc := &mockRuleService{}
	mockExec := &mockExecutor{}

	mgr := &Manager{
		ruleSvc:  mockRuleSvc,
		executor: mockExec,
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
		Output:  []interface{}{true},
		Error:   "",
	}
	mockExec.On("ExecuteScript", mock.Anything, expectedRule.LuaScript, mock.Anything).Return(ruleResult)

	// Mock action execution
	actionResult := &execPkg.ExecuteResult{
		Success: true,
		Output:  []interface{}{},
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
		ruleSvc:  mockRuleSvc,
		executor: mockExec,
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
		Output:  []interface{}{false},
		Error:   "",
	}
	mockExec.On("ExecuteScript", mock.Anything, expectedRule.LuaScript, mock.Anything).Return(ruleResult)

	mgr.executeRule(context.Background(), ruleID)

	mockRuleSvc.AssertExpectations(t)
	mockExec.AssertExpectations(t)
	// Ensure only rule script was executed, no actions
	mockExec.AssertNumberOfCalls(t, "ExecuteScript", 1)
}
