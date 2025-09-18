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
	params := map[string]any{
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

func TestRegisterAPIFunctions_LogMessage(t *testing.T) {
	service := NewService()
	L := lua.NewState()
	defer L.Close()

	service.RegisterAPIFunctions(L)

	// Test log_message function
	err := L.DoString(`
		local logger = require 'logger'

		logger.info('Test log message from Lua')
		logger.error('Test error message')
	`)

	assert.NoError(t, err)
}

func TestRegisterAPIFunctions_GetCurrentTime(t *testing.T) {
	service := NewService()
	L := lua.NewState()
	defer L.Close()

	service.RegisterAPIFunctions(L)

	// Test get_current_time function
	err := L.DoString(`
		local time = require 'time'

		local timestamp = time.now()
		assert(type(timestamp) == "string")
		assert(timestamp ~= "")
	`)

	assert.NoError(t, err)
}

func TestRegisterAPIFunctions_ErrorHandling(t *testing.T) {
	service := NewService()
	L := lua.NewState()
	defer L.Close()

	service.RegisterAPIFunctions(L)

	// Test error handling for empty device ID
	err := L.DoString(`
		local http = require 'http'

		local state, err = http.get("", {}, "")
		if not err then
			error("Expected error for empty HTTP URL")
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
	mapData := map[string]any{"key": "value"}
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

	service.RegisterAPIFunctions(L)

	// Test a complete Lua script that validates that all API functions are available
	script := `
	local http = require 'http'
	assert(type(http) == "table")
	assert(http.get ~= nil)
	assert(type(http.get) == "function")
	assert(http.post ~= nil)
	assert(type(http.post) == "function")
	assert(http.delete ~= nil)
	assert(type(http.delete) == "function")
	assert(http.put ~= nil)
	assert(type(http.put) == "function")
	assert(http.patch ~= nil)
	assert(type(http.patch) == "function")

	local logger = require 'logger'
	assert(type(logger) == "table")
	assert(logger.info ~= nil)
	assert(type(logger.info) == "function")
	assert(logger.debug ~= nil)
	assert(type(logger.debug) == "function")
	assert(logger.warn ~= nil)
	assert(type(logger.warn) == "function")
	assert(logger.error ~= nil)
	assert(type(logger.error) == "function")

	local time = require 'time'
	assert(type(time) == "table")
	assert(time.now ~= nil)
	assert(type(time.now) == "function")
	`

	err := L.DoString(script)
	assert.NoError(t, err)
}
