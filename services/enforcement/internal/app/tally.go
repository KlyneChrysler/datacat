package app

import (
	"sync"
	"time"

	"github.com/KlyneChrysler/datacat/services/enforcement/internal/domain"
)

// tallyBuckets is one hour of minute buckets — the maximum summary window.
const tallyBuckets = 60

// Tally aggregates classification events into a fixed ring of per-minute
// buckets. Record is O(1) (hot path); Summary is O(tallyBuckets), bounded
// by the constant 60 (cold path). Memory is bounded by 60 buckets × the
// number of distinct classifications.
type Tally struct {
	mu      sync.Mutex
	buckets [tallyBuckets]minuteBucket
}

func NewTally() *Tally {
	return &Tally{}
}

// Record counts one classification event in the current minute.
func (t *Tally) Record(class domain.Classification) {
	t.recordAt(class, time.Now())
}

// Summary reports the last windowMinutes of events (clamped to [1,60]).
func (t *Tally) Summary(windowMinutes int) domain.TrafficSummary {
	return t.summaryAt(windowMinutes, time.Now())
}

func (t *Tally) recordAt(class domain.Classification, now time.Time) {
	minute := now.Unix() / 60
	t.mu.Lock()
	defer t.mu.Unlock()
	bucket := &t.buckets[minute%tallyBuckets]
	if bucket.minute != minute {
		bucket.minute = minute
		bucket.counts = make(map[domain.Classification]int64, 4)
	}
	bucket.counts[class]++
}

func (t *Tally) summaryAt(windowMinutes int, now time.Time) domain.TrafficSummary {
	window := clampWindow(windowMinutes)
	oldest := now.Unix()/60 - int64(window) + 1
	summary := domain.TrafficSummary{WindowMinutes: window}
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

func addBucket(summary *domain.TrafficSummary, bucket minuteBucket) {
	for class, count := range bucket.counts {
		summary.Total += count
		switch class {
		case domain.Human:
			summary.Human += count
		case domain.VerifiedBot:
			summary.VerifiedAgent += count
		case domain.Unverified:
			summary.Unverified += count
		case domain.Abusive:
			summary.Abusive += count
		default:
			summary.Other += count
		}
	}
}
