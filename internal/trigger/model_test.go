package trigger

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTrigger(t *testing.T) {
	id := uuid.New()
	ruleID := uuid.New()
	now := time.Now()

	trig := Trigger{
		ID:              id,
		RuleID:          ruleID,
		Type:            Conditional,
		ConditionScript: "event.type == 'device.update'",
		Enabled:         true,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	assert.Equal(t, id, trig.ID)
	assert.Equal(t, ruleID, trig.RuleID)
	assert.Equal(t, Conditional, trig.Type)
	assert.Equal(t, "event.type == 'device.update'", trig.ConditionScript)
	assert.True(t, trig.Enabled)
	assert.Equal(t, now, trig.CreatedAt)
	assert.Equal(t, now, trig.UpdatedAt)
}

func TestTriggerType(t *testing.T) {
	assert.Equal(t, TriggerType("CONDITIONAL"), Conditional)
	assert.Equal(t, TriggerType("CRON"), Cron)
}
