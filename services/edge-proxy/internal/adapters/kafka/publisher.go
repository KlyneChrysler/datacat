// Package kafka holds the edge proxy Kafka adapters.
package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/KlyneChrysler/datacat/pkg/events"
	"github.com/KlyneChrysler/datacat/pkg/kafkax"
	"github.com/KlyneChrysler/datacat/services/edge-proxy/internal/ports"
)

// EventPublisher sends request events keyed by session.
type EventPublisher struct {
	producer *kafkax.Producer
	topic    string
}

var _ ports.EventPublisher = (*EventPublisher)(nil)

func NewEventPublisher(producer *kafkax.Producer, topic string) *EventPublisher {
	return &EventPublisher{producer: producer, topic: topic}
}

func (p *EventPublisher) PublishRequest(ctx context.Context, ev events.RequestEvent) error {
	payload, err := json.Marshal(ev)
	if err != nil {
		return fmt.Errorf("marshal request event: %w", err)
	}

	return p.producer.Publish(ctx, p.topic, []byte(ev.SessionID), payload)
}
