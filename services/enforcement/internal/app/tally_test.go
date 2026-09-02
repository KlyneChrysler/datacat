package app

import (
	"testing"
	"time"

	"github.com/KlyneChrysler/datacat/pkg/policy"
)

func TestTallyCountsWithinWindow(t *testing.T) {
	tally := NewTally()
	now := time.Unix(6_000_000, 0)

	tally.recordAt(policy.Human, now)
	tally.recordAt(policy.Human, now.Add(-2*time.Minute))
	tally.recordAt(policy.Abusive, now)

	summary := tally.summaryAt(15, now)

	if summary.Human != 2 || summary.Abusive != 1 || summary.Total != 3 {
		t.Errorf("summary = %+v, want human=2 abusive=1 total=3", summary)
	}
}

func TestTallyExcludesEventsOutsideWindow(t *testing.T) {
	tally := NewTally()
	now := time.Unix(6_000_000, 0)

	tally.recordAt(policy.Abusive, now.Add(-20*time.Minute))
	tally.recordAt(policy.Human, now)

	summary := tally.summaryAt(15, now)

	if summary.Abusive != 0 || summary.Human != 1 {
		t.Errorf("summary = %+v, want the 20-minute-old event excluded", summary)
	}
}

func TestTallyRingSlotReuseResetsStaleCounts(t *testing.T) {
	tally := NewTally()
	now := time.Unix(6_000_000, 0)

	tally.recordAt(policy.Human, now.Add(-time.Duration(tallyBuckets)*time.Minute))
	tally.recordAt(policy.Abusive, now) // same ring slot, one hour later

	summary := tally.summaryAt(tallyBuckets, now)

	if summary.Human != 0 || summary.Abusive != 1 {
		t.Errorf("summary = %+v, want stale slot fully replaced", summary)
	}
}

func TestTallyClampsWindow(t *testing.T) {
	tally := NewTally()
	now := time.Unix(6_000_000, 0)
	tally.recordAt(policy.Unverified, now)

	if got := tally.summaryAt(500, now).WindowMinutes; got != tallyBuckets {
		t.Errorf("window = %d, want clamped to %d", got, tallyBuckets)
	}
	if got := tally.summaryAt(0, now).WindowMinutes; got != 1 {
		t.Errorf("window = %d, want clamped to 1", got)
	}
}

func TestTallyCountsUnknownClassAsOther(t *testing.T) {
	tally := NewTally()
	now := time.Unix(6_000_000, 0)

	tally.recordAt(policy.Classification("mystery"), now)

	if summary := tally.summaryAt(5, now); summary.Other != 1 || summary.Total != 1 {
		t.Errorf("summary = %+v, want other=1 total=1", summary)
	}
}
