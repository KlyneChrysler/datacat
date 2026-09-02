package domain

// Verdict is one immutable classification of a session.
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

// WithClass returns a copy with a different class.
func (v Verdict) WithClass(c Classification) Verdict {
	v.Class = c

	return v
}
