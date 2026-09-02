// Package app holds lite use cases.
package app

import (
	"log/slog"

	"github.com/KlyneChrysler/datacat/pkg/guard"
	policy "github.com/KlyneChrysler/datacat/pkg/policy"
	"github.com/KlyneChrysler/datacat/services/lite/internal/classify"
)

// Analyzer turns observations into gate actions, all in process.
type Analyzer struct {
	tracker    *classify.Tracker
	policy     policy.Policy
	gatekeeper *guard.Gatekeeper
	log        *slog.Logger
}

func NewAnalyzer(tracker *classify.Tracker, pol policy.Policy, gatekeeper *guard.Gatekeeper, log *slog.Logger) *Analyzer {
	return &Analyzer{tracker: tracker, policy: pol, gatekeeper: gatekeeper, log: log}
}

// Observe records one request and applies a decision on class change.
func (a *Analyzer) Observe(o classify.Observation) {
	verdict, changed := a.tracker.Observe(o)
	if !changed {
		return
	}

	decision := a.policy.Decide(verdict)
	a.gatekeeper.Update(decision.SessionID, string(decision.Action))

	a.log.Info("session classified", "session_id", decision.SessionID, "class", decision.Class, "action", decision.Action)
}
