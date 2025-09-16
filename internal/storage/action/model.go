package action

import (
	"time"

	"github.com/google/uuid"
)

// Action represents an action in the storage layer
type Action struct {
	ID        uuid.UUID `json:"id" db:"id"`
	LuaScript string    `json:"lua_script" db:"lua_script"`
	Enabled   bool      `json:"enabled" db:"enabled"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// TODO: add action model
