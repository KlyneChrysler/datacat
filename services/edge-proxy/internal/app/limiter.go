package app

import (
	"sync"
	"time"
)

// maxTrackedSessions bounds limiter memory. Overflow policy is fail-closed:
// a session that cannot be tracked is throttled — it is already under a
// rate_limit decision, so erring toward restriction is the safe direction.
const maxTrackedSessions = 10_000

// RateLimiter enforces a per-session token bucket for sessions whose
// standing action is rate_limit. Allow is O(1) (hot path); expired buckets
// are swept only when the cap is hit (amortized, bounded by the cap).
type RateLimiter struct {
	ratePerSec float64
	burst      float64
	idleTTL    time.Duration
	mu         sync.Mutex
	buckets    map[string]*tokenBucket
}

func NewRateLimiter(perMinute int, burst int, idleTTL time.Duration) *RateLimiter {
	return &RateLimiter{
		ratePerSec: float64(perMinute) / 60.0,
		burst:      float64(burst),
		idleTTL:    idleTTL,
		buckets:    make(map[string]*tokenBucket),
	}
}

// Allow reports whether the session may pass right now, consuming one token.
func (l *RateLimiter) Allow(sessionID string) bool {
	return l.allowAt(sessionID, time.Now())
}

func (l *RateLimiter) allowAt(sessionID string, now time.Time) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	bucket, ok := l.buckets[sessionID]
	if !ok {
		bucket = l.track(sessionID, now)
		if bucket == nil {
			return false // cap reached even after sweeping: fail closed
		}
	}
	refill(bucket, now, l.ratePerSec, l.burst)
	if bucket.tokens < 1 {
		return false
	}
	bucket.tokens--
	return true
}

// track registers a new session, sweeping expired buckets when at capacity.
func (l *RateLimiter) track(sessionID string, now time.Time) *tokenBucket {
	if len(l.buckets) >= maxTrackedSessions {
		l.sweep(now)
	}
	if len(l.buckets) >= maxTrackedSessions {
		return nil
	}
	bucket := &tokenBucket{tokens: l.burst, lastRefill: now}
	l.buckets[sessionID] = bucket
	return bucket
}

// sweep drops buckets idle past the TTL. O(tracked sessions), runs only at
// the cap — amortized against the inserts that filled the map.
func (l *RateLimiter) sweep(now time.Time) {
	for id, bucket := range l.buckets {
		if now.Sub(bucket.lastRefill) > l.idleTTL {
			delete(l.buckets, id)
		}
	}
}

func refill(bucket *tokenBucket, now time.Time, ratePerSec, burst float64) {
	elapsed := now.Sub(bucket.lastRefill).Seconds()
	bucket.tokens = min(burst, bucket.tokens+elapsed*ratePerSec)
	bucket.lastRefill = now
}
