package modules

import (
	"context"
	"log/slog"

	lua "github.com/yuin/gopher-lua"
)

// LogLevel represents a log level
type LogLevel string

// Log levels enum
const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

// Logger interface for logging
type Logger interface {
	Info(message string, args ...any)
	Debug(message string, args ...any)
	Warn(message string, args ...any)
	Error(message string, args ...any)
}

// LoggerModuleOption allows to configure the logger module
type LoggerModuleOption func(lm *LoggerModule) *LoggerModule

// WithLogger sets a custom logger
func WithLogger(logger Logger) LoggerModuleOption {
	return func(lm *LoggerModule) *LoggerModule {
		lm.logger = logger
		return lm
	}
}

// LoggerModule provides functions that Lua scripts can call
type LoggerModule struct {
	logger Logger
}

// NewLoggerModule creates a new LoggerModule
func NewLoggerModule(opts ...LoggerModuleOption) *LoggerModule {
	lm := &LoggerModule{logger: slog.Default()}

	for _, opt := range opts {
		opt(lm)
	}
	return lm
}

// LogMessage logs a message
func (s *LoggerModule) Info(L *lua.LState) int {
	message := L.ToString(1)
	L.Pop(1)

	ctx := context.Background()
	s.LogMessage(ctx, "info", message)

	return 0
}

// LogMessage logs a message
func (s *LoggerModule) Debug(L *lua.LState) int {
	message := L.ToString(1)
	L.Pop(1)

	ctx := context.Background()
	s.LogMessage(ctx, "debug", message)

	return 0
}

// LogMessage logs a message
func (s *LoggerModule) Warn(L *lua.LState) int {
	message := L.ToString(1)
	L.Pop(1)

	ctx := context.Background()
	s.LogMessage(ctx, "warn", message)

	return 0
}

// LogMessage logs a message
func (s *LoggerModule) Error(L *lua.LState) int {
	message := L.ToString(1)
	L.Pop(1)

	ctx := context.Background()
	s.LogMessage(ctx, "error", message)

	return 0
}

// LogMessage logs a message
func (s *LoggerModule) LogMessage(ctx context.Context, level LogLevel, message string) {
	switch level {
	case LogLevelDebug:
		s.logger.Debug("Lua script message", "message", message)
	case LogLevelInfo:
		s.logger.Info("Lua script message", "message", message)
	case LogLevelWarn:
		s.logger.Warn("Lua script message", "message", message)
	case LogLevelError:
		s.logger.Error("Lua script message", "message", message)
	default:
		s.logger.Info("Lua script message", "level", level, "message", message)
	}
}

// Loader loads the Logger module into the Lua state
func (s *LoggerModule) Loader(L *lua.LState) int {
	exports := map[string]lua.LGFunction{
		"info":  s.Info,
		"debug": s.Debug,
		"warn":  s.Warn,
		"error": s.Error,
	}

	mod := L.SetFuncs(L.NewTable(), exports)
	L.Push(mod)
	return 1
}
