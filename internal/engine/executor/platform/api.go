package platform

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	lua "github.com/yuin/gopher-lua"
)

// PlatformAPI provides functions that Lua scripts can call to interact with the IoT platform
type PlatformAPI interface {
	// GetDeviceState retrieves the current state of a device
	GetDeviceState(ctx context.Context, deviceID string) (map[string]interface{}, error)

	// SendCommand sends a command to a device
	SendCommand(ctx context.Context, deviceID string, command string, params map[string]interface{}) error

	// LogMessage logs a message from a Lua script
	LogMessage(ctx context.Context, level, message string)

	// GetCurrentTime returns the current timestamp
	GetCurrentTime() time.Time

	// StoreData stores temporary data during rule execution
	StoreData(ctx context.Context, key string, value interface{})

	// GetStoredData retrieves stored data
	GetStoredData(ctx context.Context, key string) interface{}
}

// Service implements the PlatformAPI interface
type Service struct {
	// TODO: Add dependencies like device service, message bus, cache, etc.
	dataStores map[string]map[string]interface{} // Per-rule execution storage
}

// NewService creates a new platform API service
func NewService() *Service {
	return &Service{
		dataStores: make(map[string]map[string]interface{}),
	}
}

// GetDeviceState retrieves the current state of a device
func (s *Service) GetDeviceState(ctx context.Context, deviceID string) (map[string]interface{}, error) {
	// TODO: Implement actual device state retrieval from device service
	// For now, return mock data
	return map[string]interface{}{
		"id":          deviceID,
		"online":      true,
		"last_seen":   time.Now().Unix(),
		"temperature": 25.5,
		"humidity":    60.0,
	}, nil
}

// SendCommand sends a command to a device
func (s *Service) SendCommand(ctx context.Context, deviceID string, command string, params map[string]interface{}) error {
	// TODO: Implement actual command sending via message bus
	slog.Info("Sending command to device",
		"device_id", deviceID,
		"command", command,
		"params", params)

	// Simulate command sending
	return nil
}

// LogMessage logs a message from a Lua script
func (s *Service) LogMessage(ctx context.Context, level, message string) {
	switch level {
	case "debug":
		slog.Debug("Lua script message", "message", message)
	case "info":
		slog.Info("Lua script message", "message", message)
	case "warn", "warning":
		slog.Warn("Lua script message", "message", message)
	case "error":
		slog.Error("Lua script message", "message", message)
	default:
		slog.Info("Lua script message", "level", level, "message", message)
	}
}

// GetCurrentTime returns the current timestamp
func (s *Service) GetCurrentTime() time.Time {
	return time.Now()
}

// RegisterAPIFunctions registers platform API functions in the Lua state
func (s *Service) RegisterAPIFunctions(L *lua.LState, ruleID, triggerID string) {
	// get_device_state(device_id)
	L.SetGlobal("get_device_state", L.NewFunction(func(L *lua.LState) int {
		deviceID := L.ToString(1)
		if deviceID == "" {
			L.Push(lua.LNil)
			L.Push(lua.LString("device_id cannot be empty"))
			return 2
		}

		ctx := context.Background()
		state, err := s.GetDeviceState(ctx, deviceID)
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}

		// Convert Go map to Lua table
		table := L.NewTable()
		for k, v := range state {
			L.SetField(table, k, luaValueFromGo(L, v))
		}

		L.Push(table)
		L.Push(lua.LNil)
		return 2
	}))

	// send_command(device_id, command, params_table)
	L.SetGlobal("send_command", L.NewFunction(func(L *lua.LState) int {
		deviceID := L.ToString(1)
		command := L.ToString(2)
		paramsTable := L.ToTable(3)

		if deviceID == "" || command == "" {
			L.Push(lua.LString("device_id and command cannot be empty"))
			return 1
		}

		// Convert Lua table to Go map
		params := make(map[string]interface{})
		if paramsTable != nil {
			paramsTable.ForEach(func(k, v lua.LValue) {
				params[k.String()] = luaValueToGo(v)
			})
		}

		ctx := context.Background()
		err := s.SendCommand(ctx, deviceID, command, params)
		if err != nil {
			L.Push(lua.LString(err.Error()))
			return 1
		}

		L.Push(lua.LNil)
		return 1
	}))

	// log_message(level, message)
	L.SetGlobal("log_message", L.NewFunction(func(L *lua.LState) int {
		level := L.ToString(1)
		message := L.ToString(2)

		ctx := context.Background()
		s.LogMessage(ctx, level, message)

		return 0
	}))

	// get_current_time()
	L.SetGlobal("get_current_time", L.NewFunction(func(L *lua.LState) int {
		currentTime := s.GetCurrentTime()
		L.Push(lua.LNumber(currentTime.Unix()))
		return 1
	}))

	// store_data(key, value)
	L.SetGlobal("store_data", L.NewFunction(func(L *lua.LState) int {
		key := L.ToString(1)
		value := L.Get(2)

		if key == "" {
			L.Push(lua.LString("key cannot be empty"))
			return 1
		}

		// Use execution-specific key
		execKey := ruleID + ":" + triggerID
		if s.dataStores[execKey] == nil {
			s.dataStores[execKey] = make(map[string]interface{})
		}
		s.dataStores[execKey][key] = luaValueToGo(value)

		return 0
	}))

	// get_stored_data(key)
	L.SetGlobal("get_stored_data", L.NewFunction(func(L *lua.LState) int {
		key := L.ToString(1)

		// Use execution-specific key
		execKey := ruleID + ":" + triggerID
		if store, exists := s.dataStores[execKey]; exists {
			if value, exists := store[key]; exists {
				L.Push(luaValueFromGo(L, value))
				return 1
			}
		}

		L.Push(lua.LNil)
		return 1
	}))
}

// luaValueFromGo converts a Go value to a Lua value
func luaValueFromGo(L *lua.LState, v interface{}) lua.LValue {
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
	case map[string]interface{}:
		table := L.NewTable()
		for k, v := range val {
			L.SetField(table, k, luaValueFromGo(L, v))
		}
		return table
	case []interface{}:
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
		table.ForEach(func(k, v lua.LValue) {
			result[k.String()] = luaValueToGo(v)
		})
		return result
	default:
		return v.String()
	}
}
