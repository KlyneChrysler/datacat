// Package classify holds the in process session classifier.
package classify

import "time"

// Observation is one watched request.
type Observation struct {
	SessionID string
	Path      string
	UserAgent string
	Verified  bool
	At        time.Time
}
