package trigger

import (
	"time"

	"github.com/google/uuid"
)

// TriggerType represents the type of trigger
type TriggerType string

const (
	Conditional TriggerType = "CONDITIONAL"
	Cron        TriggerType = "CRON"
)

// Trigger represents a trigger in the business domain
type Trigger struct {
	ID              uuid.UUID   `json:"id"`
	Type            TriggerType `json:"type"`
	ConditionScript string      `json:"condition_script"`
	Enabled         bool        `json:"enabled"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}
