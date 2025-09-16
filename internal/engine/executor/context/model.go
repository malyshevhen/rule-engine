// Package context provides execution context for Lua scripts
package context

// ExecutionContext holds data available to Lua scripts during execution
type ExecutionContext struct {
	// TODO: Add fields like device state, user data, etc.
	RuleID    string                 `json:"rule_id"`
	TriggerID string                 `json:"trigger_id"`
	Data      map[string]interface{} `json:"data"`
}
