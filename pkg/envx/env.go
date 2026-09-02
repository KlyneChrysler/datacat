// Package envx holds the shared optional env readers.
package envx

import (
	"os"
	"strconv"
)

// Int reads a positive int env var, falling back when unset or invalid.
func Int(name string, fallback int) int {
	raw := os.Getenv(name)
	if raw == "" {
		return fallback
	}

	value, err := strconv.Atoi(raw)
	if err != nil || value < 1 {
		return fallback
	}

	return value
}
