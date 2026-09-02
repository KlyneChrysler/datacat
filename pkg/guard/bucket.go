package guard

import "time"

// tokenBucket is the stored rate limit state for one session.
type tokenBucket struct {
	tokens     float64
	lastRefill time.Time
}
