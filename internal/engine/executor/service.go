package executor

import (
	"context"
	"fmt"
	"time"

	execCtx "github.com/malyshevhen/rule-engine/internal/engine/executor/context"
	"github.com/malyshevhen/rule-engine/internal/engine/executor/platform"
	"github.com/malyshevhen/rule-engine/internal/metrics"
	lua "github.com/yuin/gopher-lua"
)

// Service handles Lua script execution
type Service struct {
	contextService *execCtx.Service
	platformAPI    *platform.Service
}

// NewService creates a new executor service
func NewService(contextService *execCtx.Service, platformAPI *platform.Service) *Service {
	return &Service{
		contextService: contextService,
		platformAPI:    platformAPI,
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

	// Inject event data into Lua globals
	for key, value := range execCtx.Data {
		L.SetGlobal(key, luaValueToLValue(value))
	}

	// Register platform API functions
	s.platformAPI.RegisterAPIFunctions(L, execCtx.RuleID, execCtx.TriggerID)

	// Execute the script
	err := L.DoString(script)
	duration := time.Since(start)

	status := "success"
	if err != nil {
		status = "failure"
		metrics.LuaExecutionErrorsTotal.WithLabelValues(execCtx.RuleID, "execution_error").Inc()
	}

	// Record metrics
	metrics.RuleExecutionsTotal.WithLabelValues(execCtx.RuleID, status).Inc()
	metrics.RuleExecutionDuration.WithLabelValues(execCtx.RuleID).Observe(duration.Seconds())

	if err != nil {
		return &ExecuteResult{
			Success:  false,
			Error:    err.Error(),
			Duration: duration,
		}
	}

	// Get the last return value (Lua scripts typically return one value)
	top := L.GetTop()
	results := make([]interface{}, 0, 1)
	if top > 0 {
		result := L.Get(top)
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
		// Avoid infinite recursion by not converting tables
		return v.String()
	default:
		return v.String()
	}
}

// luaValueToLValue converts a Go value to a Lua LValue
func luaValueToLValue(v interface{}) lua.LValue {
	switch val := v.(type) {
	case nil:
		return lua.LNil
	case bool:
		return lua.LBool(val)
	case int:
		return lua.LNumber(val)
	case int32:
		return lua.LNumber(val)
	case int64:
		return lua.LNumber(val)
	case float32:
		return lua.LNumber(val)
	case float64:
		return lua.LNumber(val)
	case string:
		return lua.LString(val)
	case []interface{}:
		// Convert slice to Lua table
		table := &lua.LTable{}
		for i, item := range val {
			table.RawSetInt(i+1, luaValueToLValue(item)) // Lua is 1-indexed
		}
		return table
	case map[string]interface{}:
		// Convert map to Lua table
		table := &lua.LTable{}
		for k, v := range val {
			table.RawSetString(k, luaValueToLValue(v))
		}
		return table
	default:
		// For unsupported types, convert to string
		return lua.LString(fmt.Sprintf("%v", v))
	}
}
