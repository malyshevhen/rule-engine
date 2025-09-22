package api

import (
	"time"

	"github.com/google/uuid"
	"github.com/malyshevhen/rule-engine/internal/action"
	"github.com/malyshevhen/rule-engine/internal/rule"
	"github.com/malyshevhen/rule-engine/internal/trigger"
)

// APIConfig holds configuration values for the API layer
type APIConfig struct {
	MaxRuleNameLength  int
	MaxLuaScriptLength int
	DefaultRulesLimit  int
	MaxRulesLimit      int
	DefaultRulesOffset int
}

// DefaultAPIConfig returns the default API configuration
func DefaultAPIConfig() *APIConfig {
	return &APIConfig{
		MaxRuleNameLength:  255,
		MaxLuaScriptLength: 10000,
		DefaultRulesLimit:  50,
		MaxRulesLimit:      1000,
		DefaultRulesOffset: 0,
	}
}

// Global API configuration instance
var apiConfig = DefaultAPIConfig()

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
	Name      string    `json:"name"`
	LuaScript string    `json:"lua_script"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateRuleRequest represents a request to create a rule
type CreateRuleRequest struct {
	Name      string `json:"name" validate:"required,rule_name_length" example:"Temperature Alert Rule"`
	LuaScript string `json:"lua_script" validate:"required,lua_script_length" example:"if event.temperature > 25 then return true end"`
	Priority  *int   `json:"priority,omitempty" example:"0"`
	Enabled   *bool  `json:"enabled,omitempty" example:"true"`
}

// UpdateRuleRequest represents a request to update a rule
type UpdateRuleRequest struct {
	Name      *string `json:"name,omitempty" validate:"omitempty,rule_name_length" example:"Updated Rule Name"`
	LuaScript *string `json:"lua_script,omitempty" validate:"omitempty,lua_script_length" example:"if event.temperature > 30 then return true end"`
	Priority  *int    `json:"priority,omitempty" example:"5"`
	Enabled   *bool   `json:"enabled,omitempty" example:"false"`
}

// CreateTriggerRequest represents a request to create a trigger
type CreateTriggerRequest struct {
	RuleID          uuid.UUID `json:"rule_id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Type            string    `json:"type" validate:"required,oneof=CONDITIONAL CRON" example:"CONDITIONAL"`
	ConditionScript string    `json:"condition_script" validate:"required,lua_script_length" example:"if event.device_id == 'sensor_1' then return true end"`
	Enabled         *bool     `json:"enabled,omitempty" example:"true"`
}

// CreateActionRequest represents a request to create an action
type CreateActionRequest struct {
	Name      string `json:"name" example:"Send Temperature Alert"`
	LuaScript string `json:"lua_script" validate:"required,lua_script_length" example:"log_message('info', 'Temperature alert triggered')"`
	Enabled   *bool  `json:"enabled,omitempty" example:"true"`
}

// EvaluateScriptRequest represents a request to evaluate a Lua script
type EvaluateScriptRequest struct {
	Script  string         `json:"script" validate:"required,lua_script_length" example:"return 2 + 2"`
	Context map[string]any `json:"context,omitempty"`
}

// EvaluateScriptResponse represents the result of script evaluation
type EvaluateScriptResponse struct {
	Success  bool   `json:"success" example:"true"`
	Result   any    `json:"result,omitempty"`
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

// JSON Patch (RFC 6902) types

// PatchOperation represents a single JSON Patch operation
type PatchOperation struct {
	Op    string `json:"op" validate:"required,oneof=add remove replace test" example:"replace"`
	Path  string `json:"path" validate:"required" example:"/name"`
	Value any    `json:"value,omitempty"`
}

// PatchRequest represents a JSON Patch request containing multiple operations
type PatchRequest []PatchOperation

// Conversion functions from domain models to DTOs

// RuleToRuleInfo converts a rule domain model to RuleInfo DTO
func RuleToRuleInfo(r *rule.Rule) *RuleInfo {
	return &RuleInfo{
		ID:        r.ID,
		Name:      r.Name,
		LuaScript: r.LuaScript,
		Priority:  r.Priority,
		Enabled:   r.Enabled,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
		Triggers:  make([]TriggerInfo, len(r.Triggers)),
		Actions:   make([]ActionInfo, len(r.Actions)),
	}
}

// RulesToRuleInfos converts a slice of rule domain models to RuleInfo DTOs
func RulesToRuleInfos(rules []*rule.Rule) []*RuleInfo {
	result := make([]*RuleInfo, len(rules))
	for i, r := range rules {
		result[i] = RuleToRuleInfo(r)
	}
	return result
}

// TriggerToTriggerInfo converts a trigger domain model to TriggerInfo DTO
func TriggerToTriggerInfo(t *trigger.Trigger) *TriggerInfo {
	return &TriggerInfo{
		ID:              t.ID,
		RuleID:          t.RuleID,
		Type:            string(t.Type),
		ConditionScript: t.ConditionScript,
		Enabled:         t.Enabled,
		CreatedAt:       t.CreatedAt,
		UpdatedAt:       t.UpdatedAt,
	}
}

// ActionToActionInfo converts an action domain model to ActionInfo DTO
func ActionToActionInfo(a *action.Action) *ActionInfo {
	return &ActionInfo{
		ID:        a.ID,
		Name:      a.Name,
		LuaScript: a.LuaScript,
		Enabled:   a.Enabled,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}
