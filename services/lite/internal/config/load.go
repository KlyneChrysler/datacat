package config

import (
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/KlyneChrysler/datacat/pkg/envx"
	"github.com/KlyneChrysler/datacat/pkg/guard"
)

// Load reads and validates all configuration, crashing on missing values.
func Load() (Config, error) {
	upstream, err := parseUpstream(os.Getenv("UPSTREAM_URL"))
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

func parseUpstream(raw string) (*url.URL, error) {
	if raw == "" {
		return nil, fmt.Errorf("config: UPSTREAM_URL is required")
	}

	upstream, err := url.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("config: UPSTREAM_URL invalid: %w", err)
	}

	return upstream, nil
}

func validate(c Config) error {
	for name, v := range map[string]string{"PORT": c.Port, "CHALLENGE_SECRET": c.ChallengeSecret} {
		if v == "" {
			return fmt.Errorf("config: %s is required", name)
		}
	}

	return nil
}
