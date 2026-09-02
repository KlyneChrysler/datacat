// Package config loads and validates all environment configuration at
// startup (twelve-factor III).
package config

import (
	"fmt"
	"os"
)

type Config struct {
	TargetURL string
}

func Load() (Config, error) {
	cfg := Config{TargetURL: os.Getenv("TARGET_URL")}
	return cfg, cfg.validate()
}

func (c Config) validate() error {
	if c.TargetURL == "" {
		return fmt.Errorf("config: TARGET_URL is required")
	}
	return nil
}
