package app

import (
	"sync"
	"time"

	"github.com/KlyneChrysler/datacat/pkg/policy"
)

// tallyBuckets is one hour of minute buckets, the largest summary window.
const tallyBuckets = 60

// Tally counts classification events in a fixed ring of minute buckets.
type Tally struct {
	mu      sync.Mutex
	buckets [tallyBuckets]minuteBucket
}

func NewTally() *Tally {
	return &Tally{}
}

// Record counts one event in the current minute.
func (t *Tally) Record(class policy.Classification) {
	t.recordAt(class, time.Now())
}

// Summary reports the last windowMinutes of events, clamped to the ring.
func (t *Tally) Summary(windowMinutes int) policy.TrafficSummary {
	return t.summaryAt(windowMinutes, time.Now())
}

func (t *Tally) recordAt(class policy.Classification, now time.Time) {
	minute := now.Unix() / 60

	t.mu.Lock()
	defer t.mu.Unlock()

	bucket := &t.buckets[minute%tallyBuckets]
	if bucket.minute != minute {
		bucket.minute = minute
		bucket.counts = make(map[policy.Classification]int64, 4)
	}

	bucket.counts[class]++
}

func (t *Tally) summaryAt(windowMinutes int, now time.Time) policy.TrafficSummary {
	window := clampWindow(windowMinutes)
	oldest := now.Unix()/60 - int64(window) + 1
	summary := policy.TrafficSummary{WindowMinutes: window}

	t.mu.Lock()
	defer t.mu.Unlock()

	for i := range t.buckets {
		if t.buckets[i].minute >= oldest {
			addBucket(&summary, t.buckets[i])
		}
	}

	return summary
}

func clampWindow(minutes int) int {
	if minutes < 1 {
		return 1
	}
	if minutes > tallyBuckets {
		return tallyBuckets
	}

	return minutes
}

func addBucket(summary *policy.TrafficSummary, bucket minuteBucket) {
	for class, count := range bucket.counts {
		summary.Total += count

		switch class {
		case policy.Human:
			summary.Human += count
		case policy.VerifiedBot:
			summary.VerifiedAgent += count
		case policy.Unverified:
			summary.Unverified += count
		case policy.Abusive:
			summary.Abusive += count
		default:
			summary.Other += count
		}
	}
}
