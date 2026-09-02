package domain

// Decision is the enforcement outcome for one session. Immutable.
type Decision struct {
	SessionID string
	Class     Classification
	Action    Action
}
