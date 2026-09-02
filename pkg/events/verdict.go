package events

import "time"

// Verdict is one classification change, keyed by session.
type Verdict struct {
	SessionID      string    `json:"session_id"`
	Timestamp      time.Time `json:"timestamp"`
	Classification string    `json:"classification"`
	Confidence     float64   `json:"confidence"`
}
