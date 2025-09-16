package context

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_CreateContext(t *testing.T) {
	service := NewService()

	ruleID := "rule-123"
	triggerID := "trigger-456"

	ctx := service.CreateContext(ruleID, triggerID)

	assert.NotNil(t, ctx)
	assert.Equal(t, ruleID, ctx.RuleID)
	assert.Equal(t, triggerID, ctx.TriggerID)
	assert.NotNil(t, ctx.Data)
	assert.Empty(t, ctx.Data)
}

func TestExecutionContext_Structure(t *testing.T) {
	ctx := &ExecutionContext{
		RuleID:    "rule-123",
		TriggerID: "trigger-456",
		Data:      make(map[string]interface{}),
	}

	// Test that we can add data to the context
	ctx.Data["device_id"] = "device-789"
	ctx.Data["temperature"] = 25.5

	assert.Equal(t, "device-789", ctx.Data["device_id"])
	assert.Equal(t, 25.5, ctx.Data["temperature"])
}
