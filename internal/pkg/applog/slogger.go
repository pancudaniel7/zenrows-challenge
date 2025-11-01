package applog

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// DefaultLogger wraps slog.Logger and implements AppLogger.
type DefaultLogger struct {
	logger *slog.Logger
}

// NewAppDefaultLogger creates a new DefaultLogger configured from application settings.
func NewAppDefaultLogger() *DefaultLogger {
	levelStr := viper.GetString("log.level")
	level := parseLogLevel(levelStr)
	return &DefaultLogger{
		logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})),
	}
}

// Info proxies structured info-level logs to slog.
func (l *DefaultLogger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

// Warn emits warning-level logs.
func (l *DefaultLogger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

// Error reports failures with error severity.
func (l *DefaultLogger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

// Debug records verbose debugging information.
func (l *DefaultLogger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

// Trace logs extremely low-level traces using a custom slog level.
func (l *DefaultLogger) Trace(msg string, args ...any) {
	l.logger.Log(context.Background(), slog.Level(-8), msg, args...)
}

// Fatal logs an error and terminates the process with exit code 1.
func (l *DefaultLogger) Fatal(msg string, args ...any) {
	l.logger.Error(msg, args...)
	os.Exit(1)
}

func parseLogLevel(s string) slog.Level {
	s = strings.TrimSpace(strings.ToLower(s))
	switch s {
	case "trace":
		return slog.Level(-8)
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	case "info", "":
		return slog.LevelInfo
	default:
		return slog.LevelInfo
	}
}
