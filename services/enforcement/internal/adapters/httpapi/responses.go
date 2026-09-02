package httpapi

// Wire shapes only (file taxonomy, standards rule 2).

type decisionResponse struct {
	SessionID string `json:"session_id"`
	Class     string `json:"classification"`
	Action    string `json:"action"`
}
