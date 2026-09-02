// Package config loads and validates all environment configuration at
// startup (twelve-factor III). Missing config crashes the process at boot.
package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	Port            string
	KafkaBrokers    string
	VerdictsTopic   string
	ConsumerGroup   string
	ShutdownTimeout time.Duration
}

func Load() (Config, error) {
	cfg := Config{
		Port:            os.Getenv("PORT"),
		KafkaBrokers:    os.Getenv("KAFKA_BROKERS"),
		VerdictsTopic:   os.Getenv("VERDICTS_TOPIC"),
		ConsumerGroup:   os.Getenv("CONSUMER_GROUP"),
		ShutdownTimeout: 10 * time.Second,
	}
	return cfg, cfg.validate()
}

func (c Config) validate() error {
	for name, v := range map[string]string{
		"PORT":           c.Port,
		"KAFKA_BROKERS":  c.KafkaBrokers,
		"VERDICTS_TOPIC": c.VerdictsTopic,
		"CONSUMER_GROUP": c.ConsumerGroup,
	} {
		if v == "" {
			return fmt.Errorf("config: %s is required", name)
		}
	}
	return nil
}
