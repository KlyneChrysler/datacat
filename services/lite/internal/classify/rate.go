package classify

import "math"

// rateScorer reads rates far above human browsing as automation.
type rateScorer struct{}

func (rateScorer) Score(f Features) float64 {
	const certainlyAutomatedRPM = 120.0

	return math.Min(f.RequestsPerMinute/certainlyAutomatedRPM, 1.0)
}
