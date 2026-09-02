package classify

import (
	"fmt"
	"testing"
	"time"

	policy "github.com/KlyneChrysler/datacat/pkg/policy"
)

func drive(t *testing.T, tracker *Tracker, session, ua string, gap time.Duration, n int, verified bool) (policy.Verdict, bool) {
	t.Helper()

	now := time.Unix(9_000_000, 0)
	var verdict policy.Verdict
	changed := false

	for i := 0; i < n; i++ {
		o := Observation{SessionID: session, Path: fmt.Sprintf("/p-%d", i), UserAgent: ua, Verified: verified, At: now}
		if v, ok := tracker.Observe(o); ok {
			verdict, changed = v, true
		}
		now = now.Add(gap)
	}

	return verdict, changed
}

func TestScraperPatternClassifiesAutomated(t *testing.T) {
	verdict, changed := drive(t, NewTracker(), "scraper", "python-requests/2.32", 400*time.Millisecond, 40, false)

	if !changed {
		t.Fatal("no verdict emitted")
	}
	if verdict.Class != policy.Abusive && verdict.Class != policy.Unverified {
		t.Errorf("class = %s, want automated", verdict.Class)
	}
}

func TestSignedAgentClassifiesVerified(t *testing.T) {
	verdict, changed := drive(t, NewTracker(), "agent", "datacat-agent/1.0", 800*time.Millisecond, 40, true)

	if !changed {
		t.Fatal("no verdict emitted")
	}
	if verdict.Class != policy.VerifiedBot {
		t.Errorf("class = %s, want verified_agent", verdict.Class)
	}
}

func TestEvaluationIsThrottled(t *testing.T) {
	tracker := NewTracker()
	now := time.Unix(9_000_000, 0)

	first := 0
	for i := 0; i < 10; i++ {
		o := Observation{SessionID: "s", Path: "/x", UserAgent: "curl", At: now.Add(time.Duration(i) * 10 * time.Millisecond)}
		if _, ok := tracker.Observe(o); ok {
			first++
		}
	}

	if first > 1 {
		t.Errorf("evaluations in 100ms = %d, want at most 1", first)
	}
}

func TestVerdictOnlyOnClassChange(t *testing.T) {
	tracker := NewTracker()

	_, changedFirst := drive(t, tracker, "steady", "curl", 3*time.Second, 20, false)
	_, changedAgain := drive(t, tracker, "steady", "curl", 3*time.Second, 5, false)

	if !changedFirst {
		t.Fatal("first classification never emitted")
	}
	if changedAgain {
		t.Error("unchanged class re-emitted a verdict")
	}
}
