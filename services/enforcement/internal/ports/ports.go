// Package ports holds the interfaces enforcement consumes.
package ports

import (
	"context"

	policy "github.com/KlyneChrysler/datacat/pkg/policy"
)

type DecisionStore interface {
	Save(ctx context.Context, d policy.Decision) error
	FindBySession(ctx context.Context, sessionID string) (policy.Decision, error)
}

type ActionApplier interface {
	Apply(ctx context.Context, d policy.Decision) error
}

type DecisionPublisher interface {
	PublishDecision(ctx context.Context, d policy.Decision) error
}

type VerdictSource interface {
	// Consume handles each verdict until ctx ends.
	Consume(ctx context.Context, handle func(context.Context, policy.Verdict) error) error
}
