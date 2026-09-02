// Package config owns environment configuration.
package config

import "time"

type Config struct {
	TargetURL string
	// Duration bounds the run, zero means until SIGTERM.
	Duration time.Duration
	// AgentKeySeed derives the agent signing key, empty means unsigned.
	AgentKeySeed string
}
