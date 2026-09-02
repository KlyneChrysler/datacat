package domain

// TrafficSummary counts classification events over a recent window. A
// session escalating (e.g. unverified → abusive) contributes one event per
// transition — this summarizes the classifier's output stream, not raw
// request volume.
type TrafficSummary struct {
	WindowMinutes int
	Human         int64
	VerifiedAgent int64
	Unverified    int64
	Abusive       int64
	Other         int64
	Total         int64
}
