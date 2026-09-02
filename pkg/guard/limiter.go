package guard

import (
	"sync"
	"time"
)

// maxTrackedSessions caps limiter memory, overflow fails closed.
const maxTrackedSessions = 10_000

// RateLimiter runs one token bucket per rate limited session.
type RateLimiter struct {
	ratePerSec float64
	burst      float64
	idleTTL    time.Duration
	mu         sync.Mutex
	buckets    map[string]*tokenBucket
}

func NewRateLimiter(perMinute int, burst int, idleTTL time.Duration) *RateLimiter {
	return &RateLimiter{ratePerSec: float64(perMinute) / 60.0, burst: float64(burst), idleTTL: idleTTL, buckets: make(map[string]*tokenBucket)}
}

// Allow consumes one token and reports whether the session may pass.
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
			return false
		}
	}

	refill(bucket, now, l.ratePerSec, l.burst)
	if bucket.tokens < 1 {
		return false
	}

	bucket.tokens--
	return true
}

// track registers a session, sweeping stale buckets when full.
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

// sweep drops buckets idle past the ttl, runs only at the cap.
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
