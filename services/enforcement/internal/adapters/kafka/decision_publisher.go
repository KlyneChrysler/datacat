package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/KlyneChrysler/datacat/pkg/events"
	"github.com/KlyneChrysler/datacat/pkg/kafkax"
	"github.com/KlyneChrysler/datacat/services/enforcement/internal/domain"
	"github.com/KlyneChrysler/datacat/services/enforcement/internal/ports"
)

// DecisionPublisher is the outbound decision adapter: domain Decisions →
// Kafka, partitioned by session ID.
type DecisionPublisher struct {
	producer *kafkax.Producer
	topic    string
}

var _ ports.DecisionPublisher = (*DecisionPublisher)(nil)

func NewDecisionPublisher(producer *kafkax.Producer, topic string) *DecisionPublisher {
	return &DecisionPublisher{producer: producer, topic: topic}
}

func (p *DecisionPublisher) PublishDecision(ctx context.Context, d domain.Decision) error {
	payload, err := json.Marshal(toWire(d))
	if err != nil {
		return fmt.Errorf("marshal decision: %w", err)
	}
	return p.producer.Publish(ctx, p.topic, []byte(d.SessionID), payload)
}

func toWire(d domain.Decision) events.Decision {
	return events.Decision{
		SessionID:      d.SessionID,
		Timestamp:      time.Now().UTC(),
		Classification: string(d.Class),
		Action:         string(d.Action),
	}
}
