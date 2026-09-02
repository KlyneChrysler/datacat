package classify

// defaultScorers is the signal registry, one line per new signal.
func defaultScorers() []Scorer {
	return []Scorer{timingScorer{}, rateScorer{}, entropyScorer{}, userAgentScorer{}}
}
