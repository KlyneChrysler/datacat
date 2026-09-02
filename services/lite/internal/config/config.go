// Package config owns environment configuration.
package config

import (
	"crypto/ed25519"
	"net/url"
	"time"
)

type Config struct {
	Port        string
	UpstreamURL *url.URL
	// ChallengeSecret comes from a Secret in real deployments.
	ChallengeSecret string
	// Optional tuning, all with working defaults.
	ChallengeDifficulty int
	RateLimitPerMinute  int
	RateLimitBurst      int
	AgentKeys           map[string]ed25519.PublicKey
	GateTTL             time.Duration
	ShutdownTimeout     time.Duration
}
