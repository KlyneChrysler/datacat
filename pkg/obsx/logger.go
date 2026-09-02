// Package obsx provides observability primitives shared by all datacat
// services: structured JSON logging to stdout (twelve-factor XI).
package obsx

import (
	"log/slog"
	"os"
)

// NewLogger returns the canonical datacat logger: JSON lines on stdout,
// tagged with the service name. The platform owns routing and retention.
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
