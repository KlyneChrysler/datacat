package domain

import (
	"testing"
	"time"
)

func TestNextDelayStaysWithinBounds(t *testing.T) {
	persona := Persona{BaseDelay: 100 * time.Millisecond, Jitter: 50 * time.Millisecond}

	for range 100 {
		delay := persona.NextDelay()
		if delay < 100*time.Millisecond || delay >= 150*time.Millisecond {
			t.Fatalf("delay %v outside [base, base+jitter)", delay)
		}
	}
}

func TestNextDelayWithoutJitterIsExact(t *testing.T) {
	persona := Persona{BaseDelay: 100 * time.Millisecond}

	if got := persona.NextDelay(); got != 100*time.Millisecond {
		t.Errorf("delay = %v, want exactly base", got)
	}
}

func TestNextRequestCarriesIdentity(t *testing.T) {
	persona := Persona{SessionID: "s-1", UserAgent: "ua/1", Paths: NewCrawlingPaths("/p-")}

	req := persona.NextRequest()

	if req.SessionID != "s-1" || req.UserAgent != "ua/1" || req.Path != "/p-0" {
		t.Errorf("request = %+v", req)
	}
}

func TestDefaultPersonasAreDistinctSessions(t *testing.T) {
	personas := DefaultPersonas()
	sessions := make(map[string]bool, len(personas))

	for _, p := range personas {
		if sessions[p.SessionID] {
			t.Fatalf("duplicate session id %q", p.SessionID)
		}
		sessions[p.SessionID] = true
	}
	if len(personas) != 3 {
		t.Errorf("personas = %d, want 3", len(personas))
	}
}
