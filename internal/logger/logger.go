package logger

import (
	"log/slog"
	"os"
	"strings"
)

// InitLogger initializes and registers a global structured slog logger.
// logLevel accepts "DEBUG", "INFO", "WARN", or "ERROR" (defaults to "INFO").
// jsonFormat determines whether to output logs in JSON or human-readable Text format.
func InitLogger(logLevelStr string, jsonFormat bool) *slog.Logger {
	var level slog.Level
	switch strings.ToUpper(logLevelStr) {
	case "DEBUG":
		level = slog.LevelDebug
	case "WARN", "WARNING":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	var handler slog.Handler
	if jsonFormat {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	l := slog.New(handler)
	slog.SetDefault(l)
	return l
}

// Fatal logs an error level message with attributes and terminates the application with exit code 1.
func Fatal(msg string, args ...any) {
	slog.Error(msg, args...)
	os.Exit(1)
}
