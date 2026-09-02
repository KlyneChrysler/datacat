// Package config loads and validates all environment configuration at
// startup (twelve-factor III).
package config

import (
	"fmt"
	"net/url"
	"os"
	"time"
)

type Config struct {
	Port            string
	UpstreamURL     *url.URL
	KafkaBrokers    string
	RequestsTopic   string
	EventBufferSize int
	ShutdownTimeout time.Duration
}

func Load() (Config, error) {
	upstream, err := parseUpstream(os.Getenv("UPSTREAM_URL"))
	if err != nil {
		return Config{}, err
	}
	cfg := Config{
		Port:            os.Getenv("PORT"),
		UpstreamURL:     upstream,
		KafkaBrokers:    os.Getenv("KAFKA_BROKERS"),
		RequestsTopic:   os.Getenv("REQUESTS_TOPIC"),
		EventBufferSize: 1024,
		ShutdownTimeout: 10 * time.Second,
	}
	return cfg, cfg.validate()
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

func (c Config) validate() error {
	for name, v := range map[string]string{
		"PORT":           c.Port,
		"KAFKA_BROKERS":  c.KafkaBrokers,
		"REQUESTS_TOPIC": c.RequestsTopic,
	} {
		if v == "" {
			return fmt.Errorf("config: %s is required", name)
		}
	}
	return nil
}
