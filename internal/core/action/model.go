package action

import (
	"time"

	"github.com/google/uuid"
)

// Action represents an action in the business domain
type Action struct {
	ID        uuid.UUID `json:"id"`
	LuaScript string    `json:"lua_script"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
