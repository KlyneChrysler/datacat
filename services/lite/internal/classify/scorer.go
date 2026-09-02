package classify

// Scorer judges one signal, zero human to one automated.
type Scorer interface {
	Score(f Features) float64
}
