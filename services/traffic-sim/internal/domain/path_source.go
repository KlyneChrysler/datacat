package domain

// PathSource yields the next path a persona visits. Implementations keep
// per-persona state and are used from a single goroutine — not concurrent.
type PathSource interface {
	Next() string
}
