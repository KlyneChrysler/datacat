// Package envx provides the shared optional-env readers (Go counterpart of
// the classifier's Env.java). Required variables stay explicit in each
// service's validate — only defaulted reads live here.
package envx

import (
	"os"
	"strconv"
)

// Int reads an env var as a positive int, returning fallback when unset,
// unparsable, or < 1.
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
