// Package kafka is the inbound verdict adapter: Kafka records → domain
// Verdicts. Wire-format knowledge (pkg/events JSON) stays here.
package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/KlyneChrysler/datacat/pkg/events"
	"github.com/KlyneChrysler/datacat/pkg/kafkax"
	"github.com/KlyneChrysler/datacat/services/enforcement/internal/domain"
	"github.com/KlyneChrysler/datacat/services/enforcement/internal/ports"
)

type VerdictSource struct {
	consumer *kafkax.Consumer
}

var _ ports.VerdictSource = (*VerdictSource)(nil)

func NewVerdictSource(consumer *kafkax.Consumer) *VerdictSource {
	return &VerdictSource{consumer: consumer}
}

func (s *VerdictSource) Consume(ctx context.Context, handle func(context.Context, domain.Verdict) error) error {
	return s.consumer.Consume(ctx, func(ctx context.Context, _, value []byte) error {
		verdict, err := decodeVerdict(value)
		if err != nil {
			return err
		}
		return handle(ctx, verdict)
	})
}

func decodeVerdict(value []byte) (domain.Verdict, error) {
	var wire events.Verdict
	if err := json.Unmarshal(value, &wire); err != nil {
		return domain.Verdict{}, fmt.Errorf("decode verdict: %w", err)
	}
	verdict, err := domain.NewVerdict(wire.SessionID, domain.Classification(wire.Classification), wire.Confidence)
	if err != nil {
		return domain.Verdict{}, fmt.Errorf("invalid verdict for session %q: %w", wire.SessionID, err)
	}
	return verdict, nil
}
