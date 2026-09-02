// Package kafkax holds the shared Kafka producer and consumer wrappers.
package kafkax

import (
	"context"
	"fmt"
	"strings"

	"github.com/twmb/franz-go/pkg/kgo"
)

// Producer publishes keyed records to one broker set.
type Producer struct {
	client *kgo.Client
}

func NewProducer(brokers string) (*Producer, error) {
	client, err := kgo.NewClient(kgo.SeedBrokers(splitBrokers(brokers)...))
	if err != nil {
		return nil, fmt.Errorf("kafka producer: %w", err)
	}

	return &Producer{client: client}, nil
}

func (p *Producer) Publish(ctx context.Context, topic string, key, value []byte) error {
	record := &kgo.Record{Topic: topic, Key: key, Value: value}

	if err := p.client.ProduceSync(ctx, record).FirstErr(); err != nil {
		return fmt.Errorf("produce to %s: %w", topic, err)
	}

	return nil
}

func (p *Producer) Close() {
	p.client.Close()
}

func splitBrokers(brokers string) []string {
	return strings.Split(brokers, ",")
}
