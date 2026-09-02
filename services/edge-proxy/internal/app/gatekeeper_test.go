package app

import (
	"testing"
	"time"

	"github.com/KlyneChrysler/datacat/pkg/events"
)

func TestGatekeeperReturnsLatestDecision(t *testing.T) {
	gk := NewGatekeeper(time.Hour)

	gk.Update(events.Decision{SessionID: "s-1", Action: "challenge"})
	gk.Update(events.Decision{SessionID: "s-1", Action: "block"})

	if got := gk.ActionFor("s-1"); got != "block" {
		t.Fatalf("action = %q, want block", got)
	}
}

func TestGatekeeperAllowsUnknownSessions(t *testing.T) {
	gk := NewGatekeeper(time.Hour)

	if got := gk.ActionFor("stranger"); got != "allow" {
		t.Fatalf("action = %q, want allow", got)
	}
}

func TestGatekeeperExpiresOldDecisions(t *testing.T) {
	gk := NewGatekeeper(time.Millisecond)
	gk.Update(events.Decision{SessionID: "s-1", Action: "block"})

	time.Sleep(5 * time.Millisecond)

	if got := gk.ActionFor("s-1"); got != "allow" {
		t.Fatalf("action = %q, want allow after expiry", got)
	}
}
