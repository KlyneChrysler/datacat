package dynamo

import (
	"testing"
	"time"

	"github.com/KlyneChrysler/datacat/services/enforcement/internal/domain"
)

func TestCodecRoundTrip(t *testing.T) {
	decision := domain.Decision{SessionID: "s-1", Class: domain.Abusive, Action: domain.Block}
	now := time.Unix(1_000_000, 0)

	rec := encodeDecision(decision, now, time.Hour)
	back := decodeDecision(rec)

	if back != decision {
		t.Errorf("round trip = %+v, want %+v", back, decision)
	}
	if rec.ExpiresAt != now.Add(time.Hour).Unix() {
		t.Errorf("expires_at = %d, want now+1h", rec.ExpiresAt)
	}
}

func TestExpiredEnforcesTTLOnRead(t *testing.T) {
	now := time.Unix(1_000_000, 0)
	rec := encodeDecision(domain.Decision{SessionID: "s-1"}, now, time.Hour)

	if expired(rec, now.Add(59*time.Minute)) {
		t.Error("record expired before its TTL")
	}
	if !expired(rec, now.Add(61*time.Minute)) {
		t.Error("record still live after its TTL")
	}
}
