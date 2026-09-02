package classify

import "math"

// timingScorer reads machine like regularity as automation.
type timingScorer struct{}

func (timingScorer) Score(f Features) float64 {
	if math.IsNaN(f.IntervalCv) {
		return 0.5
	}

	return 1.0 - math.Min(f.IntervalCv, 1.0)
}
