package domain

import "math/rand/v2"

// LoopingPaths revisits a small path set in random order, human style.
type LoopingPaths struct {
	paths []string
}

func NewLoopingPaths(paths []string) *LoopingPaths {
	return &LoopingPaths{paths: paths}
}

func (l *LoopingPaths) Next() string {
	return l.paths[rand.IntN(len(l.paths))]
}
