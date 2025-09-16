package rule

import (
	"time"

	"github.com/google/uuid"
)

// Rule represents a rule in the storage layer
type Rule struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	LuaScript string    `json:"lua_script" db:"lua_script"`
	Enabled   bool      `json:"enabled" db:"enabled"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
