package action

import (
	"time"

	"github.com/google/uuid"
)

// Action represents an action in the business domain
type Action struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Params    string    `json:"params"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// LuaScript kept for backward compatibility in API
	LuaScript string `json:"lua_script,omitempty"`
}
