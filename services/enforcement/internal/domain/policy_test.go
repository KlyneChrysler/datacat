package domain

import "testing"

func TestDefaultPolicyDecide(t *testing.T) {
	tests := []struct {
		name       string
		class      Classification
		wantAction Action
	}{
		{name: "humans are allowed", class: Human, wantAction: Allow},
		{name: "verified agents are rate limited", class: VerifiedBot, wantAction: RateLimit},
		{name: "unverified automation is challenged", class: Unverified, wantAction: Challenge},
		{name: "abusive sessions are blocked", class: Abusive, wantAction: Block},
		{name: "unknown classification falls back to challenge", class: "mystery", wantAction: Challenge},
	}
	policy := DefaultPolicy()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verdict, err := NewVerdict("s-1", tt.class, 0.9)
			if err != nil {
				t.Fatal(err)
			}

			decision := policy.Decide(verdict)

			if decision.Action != tt.wantAction {
				t.Errorf("action = %s, want %s", decision.Action, tt.wantAction)
			}
			if decision.SessionID != "s-1" {
				t.Errorf("session id = %s, want s-1", decision.SessionID)
			}
		})
	}
}
