package api

import (
	"time"

	"github.com/google/uuid"
)

// RuleInfo represents a rule for API responses
type RuleInfo struct {
	ID        uuid.UUID     `json:"id"`
	Name      string        `json:"name"`
	LuaScript string        `json:"lua_script"`
	Priority  int           `json:"priority"`
	Enabled   bool          `json:"enabled"`
	Triggers  []TriggerInfo `json:"triggers"`
	Actions   []ActionInfo  `json:"actions"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

// TriggerInfo represents a trigger for API responses
type TriggerInfo struct {
	ID              uuid.UUID `json:"id"`
	RuleID          uuid.UUID `json:"rule_id"`
	Type            string    `json:"type"`
	ConditionScript string    `json:"condition_script"`
	Enabled         bool      `json:"enabled"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ActionInfo represents an action for API responses
type ActionInfo struct {
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
	Priority  *int   `json:"priority,omitempty" example:"0"`
	Enabled   *bool  `json:"enabled,omitempty" example:"true"`
}

// UpdateRuleRequest represents a request to update a rule
type UpdateRuleRequest struct {
	Name      *string `json:"name,omitempty" example:"Updated Rule Name"`
	LuaScript *string `json:"lua_script,omitempty" example:"if event.temperature > 30 then return true end"`
	Priority  *int    `json:"priority,omitempty" example:"5"`
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

// EvaluateScriptRequest represents a request to evaluate a Lua script
type EvaluateScriptRequest struct {
	Script  string         `json:"script" validate:"required" example:"return 2 + 2"`
	Context map[string]any `json:"context,omitempty"`
}

// EvaluateScriptResponse represents the result of script evaluation
type EvaluateScriptResponse struct {
	Success  bool   `json:"success" example:"true"`
	Result   string `json:"result,omitempty" example:"4"`
	Output   []any  `json:"output,omitempty"`
	Error    string `json:"error,omitempty" example:"syntax error"`
	Duration string `json:"duration" example:"1.5ms"`
}

// AddActionToRuleRequest represents a request to add an action to a rule
type AddActionToRuleRequest struct {
	ActionID uuid.UUID `json:"action_id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// APIErrorResponse represents an error response for API documentation
type APIErrorResponse struct {
	Error string `json:"error" example:"Error message"`
}
