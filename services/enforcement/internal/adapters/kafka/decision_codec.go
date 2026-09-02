package kafka

import (
	"time"

	"github.com/KlyneChrysler/datacat/pkg/events"
	"github.com/KlyneChrysler/datacat/services/enforcement/internal/domain"
)

// encodeDecision converts one domain decision into its wire form.
func encodeDecision(d domain.Decision) events.Decision {
	return events.Decision{SessionID: d.SessionID, Timestamp: time.Now().UTC(), Classification: string(d.Class), Action: string(d.Action)}
}
