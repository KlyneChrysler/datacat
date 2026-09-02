package app

import (
	"crypto/ed25519"
	"time"
)

// agentMaxSkew bounds replay: a signature older or newer than this is
// rejected regardless of validity.
const agentMaxSkew = 5 * time.Minute

// AgentVerifier checks trusted-agent request signatures (a simplified
// Web-Bot-Auth-style profile: Ed25519 over a canonical request base, keys
// pre-registered via config). Verify is O(1): one map lookup plus one
// signature check.
type AgentVerifier struct {
	keys map[string]ed25519.PublicKey
}

func NewAgentVerifier(keys map[string]ed25519.PublicKey) *AgentVerifier {
	return &AgentVerifier{keys: keys}
}

// Verify reports whether sig is a valid signature of base by the registered
// key, with issuedAt inside the replay window.
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
