package dynamo

import (
	"time"

	"github.com/KlyneChrysler/datacat/services/enforcement/internal/domain"
)

// Domain ↔ storage conversion only (file taxonomy, standards rule 2).

func encodeDecision(d domain.Decision, now time.Time, ttl time.Duration) decisionRecord {
	return decisionRecord{
		SessionID: d.SessionID,
		Class:     string(d.Class),
		Action:    string(d.Action),
		ExpiresAt: now.Add(ttl).Unix(),
	}
}

func decodeDecision(rec decisionRecord) domain.Decision {
	return domain.Decision{
		SessionID: rec.SessionID,
		Class:     domain.Classification(rec.Class),
		Action:    domain.Action(rec.Action),
	}
}

// expired reports whether the record has outlived its TTL. DynamoDB expires
// items lazily (up to days late), so reads must enforce the boundary too.
func expired(rec decisionRecord, now time.Time) bool {
	return now.Unix() >= rec.ExpiresAt
}
