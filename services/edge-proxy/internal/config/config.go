// Package config owns environment configuration: the shape here, loading
// and validation in load.go, derived values in derive.go (twelve-factor III).
package config

import (
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
	// Token bucket for sessions under a rate_limit decision. Optional envs
	// RATE_LIMIT_PER_MINUTE / RATE_LIMIT_BURST override the defaults.
	RateLimitPerMinute int
	RateLimitBurst     int
	ShutdownTimeout    time.Duration
}
