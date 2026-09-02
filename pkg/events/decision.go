package events

import "time"

// Decision is one enforcement outcome, keyed by session.
type Decision struct {
	SessionID      string    `json:"session_id"`
	Timestamp      time.Time `json:"timestamp"`
	Classification string    `json:"classification"`
	Action         string    `json:"action"`
}
