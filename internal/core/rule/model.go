package rule

import (
	"time"

	"github.com/google/uuid"
	"github.com/malyshevhen/rule-engine/internal/core/action"
	"github.com/malyshevhen/rule-engine/internal/core/trigger"
)

// Rule represents a business rule with its triggers and actions
type Rule struct {
	ID        uuid.UUID         `json:"id"`
	Name      string            `json:"name"`
	LuaScript string            `json:"lua_script"`
	Priority  int               `json:"priority"`
	Enabled   bool              `json:"enabled"`
	Triggers  []trigger.Trigger `json:"triggers"`
	Actions   []action.Action   `json:"actions"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}
