// Package obsx holds the shared logging setup.
package obsx

import (
	"log/slog"
	"os"
)

// NewLogger returns the standard JSON logger tagged with the service name.
func NewLogger(service string) *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: levelFromEnv()})

	return slog.New(handler).With("service", service)
}

func levelFromEnv() slog.Level {
	var level slog.Level
	if err := level.UnmarshalText([]byte(os.Getenv("LOG_LEVEL"))); err != nil {
		return slog.LevelInfo
	}

	return level
}
