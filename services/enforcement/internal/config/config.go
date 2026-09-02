// Package config owns environment configuration: the shape here, loading
// and validation in load.go (twelve-factor III).
package config

import "time"

type Config struct {
	Port           string
	KafkaBrokers   string
	VerdictsTopic  string
	DecisionsTopic string
	ConsumerGroup  string
	// DecisionsTable selects the store: empty = in-memory (single replica
	// only), set = DynamoDB (horizontally scalable). Optional by design.
	DecisionsTable string
	// DynamoEndpoint overrides the DynamoDB target for DynamoDB Local;
	// empty = real AWS. Optional by design.
	DynamoEndpoint string
	// CORSOrigin enables browser access from the dashboard origin; empty
	// disables CORS. Optional by design.
	CORSOrigin      string
	DecisionTTL     time.Duration
	ShutdownTimeout time.Duration
}
