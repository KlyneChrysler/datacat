package events

import "time"

// Decision is emitted by enforcement after applying policy to a verdict;
// edge-proxy consumes it to gate traffic. Partition key: SessionID.
type Decision struct {
	SessionID      string    `json:"session_id"`
	Timestamp      time.Time `json:"timestamp"`
	Classification string    `json:"classification"`
	Action         string    `json:"action"` // allow | challenge | rate_limit | block
}
