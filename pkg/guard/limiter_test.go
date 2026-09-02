package guard

import (
	"fmt"
	"testing"
	"time"
)

func TestLimiterAllowsBurstThenThrottles(t *testing.T) {
	limiter := NewRateLimiter(60, 3, time.Hour)
	now := time.Unix(7_000_000, 0)

	for i := range 3 {
		if !limiter.allowAt("s-1", now) {
			t.Fatalf("request %d within burst was throttled", i+1)
		}
	}
	if limiter.allowAt("s-1", now) {
		t.Error("request beyond burst was allowed")
	}
}

func TestLimiterRefillsOverTime(t *testing.T) {
	limiter := NewRateLimiter(60, 1, time.Hour) // one token per second
	now := time.Unix(7_000_000, 0)

	if !limiter.allowAt("s-1", now) {
		t.Fatal("first request throttled")
	}
	if limiter.allowAt("s-1", now) {
		t.Fatal("empty bucket allowed a request")
	}
	if !limiter.allowAt("s-1", now.Add(time.Second)) {
		t.Error("bucket did not refill after one second")
	}
}

func TestLimiterIsolatesSessions(t *testing.T) {
	limiter := NewRateLimiter(60, 1, time.Hour)
	now := time.Unix(7_000_000, 0)

	if !limiter.allowAt("s-1", now) {
		t.Fatal("first session throttled")
	}
	if !limiter.allowAt("s-2", now) {
		t.Error("second session throttled by first session's bucket")
	}
}

func TestLimiterFailsClosedAtCapAndRecoversAfterSweep(t *testing.T) {
	limiter := NewRateLimiter(60, 1, time.Minute)
	now := time.Unix(7_000_000, 0)
	for i := range maxTrackedSessions {
		limiter.allowAt(fmt.Sprintf("fill-%d", i), now)
	}

	if limiter.allowAt("overflow", now) {
		t.Error("new session allowed while limiter at cap (must fail closed)")
	}
	if !limiter.allowAt("overflow", now.Add(2*time.Minute)) {
		t.Error("sweep did not reclaim idle buckets after TTL")
	}
}
