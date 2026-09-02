package domain

import "crypto/ed25519"

// AgentCredential lets a persona sign requests as a trusted agent.
type AgentCredential struct {
	KeyID string
	Key   ed25519.PrivateKey
}
