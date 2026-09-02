package app

import "github.com/KlyneChrysler/datacat/services/enforcement/internal/domain"

// minuteBucket is the stored counts for one absolute minute.
type minuteBucket struct {
	minute int64
	counts map[domain.Classification]int64
}
