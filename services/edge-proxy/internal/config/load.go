package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"
)

// Load reads and validates all configuration at startup; a missing variable
// crashes the process at boot, not at first use.
func Load() (Config, error) {
	upstream, err := parseUpstream(os.Getenv("UPSTREAM_URL"))
	if err != nil {
		return Config{}, err
	}
	cfg := Config{
		Port:               os.Getenv("PORT"),
		UpstreamURL:        upstream,
		KafkaBrokers:       os.Getenv("KAFKA_BROKERS"),
		RequestsTopic:      os.Getenv("REQUESTS_TOPIC"),
		DecisionsTopic:     os.Getenv("DECISIONS_TOPIC"),
		DecisionsGroup:     os.Getenv("DECISIONS_GROUP"),
		EventBufferSize:    1024,
		GateTTL:            time.Hour,
		RateLimitPerMinute: intFromEnv("RATE_LIMIT_PER_MINUTE", 60),
		RateLimitBurst:     intFromEnv("RATE_LIMIT_BURST", 10),
		ShutdownTimeout:    10 * time.Second,
	}
	return cfg, validate(cfg)
}

func intFromEnv(name string, fallback int) int {
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
	for name, v := range map[string]string{
		"PORT":            c.Port,
		"KAFKA_BROKERS":   c.KafkaBrokers,
		"REQUESTS_TOPIC":  c.RequestsTopic,
		"DECISIONS_TOPIC": c.DecisionsTopic,
		"DECISIONS_GROUP": c.DecisionsGroup,
	} {
		if v == "" {
			return fmt.Errorf("config: %s is required", name)
		}
	}
	return nil
}
