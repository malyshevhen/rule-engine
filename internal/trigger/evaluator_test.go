package trigger

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/malyshevhen/rule-engine/internal/engine/executor"
	execCtx "github.com/malyshevhen/rule-engine/internal/engine/executor/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockExecutor is a mock implementation of the Executor interface
type mockExecutor struct {
	mock.Mock
}

func (m *mockExecutor) GetContextService() *execCtx.Service {
	args := m.Called()
	return args.Get(0).(*execCtx.Service)
}

func (m *mockExecutor) ExecuteScript(ctx context.Context, script string, execCtx *execCtx.ExecutionContext) *executor.ExecuteResult {
	args := m.Called(ctx, script, execCtx)
	return args.Get(0).(*executor.ExecuteResult)
}

func TestEvaluator_EvaluateCondition(t *testing.T) {
	tests := []struct {
		name            string
		conditionScript string
		eventData       map[string]any
		mockResult      *executor.ExecuteResult
		expectedMatch   bool
		expectError     bool
	}{
		{
			name:            "condition matches",
			conditionScript: "return event.temperature > 25",
			eventData:       map[string]any{"temperature": 30.0},
			mockResult: &executor.ExecuteResult{
				Success: true,
				Output:  []any{true},
			},
			expectedMatch: true,
			expectError:   false,
		},
		{
			name:            "condition does not match",
			conditionScript: "return event.temperature > 25",
			eventData:       map[string]any{"temperature": 20.0},
			mockResult: &executor.ExecuteResult{
				Success: true,
				Output:  []any{false},
			},
			expectedMatch: false,
			expectError:   false,
		},
		{
			name:            "condition with complex logic",
			conditionScript: "return event.device_id == 'sensor_1' and event.value > 10",
			eventData:       map[string]any{"device_id": "sensor_1", "value": 15.0},
			mockResult: &executor.ExecuteResult{
				Success: true,
				Output:  []any{true},
			},
			expectedMatch: true,
			expectError:   false,
		},
		{
			name:            "execution error",
			conditionScript: "return invalid_syntax",
			eventData:       map[string]any{"temperature": 25.0},
			mockResult: &executor.ExecuteResult{
				Success: false,
				Error:   "syntax error",
				Output:  nil,
			},
			expectedMatch: false,
			expectError:   true,
		},
		{
			name:            "non-boolean result treated as true",
			conditionScript: "return 'matched'",
			eventData:       map[string]any{"status": "active"},
			mockResult: &executor.ExecuteResult{
				Success: true,
				Output:  []any{"matched"},
			},
			expectedMatch: true,
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockExec := &mockExecutor{}
			evaluator := NewEvaluator(mockExec)

			triggerID := uuid.New()
			ruleID := uuid.New()

			// Mock the context service
			contextSvc := execCtx.NewService()
			mockExec.On("GetContextService").Return(contextSvc)

			// Mock the script execution
			mockExec.On("ExecuteScript", mock.Anything, tt.conditionScript, mock.MatchedBy(func(ctx *execCtx.ExecutionContext) bool {
				return ctx.Data["event"] != nil
			})).Return(tt.mockResult)

			// Execute
			result := evaluator.EvaluateCondition(context.Background(), triggerID, ruleID, tt.conditionScript, tt.eventData)

			// Assert
			assert.Equal(t, triggerID, result.TriggerID)
			assert.Equal(t, ruleID, result.RuleID)
			assert.Equal(t, tt.expectedMatch, result.Matched)
			if tt.expectError {
				assert.NotEmpty(t, result.Error)
			} else {
				assert.Empty(t, result.Error)
			}

			mockExec.AssertExpectations(t)
		})
	}
}

func TestEvaluator_EvaluateTriggers(t *testing.T) {
	// Setup
	mockExec := &mockExecutor{}
	evaluator := NewEvaluator(mockExec)

	contextSvc := execCtx.NewService()
	mockExec.On("GetContextService").Return(contextSvc)

	triggerID1 := uuid.New()
	ruleID1 := uuid.New()
	triggerID2 := uuid.New()
	ruleID2 := uuid.New()

	triggers := []*Trigger{
		{
			ID:              triggerID1,
			RuleID:          ruleID1,
			Type:            Conditional,
			ConditionScript: "return event.temp > 25",
			Enabled:         true,
		},
		{
			ID:              triggerID2,
			RuleID:          ruleID2,
			Type:            Conditional,
			ConditionScript: "return event.status == 'active'",
			Enabled:         false, // Disabled trigger should not be evaluated
		},
	}

	eventData := map[string]any{"temp": 30.0, "status": "active"}

	// Mock executions - only the enabled trigger should be evaluated
	mockExec.On("ExecuteScript", mock.Anything, "return event.temp > 25", mock.Anything).Return(&executor.ExecuteResult{
		Success: true,
		Output:  []any{true},
	})

	// Execute
	results := evaluator.EvaluateTriggers(context.Background(), triggers, eventData)

	// Assert
	assert.Len(t, results, 1) // Only the enabled trigger should be evaluated
	assert.Equal(t, triggerID1, results[0].TriggerID)
	assert.Equal(t, ruleID1, results[0].RuleID)
	assert.True(t, results[0].Matched)

	mockExec.AssertExpectations(t)
}

func TestEvaluator_EvaluateCondition_EventDataInContext(t *testing.T) {
	// Test that event data is properly passed to the execution context
	mockExec := &mockExecutor{}
	evaluator := NewEvaluator(mockExec)

	contextSvc := execCtx.NewService()
	mockExec.On("GetContextService").Return(contextSvc)

	var capturedContext *execCtx.ExecutionContext
	mockExec.On("ExecuteScript", mock.Anything, mock.Anything, mock.MatchedBy(func(ctx *execCtx.ExecutionContext) bool {
		capturedContext = ctx
		return true
	})).Return(&executor.ExecuteResult{
		Success: true,
		Output:  []any{true},
	})

	eventData := map[string]any{"temperature": 28.5, "device": "sensor_1"}
	result := evaluator.EvaluateCondition(context.Background(), uuid.New(), uuid.New(), "return true", eventData)

	assert.True(t, result.Matched)
	assert.NotNil(t, capturedContext)
	assert.Equal(t, eventData, capturedContext.Data["event"])

	mockExec.AssertExpectations(t)
}
