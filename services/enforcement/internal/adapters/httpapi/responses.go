package httpapi

// Wire shapes for the HTTP API.

type decisionResponse struct {
	SessionID string `json:"session_id"`
	Class     string `json:"classification"`
	Action    string `json:"action"`
}

type trafficSummaryResponse struct {
	WindowMinutes int   `json:"window_minutes"`
	Human         int64 `json:"human"`
	VerifiedAgent int64 `json:"verified_agent"`
	Unverified    int64 `json:"unverified_automation"`
	Abusive       int64 `json:"abusive"`
	Other         int64 `json:"other"`
	Total         int64 `json:"total"`
}
