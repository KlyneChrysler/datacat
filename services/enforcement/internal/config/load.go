package config

import (
	"fmt"
	"os"
	"time"
)

// Load reads and validates all configuration at startup; a missing variable
// crashes the process at boot, not at first use.
func Load() (Config, error) {
	cfg := Config{
		Port:            os.Getenv("PORT"),
		KafkaBrokers:    os.Getenv("KAFKA_BROKERS"),
		VerdictsTopic:   os.Getenv("VERDICTS_TOPIC"),
		DecisionsTopic:  os.Getenv("DECISIONS_TOPIC"),
		ConsumerGroup:   os.Getenv("CONSUMER_GROUP"),
		DecisionsTable:  os.Getenv("DECISIONS_TABLE"),
		DynamoEndpoint:  os.Getenv("DYNAMO_ENDPOINT"),
		DecisionTTL:     time.Hour,
		ShutdownTimeout: 10 * time.Second,
	}
	return cfg, validate(cfg)
}

func validate(c Config) error {
	for name, v := range map[string]string{
		"PORT":            c.Port,
		"KAFKA_BROKERS":   c.KafkaBrokers,
		"VERDICTS_TOPIC":  c.VerdictsTopic,
		"DECISIONS_TOPIC": c.DecisionsTopic,
		"CONSUMER_GROUP":  c.ConsumerGroup,
	} {
		if v == "" {
			return fmt.Errorf("config: %s is required", name)
		}
	}
	return nil
}
