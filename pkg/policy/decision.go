package policy

// Decision is one immutable enforcement outcome for a session.
type Decision struct {
	SessionID string
	Class     Classification
	Action    Action
}
