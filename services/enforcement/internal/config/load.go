package config

import (
	"os"
	"time"

	"github.com/KlyneChrysler/datacat/pkg/envx"
)

// Load reads and validates all configuration, crashing on missing values.
func Load() (Config, error) {
	cfg := Config{
		Port:            os.Getenv("PORT"),
		KafkaBrokers:    os.Getenv("KAFKA_BROKERS"),
		VerdictsTopic:   os.Getenv("VERDICTS_TOPIC"),
		DecisionsTopic:  os.Getenv("DECISIONS_TOPIC"),
		ConsumerGroup:   os.Getenv("CONSUMER_GROUP"),
		DecisionsTable:  os.Getenv("DECISIONS_TABLE"),
		DynamoEndpoint:  os.Getenv("DYNAMO_ENDPOINT"),
		CORSOrigin:      os.Getenv("CORS_ORIGIN"),
		DecisionTTL:     time.Hour,
		ShutdownTimeout: 10 * time.Second,
	}

	return cfg, validate(cfg)
}

func validate(c Config) error {
	return envx.Require(map[string]string{
		"PORT":            c.Port,
		"KAFKA_BROKERS":   c.KafkaBrokers,
		"VERDICTS_TOPIC":  c.VerdictsTopic,
		"DECISIONS_TOPIC": c.DecisionsTopic,
		"CONSUMER_GROUP":  c.ConsumerGroup,
	})
}
