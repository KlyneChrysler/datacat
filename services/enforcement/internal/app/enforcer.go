// Package app holds enforcement use cases.
package app

import (
	"context"
	"fmt"

	"github.com/KlyneChrysler/datacat/services/enforcement/internal/domain"
	"github.com/KlyneChrysler/datacat/services/enforcement/internal/ports"
)

type Enforcer struct {
	policy    domain.Policy
	store     ports.DecisionStore
	publisher ports.DecisionPublisher
	applier   ports.ActionApplier
	tally     *Tally
}

func NewEnforcer(policy domain.Policy, store ports.DecisionStore,
	publisher ports.DecisionPublisher, applier ports.ActionApplier, tally *Tally) *Enforcer {
	return &Enforcer{policy: policy, store: store, publisher: publisher, applier: applier, tally: tally}
}

// HandleVerdict is an orchestrator: every line delegates.
func (e *Enforcer) HandleVerdict(ctx context.Context, v domain.Verdict) error {
	e.tally.Record(v.Class)
	decision := e.policy.Decide(v)
	if err := e.store.Save(ctx, decision); err != nil {
		return fmt.Errorf("save decision for session %s: %w", v.SessionID, err)
	}
	if err := e.publisher.PublishDecision(ctx, decision); err != nil {
		return fmt.Errorf("publish decision for session %s: %w", v.SessionID, err)
	}
	if err := e.applier.Apply(ctx, decision); err != nil {
		return fmt.Errorf("apply %s for session %s: %w", decision.Action, v.SessionID, err)
	}
	return nil
}

func (e *Enforcer) Lookup(ctx context.Context, sessionID string) (domain.Decision, error) {
	return e.store.FindBySession(ctx, sessionID)
}

func (e *Enforcer) TrafficSummary(windowMinutes int) domain.TrafficSummary {
	return e.tally.Summary(windowMinutes)
}
