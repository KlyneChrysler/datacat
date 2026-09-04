package config

import (
	"os"
	"time"

	"github.com/KlyneChrysler/datacat/pkg/envx"
)

// Load reads and validates all configuration.
func Load() (Config, error) {
	cfg := Config{
		TargetURL:    os.Getenv("TARGET_URL"),
		Duration:     time.Duration(envx.Int("DURATION_SECONDS", 0)) * time.Second,
		AgentKeySeed: os.Getenv("AGENT_KEY_SEED"),
	}

	return cfg, validate(cfg)
}

func validate(c Config) error {
	return envx.Require(map[string]string{"TARGET_URL": c.TargetURL})
}
