package events

import "time"

// Verdict is emitted by classifier-job for every classification change.
// Partition key: SessionID.
type Verdict struct {
	SessionID      string    `json:"session_id"`
	Timestamp      time.Time `json:"timestamp"`
	Classification string    `json:"classification"` // human | verified_agent | unverified_automation | abusive
	Confidence     float64   `json:"confidence"`     // [0,1]
}
