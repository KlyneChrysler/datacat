package app

import "github.com/KlyneChrysler/datacat/pkg/policy"

// minuteBucket is the stored counts for one absolute minute.
type minuteBucket struct {
	minute int64
	counts map[policy.Classification]int64
}
