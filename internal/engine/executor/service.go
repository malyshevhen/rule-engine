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

// GetContextService returns the context service
func (s *Service) GetContextService() *execCtx.Service {
	return s.contextService
}

// ExecuteResult represents the result of script execution
type ExecuteResult struct {
	Success  bool          `json:"success"`
	Output   []interface{} `json:"output"`
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

	// Open only essential safe libraries
	lua.OpenBase(L)   // _G, basic functions
	lua.OpenTable(L)  // table library
	lua.OpenString(L) // string library
	lua.OpenMath(L)   // math library
	// Explicitly do NOT open: io, os, debug, package, coroutine (if not needed)

	// Remove any potentially unsafe globals that might be set
	L.SetGlobal("io", lua.LNil)
	L.SetGlobal("os", lua.LNil)
	L.SetGlobal("debug", lua.LNil)
	L.SetGlobal("package", lua.LNil)
	L.SetGlobal("coroutine", lua.LNil)
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

	// Get all return values
	top := L.GetTop()
	results := make([]interface{}, 0, top)
	for i := 1; i <= top; i++ {
		result := L.Get(i)
		results = append(results, luaValueToGo(result))
	}

	return &ExecuteResult{
		Success:  true,
		Output:   results,
		Duration: duration,
	}
}

// luaValueToGo converts a Lua value to a Go interface{}
func luaValueToGo(v lua.LValue) interface{} {
	switch v.Type() {
	case lua.LTNil:
		return nil
	case lua.LTBool:
		return lua.LVAsBool(v)
	case lua.LTNumber:
		return float64(v.(lua.LNumber))
	case lua.LTString:
		return string(v.(lua.LString))
	case lua.LTTable:
		// Convert table to map
		table := v.(*lua.LTable)
		result := make(map[string]interface{})
		table.ForEach(func(key, value lua.LValue) {
			if key.Type() == lua.LTString {
				result[string(key.(lua.LString))] = luaValueToGo(value)
			}
		})
		return result
	default:
		return v.String()
	}
}
