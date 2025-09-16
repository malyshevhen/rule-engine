package api

import (
	"time"

	"github.com/google/uuid"
)

// RuleDTO represents a rule for API responses
type RuleDTO struct {
	ID        uuid.UUID    `json:"id"`
	Name      string       `json:"name"`
	LuaScript string       `json:"lua_script"`
	Enabled   bool         `json:"enabled"`
	Triggers  []TriggerDTO `json:"triggers"`
	Actions   []ActionDTO  `json:"actions"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

// TriggerDTO represents a trigger for API responses
type TriggerDTO struct {
	ID              uuid.UUID `json:"id"`
	Type            string    `json:"type"`
	ConditionScript string    `json:"condition_script"`
	Enabled         bool      `json:"enabled"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ActionDTO represents an action for API responses
type ActionDTO struct {
	ID        uuid.UUID `json:"id"`
	LuaScript string    `json:"lua_script"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateRuleRequest represents a request to create a rule
type CreateRuleRequest struct {
	Name      string `json:"name" validate:"required" example:"Temperature Alert Rule"`
	LuaScript string `json:"lua_script" validate:"required" example:"if event.temperature > 25 then return true end"`
	Enabled   *bool  `json:"enabled,omitempty" example:"true"`
}

// UpdateRuleRequest represents a request to update a rule
type UpdateRuleRequest struct {
	Name      *string `json:"name,omitempty" example:"Updated Rule Name"`
	LuaScript *string `json:"lua_script,omitempty" example:"if event.temperature > 30 then return true end"`
	Enabled   *bool   `json:"enabled,omitempty" example:"false"`
}

// CreateTriggerRequest represents a request to create a trigger
type CreateTriggerRequest struct {
	RuleID          uuid.UUID `json:"rule_id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Type            string    `json:"type" validate:"required,oneof=CONDITIONAL CRON" example:"CONDITIONAL"`
	ConditionScript string    `json:"condition_script" validate:"required" example:"if event.device_id == 'sensor_1' then return true end"`
	Enabled         *bool     `json:"enabled,omitempty" example:"true"`
}

// CreateActionRequest represents a request to create an action
type CreateActionRequest struct {
	LuaScript string `json:"lua_script" validate:"required" example:"log_message('info', 'Temperature alert triggered')"`
	Enabled   *bool  `json:"enabled,omitempty" example:"true"`
}

// APIErrorResponse represents an error response for API documentation
type APIErrorResponse struct {
	Error string `json:"error" example:"Error message"`
}
