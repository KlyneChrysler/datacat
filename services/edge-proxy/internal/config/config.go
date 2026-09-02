// Package config owns environment configuration.
package config

import (
	"crypto/ed25519"
	"net/url"
	"time"
)

type Config struct {
	Port            string
	UpstreamURL     *url.URL
	KafkaBrokers    string
	RequestsTopic   string
	DecisionsTopic  string
	DecisionsGroup  string
	EventBufferSize int
	GateTTL         time.Duration
	// Optional rate limit tuning.
	RateLimitPerMinute int
	RateLimitBurst     int
	// ChallengeSecret comes from a Secret in real deployments.
	ChallengeSecret string
	// ChallengeDifficulty is the proof of work bit count.
	ChallengeDifficulty int
	// AgentKeys registers trusted agent public keys, empty means none.
	AgentKeys       map[string]ed25519.PublicKey
	ShutdownTimeout time.Duration
}
