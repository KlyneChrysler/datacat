// Package app holds enforcement use cases.
package app

import (
	"context"
	"fmt"

	policy "github.com/KlyneChrysler/datacat/pkg/policy"
	"github.com/KlyneChrysler/datacat/services/enforcement/internal/ports"
)

// Enforcer turns verdicts into stored, published, applied decisions.
type Enforcer struct {
	policy    policy.Policy
	store     ports.DecisionStore
	publisher ports.DecisionPublisher
	applier   ports.ActionApplier
	tally     *Tally
}

func NewEnforcer(policy policy.Policy, store ports.DecisionStore, publisher ports.DecisionPublisher, applier ports.ActionApplier, tally *Tally) *Enforcer {
	return &Enforcer{policy: policy, store: store, publisher: publisher, applier: applier, tally: tally}
}

// HandleVerdict runs the full decision pipeline for one verdict.
func (e *Enforcer) HandleVerdict(ctx context.Context, v policy.Verdict) error {
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

func (e *Enforcer) Lookup(ctx context.Context, sessionID string) (policy.Decision, error) {
	return e.store.FindBySession(ctx, sessionID)
}

func (e *Enforcer) TrafficSummary(windowMinutes int) policy.TrafficSummary {
	return e.tally.Summary(windowMinutes)
}
