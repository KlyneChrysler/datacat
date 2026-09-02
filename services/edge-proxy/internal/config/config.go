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
	ShutdownTimeout time.Duration
}
