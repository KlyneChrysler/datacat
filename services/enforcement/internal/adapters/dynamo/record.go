package dynamo

// decisionRecord is the storage shape, expires_at drives the table ttl.
type decisionRecord struct {
	SessionID string `dynamodbav:"session_id"`
	Class     string `dynamodbav:"classification"`
	Action    string `dynamodbav:"action"`
	ExpiresAt int64  `dynamodbav:"expires_at"`
}
