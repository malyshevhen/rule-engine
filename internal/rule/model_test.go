package rule

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRule(t *testing.T) {
	id := uuid.New()
	now := time.Now()

	rule := Rule{
		ID:        id,
		Name:      "Test Rule",
		LuaScript: "return true",
		Enabled:   true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	assert.Equal(t, id, rule.ID)
	assert.Equal(t, "Test Rule", rule.Name)
	assert.Equal(t, "return true", rule.LuaScript)
	assert.True(t, rule.Enabled)
	assert.Equal(t, now, rule.CreatedAt)
	assert.Equal(t, now, rule.UpdatedAt)
	assert.Empty(t, rule.Triggers)
	assert.Empty(t, rule.Actions)
}
