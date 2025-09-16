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

// Trigger represents a trigger in the storage layer
type Trigger struct {
	ID              uuid.UUID   `json:"id" db:"id"`
	RuleID          uuid.UUID   `json:"rule_id" db:"rule_id"`
	Type            TriggerType `json:"type" db:"type"`
	ConditionScript string      `json:"condition_script" db:"condition_script"`
	Enabled         bool        `json:"enabled" db:"enabled"`
	CreatedAt       time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at" db:"updated_at"`
}
