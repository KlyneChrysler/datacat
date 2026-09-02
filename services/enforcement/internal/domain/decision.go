package domain

type Action string

const (
	Allow     Action = "allow"
	Challenge Action = "challenge"
	RateLimit Action = "rate_limit"
	Block     Action = "block"
)

// Decision is the enforcement outcome for one session. Immutable.
type Decision struct {
	SessionID string
	Class     Classification
	Action    Action
}

// Policy maps classifications to actions. The zero value is unusable on
// purpose — construct via DefaultPolicy so the unknown-class fallback is
// always defined.
type Policy struct {
	actions  map[Classification]Action
	fallback Action
}

// DefaultPolicy is deliberately conservative: anything unrecognized is
// challenged, never silently allowed.
func DefaultPolicy() Policy {
	return Policy{
		actions: map[Classification]Action{
			Human:       Allow,
			VerifiedBot: RateLimit,
			Unverified:  Challenge,
			Abusive:     Block,
		},
		fallback: Challenge,
	}
}

func (p Policy) Decide(v Verdict) Decision {
	action, known := p.actions[v.Class]
	if !known {
		action = p.fallback
	}
	return Decision{SessionID: v.SessionID, Class: v.Class, Action: action}
}
