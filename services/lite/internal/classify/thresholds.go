package classify

import "github.com/KlyneChrysler/datacat/pkg/policy"

const (
	humanBelow         = 0.45
	abusiveFrom        = 0.75
	verifiedShareFloor = 0.9
)

// classifyScores mirrors the Flink thresholds, identity wins unless abusive.
func classifyScores(scores []float64, verifiedShare float64) (policy.Classification, float64) {
	average := averageOf(scores)

	if verifiedShare >= verifiedShareFloor && average < abusiveFrom {
		return policy.VerifiedBot, verifiedShare
	}
	if average < humanBelow {
		return policy.Human, confidenceOf(average)
	}
	if average < abusiveFrom {
		return policy.Unverified, confidenceOf(average)
	}

	return policy.Abusive, confidenceOf(average)
}

// confidenceOf is the distance from the nearest boundary, normalized.
func confidenceOf(average float64) float64 {
	toBoundary := min(abs(average-humanBelow), abs(average-abusiveFrom))
	widest := max(humanBelow, max(abusiveFrom-humanBelow, 1-abusiveFrom))

	return min(toBoundary/widest, 1.0)
}

func averageOf(values []float64) float64 {
	sum := 0.0
	for _, v := range values {
		sum += v
	}

	return sum / float64(len(values))
}

func abs(v float64) float64 {
	if v < 0 {
		return -v
	}

	return v
}
