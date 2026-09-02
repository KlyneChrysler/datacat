// Package config loads and validates all environment configuration at
// startup (twelve-factor III). Missing config crashes the process at boot.
package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	Port            string
	ShutdownTimeout time.Duration
}

func Load() (Config, error) {
	cfg := Config{
		Port:            os.Getenv("PORT"),
		ShutdownTimeout: 10 * time.Second,
	}
	return cfg, cfg.validate()
}

func (c Config) validate() error {
	if c.Port == "" {
		return fmt.Errorf("config: PORT is required")
	}
	return nil
}
