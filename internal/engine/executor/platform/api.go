package platform

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/malyshevhen/rule-engine/internal/engine/executor/platform/modules"
	lua "github.com/yuin/gopher-lua"
)

type Module interface {
	Name() string
	Loader(L *lua.LState) int
}

// Service implements the PlatformAPI interface
type Service struct {
	ms []Module
}

// NewService creates a new platform API service
func NewService() *Service {
	ms := []Module{
		modules.NewLoggerModule(),
		modules.NewHTTPModule(),
		modules.NewTimeModule(),
	}
	return &Service{ms: ms}
}

// GetDeviceState retrieves the current state of a device
func (s *Service) GetDeviceState(ctx context.Context, deviceID string) (map[string]any, error) {
	// TODO: Implement actual device state retrieval from device service
	// For now, return mock data
	return map[string]any{
		"id":          deviceID,
		"online":      true,
		"last_seen":   time.Now().Unix(),
		"temperature": 25.5,
		"humidity":    60.0,
	}, nil
}

// SendCommand sends a command to a device
func (s *Service) SendCommand(ctx context.Context, deviceID string, command string, params map[string]any) error {
	// TODO: Implement actual command sending via message bus
	slog.Info("Sending command to device",
		"device_id", deviceID,
		"command", command,
		"params", params)

	// Simulate command sending
	return nil
}

// GetCurrentTime returns the current timestamp
func (s *Service) GetCurrentTime() time.Time {
	return time.Now()
}

// RegisterAPIFunctions registers platform API functions in the Lua state
func (s *Service) RegisterAPIFunctions(L *lua.LState) {
	// Register modules
	for _, module := range s.ms {
		L.PreloadModule(module.Name(), module.Loader)
	}
}

// luaValueFromGo converts a Go value to a Lua value
func luaValueFromGo(L *lua.LState, v any) lua.LValue {
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
	case map[string]any:
		table := L.NewTable()
		for k, v := range val {
			L.SetField(table, k, luaValueFromGo(L, v))
		}
		return table
	case []any:
		table := L.NewTable()
		for i, v := range val {
			table.Insert(i+1, luaValueFromGo(L, v))
		}
		return table
	default:
		// For complex types, convert to JSON string
		if jsonBytes, err := json.Marshal(v); err == nil {
			return lua.LString(jsonBytes)
		}
		return lua.LString(fmt.Sprintf("%v", v))
	}
}

// luaValueToGo converts a Lua value to a Go interface{}
func luaValueToGo(v lua.LValue) any {
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
		result := make(map[string]any)
		table.ForEach(func(k, v lua.LValue) {
			result[k.String()] = luaValueToGo(v)
		})
		return result
	default:
		return v.String()
	}
}
