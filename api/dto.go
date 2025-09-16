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
	Name      string `json:"name" validate:"required"`
	LuaScript string `json:"lua_script" validate:"required"`
	Enabled   *bool  `json:"enabled,omitempty"`
}

// UpdateRuleRequest represents a request to update a rule
type UpdateRuleRequest struct {
	Name      *string `json:"name,omitempty"`
	LuaScript *string `json:"lua_script,omitempty"`
	Enabled   *bool   `json:"enabled,omitempty"`
}

// CreateTriggerRequest represents a request to create a trigger
type CreateTriggerRequest struct {
	RuleID          uuid.UUID `json:"rule_id" validate:"required"`
	Type            string    `json:"type" validate:"required,oneof=CONDITIONAL CRON"`
	ConditionScript string    `json:"condition_script" validate:"required"`
	Enabled         *bool     `json:"enabled,omitempty"`
}

// CreateActionRequest represents a request to create an action
type CreateActionRequest struct {
	LuaScript string `json:"lua_script" validate:"required"`
	Enabled   *bool  `json:"enabled,omitempty"`
}
