package dynamo

import (
	"time"

	"github.com/KlyneChrysler/datacat/services/enforcement/internal/domain"
)

// Domain to storage conversion.

func encodeDecision(d domain.Decision, now time.Time, ttl time.Duration) decisionRecord {
	return decisionRecord{SessionID: d.SessionID, Class: string(d.Class), Action: string(d.Action), ExpiresAt: now.Add(ttl).Unix()}
}

func decodeDecision(rec decisionRecord) domain.Decision {
	return domain.Decision{SessionID: rec.SessionID, Class: domain.Classification(rec.Class), Action: domain.Action(rec.Action)}
}

// expired guards reads because DynamoDB expires items lazily.
func expired(rec decisionRecord, now time.Time) bool {
	return now.Unix() >= rec.ExpiresAt
}
