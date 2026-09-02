package config

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/KlyneChrysler/datacat/pkg/envx"
)

// Load reads and validates all configuration, crashing on missing values.
func Load() (Config, error) {
	upstream, err := parseUpstream(os.Getenv("UPSTREAM_URL"))
	if err != nil {
		return Config{}, err
	}
	agentKeys, err := parseAgentKeys(os.Getenv("AGENT_KEYS"))
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		Port:                os.Getenv("PORT"),
		UpstreamURL:         upstream,
		KafkaBrokers:        os.Getenv("KAFKA_BROKERS"),
		RequestsTopic:       os.Getenv("REQUESTS_TOPIC"),
		DecisionsTopic:      os.Getenv("DECISIONS_TOPIC"),
		DecisionsGroup:      os.Getenv("DECISIONS_GROUP"),
		EventBufferSize:     1024,
		GateTTL:             time.Hour,
		RateLimitPerMinute:  envx.Int("RATE_LIMIT_PER_MINUTE", 60),
		RateLimitBurst:      envx.Int("RATE_LIMIT_BURST", 10),
		ChallengeSecret:     os.Getenv("CHALLENGE_SECRET"),
		ChallengeDifficulty: envx.Int("CHALLENGE_DIFFICULTY_BITS", 16),
		AgentKeys:           agentKeys,
		ShutdownTimeout:     10 * time.Second,
	}

	return cfg, validate(cfg)
}

// parseAgentKeys reads keyid=hexpubkey pairs, failing on any bad entry.
func parseAgentKeys(raw string) (map[string]ed25519.PublicKey, error) {
	keys := make(map[string]ed25519.PublicKey)
	if raw == "" {
		return keys, nil
	}

	for _, entry := range strings.Split(raw, ";") {
		keyID, hexKey, found := strings.Cut(entry, "=")
		if !found {
			return nil, fmt.Errorf("config: AGENT_KEYS entry %q is not keyid=hexpubkey", entry)
		}
		pub, err := hex.DecodeString(hexKey)
		if err != nil || len(pub) != ed25519.PublicKeySize {
			return nil, fmt.Errorf("config: AGENT_KEYS %q has an invalid public key", keyID)
		}
		keys[keyID] = ed25519.PublicKey(pub)
	}

	return keys, nil
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
		"PORT":             c.Port,
		"KAFKA_BROKERS":    c.KafkaBrokers,
		"REQUESTS_TOPIC":   c.RequestsTopic,
		"DECISIONS_TOPIC":  c.DecisionsTopic,
		"DECISIONS_GROUP":  c.DecisionsGroup,
		"CHALLENGE_SECRET": c.ChallengeSecret,
	} {
		if v == "" {
			return fmt.Errorf("config: %s is required", name)
		}
	}

	return nil
}
