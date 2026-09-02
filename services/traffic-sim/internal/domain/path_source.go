package domain

// PathSource yields the next path a persona visits, single goroutine use.
type PathSource interface {
	Next() string
}
