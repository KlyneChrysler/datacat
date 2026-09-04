package config

import (
	"os"
	"time"

	"github.com/KlyneChrysler/datacat/pkg/envx"
	"github.com/KlyneChrysler/datacat/pkg/guard"
)

// Load reads and validates all configuration, crashing on missing values.
func Load() (Config, error) {
	upstream, err := envx.Upstream()
	if err != nil {
		return Config{}, err
	}
	agentKeys, err := guard.ParseAgentKeys(os.Getenv("AGENT_KEYS"))
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		Port:                os.Getenv("PORT"),
		UpstreamURL:         upstream,
		ChallengeSecret:     os.Getenv("CHALLENGE_SECRET"),
		ChallengeDifficulty: envx.Int("CHALLENGE_DIFFICULTY_BITS", 16),
		RateLimitPerMinute:  envx.Int("RATE_LIMIT_PER_MINUTE", 60),
		RateLimitBurst:      envx.Int("RATE_LIMIT_BURST", 10),
		AgentKeys:           agentKeys,
		GateTTL:             time.Hour,
		ShutdownTimeout:     10 * time.Second,
	}

	return cfg, validate(cfg)
}

func validate(c Config) error {
	return envx.Require(map[string]string{"PORT": c.Port, "CHALLENGE_SECRET": c.ChallengeSecret})
}
