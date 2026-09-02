package classify

import (
	"math"
	"testing"
	"time"
)

func at(secs float64) time.Time {
	return time.Unix(9_000_000, 0).Add(time.Duration(secs * float64(time.Second)))
}

func TestIntervalCvRegularIsNearZero(t *testing.T) {
	samples := []time.Time{at(0), at(0.5), at(1.0), at(1.5), at(2.0)}

	if cv := intervalCv(samples); cv > 0.001 {
		t.Errorf("cv = %f, want near zero for metronome", cv)
	}
}

func TestIntervalCvJitteryIsHigh(t *testing.T) {
	samples := []time.Time{at(0), at(0.2), at(4.2), at(4.5), at(30)}

	if cv := intervalCv(samples); cv < 0.8 {
		t.Errorf("cv = %f, want high for jitter", cv)
	}
}

func TestIntervalCvTooFewIsNaN(t *testing.T) {
	if !math.IsNaN(intervalCv([]time.Time{at(0), at(1)})) {
		t.Error("two samples must give NaN")
	}
}

func TestEntropyNeverRevisitingIsOne(t *testing.T) {
	counts := map[string]int64{"/a": 1, "/b": 1, "/c": 1, "/d": 1}

	if e := normalizedEntropy(counts, 4); math.Abs(e-1.0) > 1e-9 {
		t.Errorf("entropy = %f, want 1", e)
	}
}

func TestEntropySinglePathIsZero(t *testing.T) {
	counts := map[string]int64{"/login": 3}

	if e := normalizedEntropy(counts, 3); e != 0 {
		t.Errorf("entropy = %f, want 0", e)
	}
}
