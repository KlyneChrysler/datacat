package classify

// Features are the facts about one session, judgment lives in scorers.
type Features struct {
	RequestCount      int64
	RequestsPerMinute float64
	IntervalCv        float64
	PathEntropy       float64
	UserAgent         string
	VerifiedShare     float64
}
