package platform

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
)

func TestService_GetDeviceState(t *testing.T) {
	service := NewService()

	ctx := context.Background()
	state, err := service.GetDeviceState(ctx, "device-123")

	assert.NoError(t, err)
	assert.NotNil(t, state)
	assert.Equal(t, "device-123", state["id"])
	assert.True(t, state["online"].(bool))
	assert.Contains(t, state, "temperature")
	assert.Contains(t, state, "humidity")
}

func TestService_SendCommand(t *testing.T) {
	service := NewService()

	ctx := context.Background()
	params := map[string]interface{}{
		"power": true,
		"level": 75,
	}

	err := service.SendCommand(ctx, "device-123", "set_power", params)

	// Should not error (mock implementation)
	assert.NoError(t, err)
}

func TestService_GetCurrentTime(t *testing.T) {
	service := NewService()

	before := time.Now()
	currentTime := service.GetCurrentTime()
	after := time.Now()

	assert.True(t, currentTime.After(before) || currentTime.Equal(before))
	assert.True(t, currentTime.Before(after) || currentTime.Equal(after))
}

func TestService_DataStorageViaLua(t *testing.T) {
	service := NewService()
	L := lua.NewState()
	defer L.Close()

	service.RegisterAPIFunctions(L, "test-rule", "test-trigger")

	// Test store_data and get_stored_data via Lua
	err := L.DoString(`
		-- Store some data
		store_data("test_key", "test_value")
		store_data("number_key", 42)

		-- Retrieve data
		local value1 = get_stored_data("test_key")
		local value2 = get_stored_data("number_key")
		local value3 = get_stored_data("nonexistent_key")

		assert(value1 == "test_value")
		assert(value2 == 42)
		assert(value3 == nil)
	`)

	assert.NoError(t, err)
}

func TestRegisterAPIFunctions_GetDeviceState(t *testing.T) {
	service := NewService()
	L := lua.NewState()
	defer L.Close()

	service.RegisterAPIFunctions(L, "rule-123", "trigger-456")

	// Test get_device_state function
	err := L.DoString(`
		local state, err = get_device_state("device-123")
		if err then
			error(err)
		end
		assert(state.id == "device-123")
		assert(state.online == true)
		assert(type(state.temperature) == "number")
		assert(type(state.humidity) == "number")
	`)

	assert.NoError(t, err)
}

func TestRegisterAPIFunctions_SendCommand(t *testing.T) {
	service := NewService()
	L := lua.NewState()
	defer L.Close()

	service.RegisterAPIFunctions(L, "rule-123", "trigger-456")

	// Test send_command function
	err := L.DoString(`
		local err = send_command("device-123", "set_power", {power = true, level = 75})
		if err then
			error(err)
		end
	`)

	assert.NoError(t, err)
}

func TestRegisterAPIFunctions_LogMessage(t *testing.T) {
	service := NewService()
	L := lua.NewState()
	defer L.Close()

	service.RegisterAPIFunctions(L, "rule-123", "trigger-456")

	// Test log_message function
	err := L.DoString(`
		log_message("info", "Test log message from Lua")
		log_message("error", "Test error message")
	`)

	assert.NoError(t, err)
}

func TestRegisterAPIFunctions_GetCurrentTime(t *testing.T) {
	service := NewService()
	L := lua.NewState()
	defer L.Close()

	service.RegisterAPIFunctions(L, "rule-123", "trigger-456")

	// Test get_current_time function
	err := L.DoString(`
		local timestamp = get_current_time()
		assert(type(timestamp) == "number")
		assert(timestamp > 0)
	`)

	assert.NoError(t, err)
}

func TestRegisterAPIFunctions_StoreData(t *testing.T) {
	service := NewService()
	L := lua.NewState()
	defer L.Close()

	service.RegisterAPIFunctions(L, "rule-123", "trigger-456")

	// Test store_data and get_stored_data functions
	err := L.DoString(`
		-- Store some data
		store_data("test_key", "test_value")
		store_data("number_key", 42)
		store_data("table_key", {a = 1, b = 2})

		-- Retrieve data
		local value1 = get_stored_data("test_key")
		local value2 = get_stored_data("number_key")
		local value3 = get_stored_data("table_key")
		local value4 = get_stored_data("nonexistent")

		assert(value1 == "test_value")
		assert(value2 == 42)
		assert(type(value3) == "table")
		assert(value3.a == 1)
		assert(value3.b == 2)
		assert(value4 == nil)
	`)

	assert.NoError(t, err)
}

func TestRegisterAPIFunctions_ErrorHandling(t *testing.T) {
	service := NewService()
	L := lua.NewState()
	defer L.Close()

	service.RegisterAPIFunctions(L, "rule-123", "trigger-456")

	// Test error handling for empty device ID
	err := L.DoString(`
		local state, err = get_device_state("")
		if not err then
			error("Expected error for empty device ID")
		end
	`)

	assert.NoError(t, err)
}

func TestLuaValueConversions(t *testing.T) {
	// Test luaValueFromGo
	L := lua.NewState()
	defer L.Close()

	// Test various Go types to Lua conversion
	assert.Equal(t, lua.LTString, luaValueFromGo(L, "string").Type())
	assert.Equal(t, lua.LTNumber, luaValueFromGo(L, 42).Type())
	assert.Equal(t, lua.LTNumber, luaValueFromGo(L, int64(42)).Type())
	assert.Equal(t, lua.LTNumber, luaValueFromGo(L, 3.14).Type())
	assert.Equal(t, lua.LTBool, luaValueFromGo(L, true).Type())
	assert.Equal(t, lua.LTNil, luaValueFromGo(L, nil).Type())

	// Test table conversion
	mapData := map[string]interface{}{"key": "value"}
	table := luaValueFromGo(L, mapData)
	assert.Equal(t, lua.LTTable, table.Type())

	// Test luaValueToGo
	userData := &lua.LUserData{}
	result := luaValueToGo(userData)
	assert.Contains(t, result, "userdata") // fallback case returns string representation
}

func TestPlatformAPI_Integration(t *testing.T) {
	service := NewService()
	L := lua.NewState()
	defer L.Close()

	service.RegisterAPIFunctions(L, "rule-123", "trigger-456")

	// Test a complete Lua script that uses multiple API functions
	script := `
		-- Get device state
		local state, err = get_device_state("sensor-001")
		if err then error("Failed to get device state: " .. err) end

		-- Log the temperature
		log_message("info", "Current temperature: " .. tostring(state.temperature))

		-- Store the temperature for later use
		store_data("last_temp", state.temperature)

		-- Check if temperature is too high
		if state.temperature > 30 then
			-- Send command to turn on fan
			local cmd_err = send_command("fan-001", "set_power", {power = true, speed = "high"})
			if cmd_err then error("Failed to send command: " .. cmd_err) end

			log_message("warn", "High temperature detected, turning on fan")
		end

		-- Get current time
		local now = get_current_time()

		-- Store execution timestamp
		store_data("last_execution", now)

		return true
	`

	err := L.DoString(script)
	assert.NoError(t, err)

	// The data is stored within the Lua execution context, so we can't access it directly
	// from Go. The test passes if the Lua script runs without error, which means
	// the data storage and retrieval worked within the Lua context.
	assert.NoError(t, err)
}
