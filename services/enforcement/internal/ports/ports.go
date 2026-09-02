// Package ports holds the interfaces enforcement consumes.
package ports

import (
	"context"

	"github.com/KlyneChrysler/datacat/services/enforcement/internal/domain"
)

type DecisionStore interface {
	Save(ctx context.Context, d domain.Decision) error
	FindBySession(ctx context.Context, sessionID string) (domain.Decision, error)
}

type ActionApplier interface {
	Apply(ctx context.Context, d domain.Decision) error
}

type DecisionPublisher interface {
	PublishDecision(ctx context.Context, d domain.Decision) error
}

type VerdictSource interface {
	// Consume handles each verdict until ctx ends.
	Consume(ctx context.Context, handle func(context.Context, domain.Verdict) error) error
}
