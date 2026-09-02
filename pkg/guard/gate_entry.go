package guard

import "time"

// gateEntry is the stored decision state for one session.
type gateEntry struct {
	action  string
	savedAt time.Time
}
