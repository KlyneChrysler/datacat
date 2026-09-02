package classify

import "math"

// entropyScorer reads never revisiting paths as crawling, damped when small.
type entropyScorer struct{}

func (entropyScorer) Score(f Features) float64 {
	const fullWeightRequests = 20.0

	volumeWeight := math.Min(float64(f.RequestCount)/fullWeightRequests, 1.0)

	return f.PathEntropy * volumeWeight
}
