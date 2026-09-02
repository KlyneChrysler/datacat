package domain

import "crypto/ed25519"

// AgentCredential lets a persona sign its requests as a declared, trusted
// agent (the proxy holds the matching public key).
type AgentCredential struct {
	KeyID string
	Key   ed25519.PrivateKey
}
