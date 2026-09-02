package dynamo

import (
	"time"

	policy "github.com/KlyneChrysler/datacat/pkg/policy"
)

// Domain to storage conversion.

func encodeDecision(d policy.Decision, now time.Time, ttl time.Duration) decisionRecord {
	return decisionRecord{SessionID: d.SessionID, Class: string(d.Class), Action: string(d.Action), ExpiresAt: now.Add(ttl).Unix()}
}

func decodeDecision(rec decisionRecord) policy.Decision {
	return policy.Decision{SessionID: rec.SessionID, Class: policy.Classification(rec.Class), Action: policy.Action(rec.Action)}
}

// expired guards reads because DynamoDB expires items lazily.
func expired(rec decisionRecord, now time.Time) bool {
	return now.Unix() >= rec.ExpiresAt
}
