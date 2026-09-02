// Package config owns environment configuration.
package config

import "time"

type Config struct {
	Port           string
	KafkaBrokers   string
	VerdictsTopic  string
	DecisionsTopic string
	ConsumerGroup  string
	// DecisionsTable selects the store, empty means in memory.
	DecisionsTable string
	// DynamoEndpoint targets DynamoDB Local, empty means real AWS.
	DynamoEndpoint string
	// CORSOrigin enables browser access, empty disables cors.
	CORSOrigin      string
	DecisionTTL     time.Duration
	ShutdownTimeout time.Duration
}
