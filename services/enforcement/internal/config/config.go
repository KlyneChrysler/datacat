// Package config owns environment configuration: the shape here, loading
// and validation in load.go (twelve-factor III).
package config

import "time"

type Config struct {
	Port            string
	KafkaBrokers    string
	VerdictsTopic   string
	DecisionsTopic  string
	ConsumerGroup   string
	ShutdownTimeout time.Duration
}
