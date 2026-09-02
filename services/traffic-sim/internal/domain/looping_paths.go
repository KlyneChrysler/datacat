package domain

import "math/rand/v2"

// LoopingPaths revisits a small fixed set of paths in random order —
// human-like navigation. Next is O(1).
type LoopingPaths struct {
	paths []string
}

func NewLoopingPaths(paths []string) *LoopingPaths {
	return &LoopingPaths{paths: paths}
}

func (l *LoopingPaths) Next() string {
	return l.paths[rand.IntN(len(l.paths))]
}
