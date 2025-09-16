package executor

import (
	"context"
	"time"

	execCtx "github.com/malyshevhen/rule-engine/internal/engine/executor/context"
	lua "github.com/yuin/gopher-lua"
)

// Service handles Lua script execution
type Service struct {
	contextService *execCtx.Service
}

// NewService creates a new executor service
func NewService(contextService *execCtx.Service) *Service {
	return &Service{
		contextService: contextService,
	}
}

// ExecuteResult represents the result of script execution
type ExecuteResult struct {
	Success  bool          `json:"success"`
	Output   interface{}   `json:"output"`
	Error    string        `json:"error,omitempty"`
	Duration time.Duration `json:"duration"`
}

// ExecuteScript executes a Lua script with sandboxing
func (s *Service) ExecuteScript(ctx context.Context, script string, execCtx *execCtx.ExecutionContext) *ExecuteResult {
	start := time.Now()

	// Create a new Lua state with sandboxed options
	L := lua.NewState(lua.Options{
		SkipOpenLibs: true, // Don't open default libraries
	})

	defer L.Close()

	// Open only safe libraries
	L.OpenLibs() // This opens all, but we can selectively open

	// Remove unsafe libraries
	L.SetGlobal("io", lua.LNil)
	L.SetGlobal("os", lua.LNil)
	// TODO: Remove networking libraries if any

	// Set execution context in Lua
	L.SetGlobal("rule_id", lua.LString(execCtx.RuleID))
	L.SetGlobal("trigger_id", lua.LString(execCtx.TriggerID))
	// TODO: Set other context data

	// Execute the script
	err := L.DoString(script)
	duration := time.Since(start)

	if err != nil {
		return &ExecuteResult{
			Success:  false,
			Error:    err.Error(),
			Duration: duration,
		}
	}

	// Get the result (assume script returns a value)
	result := L.Get(-1)
	output := result.String()

	return &ExecuteResult{
		Success:  true,
		Output:   output,
		Duration: duration,
	}
}
