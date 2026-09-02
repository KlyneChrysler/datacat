// Package domain holds pure simulation logic.
package domain

// Request is one synthetic request, credential is nil for unsigned personas.
type Request struct {
	SessionID  string
	Path       string
	UserAgent  string
	Credential *AgentCredential
}
