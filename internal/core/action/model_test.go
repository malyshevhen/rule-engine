package action

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAction(t *testing.T) {
	id := uuid.New()
	now := time.Now()

	action := Action{
		ID:        id,
		LuaScript: "print('action executed')",
		Enabled:   true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	assert.Equal(t, id, action.ID)
	assert.Equal(t, "print('action executed')", action.LuaScript)
	assert.True(t, action.Enabled)
	assert.Equal(t, now, action.CreatedAt)
	assert.Equal(t, now, action.UpdatedAt)
}
