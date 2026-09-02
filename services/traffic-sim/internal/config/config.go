// Package config owns environment configuration: the shape here, loading
// and validation in load.go (twelve-factor III).
package config

import "time"

type Config struct {
	TargetURL string
	// Duration bounds the simulation; zero runs until SIGTERM. Optional env
	// DURATION_SECONDS.
	Duration time.Duration
}
