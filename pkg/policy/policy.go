package policy

// Policy maps classifications to actions.
type Policy struct {
	actions  map[Classification]Action
	fallback Action
}

// DefaultPolicy challenges anything unrecognized, never silently allows.
func DefaultPolicy() Policy {
	actions := map[Classification]Action{Human: Allow, VerifiedBot: RateLimit, Unverified: Challenge, Abusive: Block}

	return Policy{actions: actions, fallback: Challenge}
}

func (p Policy) Decide(v Verdict) Decision {
	action, known := p.actions[v.Class]
	if !known {
		action = p.fallback
	}

	return Decision{SessionID: v.SessionID, Class: v.Class, Action: action}
}
