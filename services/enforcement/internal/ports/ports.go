// Package ports holds the interfaces enforcement consumes. Implementations
// live in adapters and meet these contracts only in the composition root.
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

type VerdictSource interface {
	// Consume blocks, invoking handle for each verdict until ctx is done.
	Consume(ctx context.Context, handle func(context.Context, domain.Verdict) error) error
}
