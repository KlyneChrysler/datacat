package guard

import (
	"testing"
	"time"
)

func TestGatekeeperReturnsLatestDecision(t *testing.T) {
	gk := NewGatekeeper(time.Hour)

	gk.Update("s-1", "challenge")
	gk.Update("s-1", "block")

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
	gk.Update("s-1", "block")

	time.Sleep(5 * time.Millisecond)

	if got := gk.ActionFor("s-1"); got != "allow" {
		t.Fatalf("action = %q, want allow after expiry", got)
	}
}
