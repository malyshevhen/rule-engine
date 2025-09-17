package tools

import "log/slog"

type LoggerAPI interface {
	LogMessage(level, message string)
}

// LoggerAPIService implements the LoggerAPI interface
type LoggerAPIService struct {
}

// NewLoggerAPIService creates a new LoggerAPIService
func NewLoggerAPIService() *LoggerAPIService {
	return &LoggerAPIService{}
}

// LogMessage logs a message
func (s *LoggerAPIService) LogMessage(level, message string) {
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
