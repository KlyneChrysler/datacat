package domain

// TrafficSummary counts classification events over a recent window.
type TrafficSummary struct {
	WindowMinutes int
	Human         int64
	VerifiedAgent int64
	Unverified    int64
	Abusive       int64
	Other         int64
	Total         int64
}
