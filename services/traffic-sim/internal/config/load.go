package config

import (
	"fmt"
	"os"
	"time"

	"github.com/KlyneChrysler/datacat/pkg/envx"
)

// Load reads and validates all configuration at startup.
func Load() (Config, error) {
	cfg := Config{
		TargetURL:    os.Getenv("TARGET_URL"),
		Duration:     time.Duration(envx.Int("DURATION_SECONDS", 0)) * time.Second,
		AgentKeySeed: os.Getenv("AGENT_KEY_SEED"),
	}
	return cfg, validate(cfg)
}

func validate(c Config) error {
	if c.TargetURL == "" {
		return fmt.Errorf("config: TARGET_URL is required")
	}
	return nil
}
