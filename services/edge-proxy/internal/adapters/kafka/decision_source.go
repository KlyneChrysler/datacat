package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/KlyneChrysler/datacat/pkg/events"
	"github.com/KlyneChrysler/datacat/pkg/kafkax"
	"github.com/KlyneChrysler/datacat/services/edge-proxy/internal/ports"
)

// DecisionSource feeds decisions from Kafka to the gatekeeper.
type DecisionSource struct {
	consumer *kafkax.Consumer
}

var _ ports.DecisionSource = (*DecisionSource)(nil)

func NewDecisionSource(consumer *kafkax.Consumer) *DecisionSource {
	return &DecisionSource{consumer: consumer}
}

func (s *DecisionSource) Consume(ctx context.Context, handle func(events.Decision)) error {
	return s.consumer.Consume(ctx, func(_ context.Context, _, value []byte) error {
		var decision events.Decision
		if err := json.Unmarshal(value, &decision); err != nil {
			return fmt.Errorf("decode decision: %w", err)
		}

		handle(decision)
		return nil
	})
}
