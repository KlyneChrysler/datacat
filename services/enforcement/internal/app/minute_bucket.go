package app

import "github.com/KlyneChrysler/datacat/services/enforcement/internal/domain"

// minuteBucket is the Tally's stored state per minute (shape file,
// standards rule 2). minute identifies which absolute minute the counts
// belong to, so a reused ring slot from an hour ago is detected and reset.
type minuteBucket struct {
	minute int64
	counts map[domain.Classification]int64
}
