// Package domain holds pure business logic. It imports zero infrastructure.
package domain

type Classification string

const (
	Human       Classification = "human"
	VerifiedBot Classification = "verified_agent"
	Unverified  Classification = "unverified_automation"
	Abusive     Classification = "abusive"
)

// Verdict is immutable: constructed once, never mutated. "Changes" return a
// new value (see WithClass).
type Verdict struct {
	SessionID  string
	Class      Classification
	Confidence float64
}

func NewVerdict(sessionID string, class Classification, confidence float64) (Verdict, error) {
	if sessionID == "" {
		return Verdict{}, ErrEmptySessionID
	}
	if confidence < 0 || confidence > 1 {
		return Verdict{}, ErrConfidenceRange
	}
	return Verdict{SessionID: sessionID, Class: class, Confidence: confidence}, nil
}

func (v Verdict) WithClass(c Classification) Verdict {
	v.Class = c // value receiver: mutates the copy, not the original
	return v
}
