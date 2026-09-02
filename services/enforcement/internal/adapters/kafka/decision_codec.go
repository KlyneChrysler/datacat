package kafka

import (
	"time"

	"github.com/KlyneChrysler/datacat/pkg/events"
	policy "github.com/KlyneChrysler/datacat/pkg/policy"
)

// encodeDecision converts one domain decision into its wire form.
func encodeDecision(d policy.Decision) events.Decision {
	return events.Decision{SessionID: d.SessionID, Timestamp: time.Now().UTC(), Classification: string(d.Class), Action: string(d.Action)}
}
