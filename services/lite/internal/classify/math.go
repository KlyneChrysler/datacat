package classify

import (
	"math"
	"time"
)

// intervalCv is the variation of gaps between samples, NaN means no evidence.
func intervalCv(timestamps []time.Time) float64 {
	if len(timestamps) < 3 {
		return math.NaN()
	}

	intervals := make([]float64, 0, len(timestamps)-1)
	for i := 1; i < len(timestamps); i++ {
		intervals = append(intervals, timestamps[i].Sub(timestamps[i-1]).Seconds())
	}

	mean := meanOf(intervals)
	if mean == 0 {
		return 0
	}

	return math.Sqrt(varianceOf(intervals, mean)) / mean
}

// normalizedEntropy is one for never revisiting, zero for one path.
func normalizedEntropy(counts map[string]int64, total int64) float64 {
	if total < 2 {
		return 0
	}

	entropy := 0.0
	for _, count := range counts {
		p := float64(count) / float64(total)
		entropy -= p * math.Log2(p)
	}

	max := math.Log2(float64(total))
	if max == 0 {
		return 0
	}

	return math.Min(entropy/max, 1.0)
}

func meanOf(values []float64) float64 {
	sum := 0.0
	for _, v := range values {
		sum += v
	}

	return sum / float64(len(values))
}

func varianceOf(values []float64, mean float64) float64 {
	sum := 0.0
	for _, v := range values {
		sum += (v - mean) * (v - mean)
	}

	return sum / float64(len(values))
}
