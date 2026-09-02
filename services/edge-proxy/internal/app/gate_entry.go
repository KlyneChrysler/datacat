package app

import "time"

// gateEntry is the Gatekeeper's stored state per session (shape file,
// standards rule 2).
type gateEntry struct {
	action  string
	savedAt time.Time
}
