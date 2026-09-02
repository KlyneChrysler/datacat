package dynamo

// decisionRecord is the storage shape (file taxonomy, standards rule 2).
// expires_at drives the table's TTL so stale decisions age out server-side,
// mirroring the edge-proxy gate TTL.
type decisionRecord struct {
	SessionID string `dynamodbav:"session_id"`
	Class     string `dynamodbav:"classification"`
	Action    string `dynamodbav:"action"`
	ExpiresAt int64  `dynamodbav:"expires_at"`
}
