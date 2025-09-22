package client

import (
	"time"

	"github.com/google/uuid"
)

// RuleInfo represents a rule in the system
type RuleInfo struct {
	ID        uuid.UUID     `json:"id"`
	Name      string        `json:"name"`
	LuaScript string        `json:"lua_script"`
	Priority  int           `json:"priority"`
	Enabled   bool          `json:"enabled"`
	Triggers  []TriggerInfo `json:"triggers,omitempty"`
	Actions   []ActionInfo  `json:"actions,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

// TriggerInfo represents a trigger in the system
type TriggerInfo struct {
	ID              uuid.UUID `json:"id"`
	RuleID          uuid.UUID `json:"rule_id"`
	Type            string    `json:"type"` // CONDITIONAL or CRON
	ConditionScript string    `json:"condition_script"`
	Enabled         bool      `json:"enabled"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ActionInfo represents an action in the system
type ActionInfo struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	LuaScript string    `json:"lua_script"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateRuleRequest represents a request to create a rule
type CreateRuleRequest struct {
	Name      string `json:"name"`
	LuaScript string `json:"lua_script"`
	Priority  *int   `json:"priority,omitempty"`
	Enabled   *bool  `json:"enabled,omitempty"`
}

// CreateTriggerRequest represents a request to create a trigger
type CreateTriggerRequest struct {
	RuleID          uuid.UUID `json:"rule_id"`
	Type            string    `json:"type"` // CONDITIONAL or CRON
	ConditionScript string    `json:"condition_script"`
	Enabled         *bool     `json:"enabled,omitempty"`
}

// CreateActionRequest represents a request to create an action
type CreateActionRequest struct {
	Name      string `json:"name,omitempty"`
	LuaScript string `json:"lua_script"`
	Enabled   *bool  `json:"enabled,omitempty"`
}

// EvaluateScriptRequest represents a request to evaluate a Lua script
type EvaluateScriptRequest struct {
	Script  string         `json:"script"`
	Context map[string]any `json:"context,omitempty"`
}

// EvaluateScriptResponse represents the response from script evaluation
type EvaluateScriptResponse struct {
	Success  bool   `json:"success"`
	Result   any    `json:"result,omitempty"`
	Output   []any  `json:"output,omitempty"`
	Error    string `json:"error,omitempty"`
	Duration string `json:"duration"`
}

// AddActionToRuleRequest represents a request to add an action to a rule
type AddActionToRuleRequest struct {
	ActionID uuid.UUID `json:"action_id"`
}

// PatchOperation represents a JSON Patch operation
type PatchOperation struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value any    `json:"value,omitempty"`
}

// PatchRequest represents a JSON Patch request
type PatchRequest []PatchOperation

// UpdateActionRequest represents a request to update an action
type UpdateActionRequest struct {
	Patches PatchRequest `json:"-"` // Not serialized, used for JSON Patch
}

// UpdateTriggerRequest represents a request to update a trigger
type UpdateTriggerRequest struct {
	Patches PatchRequest `json:"-"` // Not serialized, used for JSON Patch
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Database string `json:"database"` // "ok", "error"
	Redis    string `json:"redis"`    // "ok", "error", "disabled"
}

// PaginatedRulesResponse represents a paginated list of rules
type PaginatedRulesResponse struct {
	Rules  []RuleInfo `json:"rules"`
	Limit  int        `json:"limit"`
	Offset int        `json:"offset"`
	Count  int        `json:"count"`
	Total  int        `json:"total"`
}

// PaginatedTriggersResponse represents a paginated list of triggers
type PaginatedTriggersResponse struct {
	Triggers []TriggerInfo `json:"triggers"`
	Limit    int           `json:"limit"`
	Offset   int           `json:"offset"`
	Count    int           `json:"count"`
	Total    int           `json:"total"`
}

// PaginatedActionsResponse represents a paginated list of actions
type PaginatedActionsResponse struct {
	Actions []ActionInfo `json:"actions"`
	Limit   int          `json:"limit"`
	Offset  int          `json:"offset"`
	Count   int          `json:"count"`
	Total   int          `json:"total"`
}

// ErrorResponse represents an error response from the API
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail represents the details of an error
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
