package app

import "time"

// tokenBucket is the RateLimiter's stored state per session (shape file,
// standards rule 2). tokens refills continuously at the limiter's rate up
// to its burst capacity.
type tokenBucket struct {
	tokens     float64
	lastRefill time.Time
}
