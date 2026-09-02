package app

import (
	"crypto/ed25519"
	"time"
)

// agentMaxSkew bounds signature replay.
const agentMaxSkew = 5 * time.Minute

// AgentVerifier checks trusted agent request signatures.
type AgentVerifier struct {
	keys map[string]ed25519.PublicKey
}

func NewAgentVerifier(keys map[string]ed25519.PublicKey) *AgentVerifier {
	return &AgentVerifier{keys: keys}
}

// Verify accepts a fresh signature of base by a registered key.
func (v *AgentVerifier) Verify(keyID, base string, sig []byte, issuedAt, now time.Time) bool {
	if skew := now.Sub(issuedAt); skew > agentMaxSkew || skew < -agentMaxSkew {
		return false
	}

	key, trusted := v.keys[keyID]
	if !trusted {
		return false
	}

	return ed25519.Verify(key, []byte(base), sig)
}
