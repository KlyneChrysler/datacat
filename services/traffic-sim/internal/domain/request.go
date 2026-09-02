// Package domain holds pure simulation logic. It imports zero
// infrastructure. One concept per file (file taxonomy, standards rule 2).
package domain

// Request is one synthetic request to send. Credential is nil for
// unsigned personas.
type Request struct {
	SessionID  string
	Path       string
	UserAgent  string
	Credential *AgentCredential
}
