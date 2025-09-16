package action

import (
	"time"

	"github.com/google/uuid"
)

// Action represents an action in the storage layer
type Action struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Type      string    `json:"type" db:"type"`
	Params    string    `json:"params" db:"params"` // JSON string for parameters
	Enabled   bool      `json:"enabled" db:"enabled"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
