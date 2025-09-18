package modules

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
)

// mockLogger implements the Logger interface for testing
type mockLogger struct {
	calls []logCall
}

type logCall struct {
	level   string
	message string
	args    []any
}

func (m *mockLogger) Info(message string, args ...any) {
	m.calls = append(m.calls, logCall{level: "info", message: message, args: args})
}

func (m *mockLogger) Debug(message string, args ...any) {
	m.calls = append(m.calls, logCall{level: "debug", message: message, args: args})
}

func (m *mockLogger) Warn(message string, args ...any) {
	m.calls = append(m.calls, logCall{level: "warn", message: message, args: args})
}

func (m *mockLogger) Error(message string, args ...any) {
	m.calls = append(m.calls, logCall{level: "error", message: message, args: args})
}

func TestNewLoggerModule(t *testing.T) {
	mod := NewLoggerModule()
	if mod == nil {
		t.Fatal("NewLoggerModule returned nil")
	}
	if mod.logger == nil {
		t.Fatal("Logger not set")
	}
}

func TestNewLoggerModuleWithCustomLogger(t *testing.T) {
	mock := &mockLogger{}
	mod := NewLoggerModule(WithLogger(mock))
	if mod.logger != mock {
		t.Fatal("Custom logger not set")
	}
}

func TestLogMessage(t *testing.T) {
	mock := &mockLogger{}
	mod := NewLoggerModule(WithLogger(mock))
	ctx := context.Background()

	testCases := []struct {
		level   LogLevel
		message string
	}{
		{LogLevelDebug, "debug message"},
		{LogLevelInfo, "info message"},
		{LogLevelWarn, "warn message"},
		{LogLevelError, "error message"},
		{"unknown", "unknown level message"},
	}

	for _, tc := range testCases {
		t.Run(string(tc.level), func(t *testing.T) {
			mock.calls = nil // reset
			mod.LogMessage(ctx, tc.level, tc.message)

			if len(mock.calls) != 1 {
				t.Fatalf("Expected 1 call, got %d", len(mock.calls))
			}

			call := mock.calls[0]
			if call.message != "Lua script message" {
				t.Errorf("Expected message 'Lua script message', got '%s'", call.message)
			}

			if tc.level == "unknown" {
				// For unknown level, should be info with level and message
				if call.level != "info" {
					t.Errorf("Expected level 'info', got '%s'", call.level)
				}
				expectedArgs := []any{"level", tc.level, "message", tc.message}
				if len(call.args) != len(expectedArgs) {
					t.Errorf("Expected %d args, got %d", len(expectedArgs), len(call.args))
				}
				for i, arg := range expectedArgs {
					if i >= len(call.args) || call.args[i] != arg {
						t.Errorf("Arg %d: expected %v, got %v", i, arg, call.args[i])
					}
				}
			} else {
				if call.level != string(tc.level) {
					t.Errorf("Expected level '%s', got '%s'", tc.level, call.level)
				}
				expectedArgs := []any{"message", tc.message}
				if len(call.args) != len(expectedArgs) {
					t.Errorf("Expected %d args, got %d", len(expectedArgs), len(call.args))
				}
				for i, arg := range expectedArgs {
					if i >= len(call.args) || call.args[i] != arg {
						t.Errorf("Arg %d: expected %v, got %v", i, arg, call.args[i])
					}
				}
			}
		})
	}
}

func TestInfo(t *testing.T) {
	mock := &mockLogger{}
	mod := NewLoggerModule(WithLogger(mock))
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("logger", mod.Loader)

	message := "test info message"
	script := `
		local logger = require 'logger'
		logger.info('` + message + `')
	`

	err := L.DoString(script)
	require.NoError(t, err)

	if L.GetTop() != 0 {
		t.Errorf("Stack not empty after Info call, top=%d", L.GetTop())
	}

	if len(mock.calls) != 1 {
		t.Fatalf("Expected 1 log call, got %d", len(mock.calls))
	}

	call := mock.calls[0]
	if call.level != "info" || call.args[1] != message {
		t.Errorf("Expected info level with message '%s', got level='%s' message='%s'", message, call.level, call.args[1])
	}
}

func TestDebug(t *testing.T) {
	mock := &mockLogger{}
	mod := NewLoggerModule(WithLogger(mock))
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("logger", mod.Loader)

	message := "test debug message"
	script := `
		local logger = require 'logger'
		logger.debug('` + message + `')
	`

	err := L.DoString(script)
	require.NoError(t, err)

	if L.GetTop() != 0 {
		t.Errorf("Stack not empty after Debug call, top=%d", L.GetTop())
	}

	if len(mock.calls) != 1 {
		t.Fatalf("Expected 1 log call, got %d", len(mock.calls))
	}

	call := mock.calls[0]
	if call.level != "debug" || call.args[1] != message {
		t.Errorf("Expected debug level with message '%s', got level='%s' message='%s'", message, call.level, call.args[1])
	}
}

func TestWarn(t *testing.T) {
	mock := &mockLogger{}
	mod := NewLoggerModule(WithLogger(mock))
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("logger", mod.Loader)

	message := "test warn message"
	script := `
		local logger = require 'logger'
		logger.warn('` + message + `')
	`

	err := L.DoString(script)
	require.NoError(t, err)

	if L.GetTop() != 0 {
		t.Errorf("Stack not empty after Warn call, top=%d", L.GetTop())
	}

	if len(mock.calls) != 1 {
		t.Fatalf("Expected 1 log call, got %d", len(mock.calls))
	}

	call := mock.calls[0]
	if call.level != "warn" || call.args[1] != message {
		t.Errorf("Expected warn level with message '%s', got level='%s' message='%s'", message, call.level, call.args[1])
	}
}

func TestError(t *testing.T) {
	mock := &mockLogger{}
	mod := NewLoggerModule(WithLogger(mock))
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("logger", mod.Loader)

	message := "test error message"
	script := `
		local logger = require 'logger'
		logger.error('` + message + `')
	`

	err := L.DoString(script)
	require.NoError(t, err)

	if L.GetTop() != 0 {
		t.Errorf("Stack not empty after Error call, top=%d", L.GetTop())
	}

	if len(mock.calls) != 1 {
		t.Fatalf("Expected 1 log call, got %d", len(mock.calls))
	}

	call := mock.calls[0]
	if call.level != "error" || call.args[1] != message {
		t.Errorf("Expected error level with message '%s', got level='%s' message='%s'", message, call.level, call.args[1])
	}
}
