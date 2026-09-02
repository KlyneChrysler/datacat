// Package kafka holds the enforcement Kafka adapters.
package kafka

import (
	"context"

	"github.com/KlyneChrysler/datacat/pkg/kafkax"
	"github.com/KlyneChrysler/datacat/services/enforcement/internal/domain"
	"github.com/KlyneChrysler/datacat/services/enforcement/internal/ports"
)

// VerdictSource feeds verdicts from Kafka to the enforcer.
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
