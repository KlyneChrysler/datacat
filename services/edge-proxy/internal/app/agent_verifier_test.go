package app

import (
	"crypto/ed25519"
	"crypto/sha256"
	"testing"
	"time"
)

func testKeyPair() (ed25519.PrivateKey, ed25519.PublicKey) {
	seed := sha256.Sum256([]byte("test-seed"))
	priv := ed25519.NewKeyFromSeed(seed[:])
	return priv, priv.Public().(ed25519.PublicKey)
}

func TestAgentVerifier(t *testing.T) {
	priv, pub := testKeyPair()
	verifier := NewAgentVerifier(map[string]ed25519.PublicKey{"agent-1": pub})
	now := time.Unix(9_000_000, 0)
	base := "GET\n/api/catalog\nsim-agent\n9000000"
	sig := ed25519.Sign(priv, []byte(base))

	if !verifier.Verify("agent-1", base, sig, now, now) {
		t.Error("valid signature rejected")
	}
	if verifier.Verify("unknown", base, sig, now, now) {
		t.Error("unregistered key accepted")
	}
	if verifier.Verify("agent-1", base+"tampered", sig, now, now) {
		t.Error("tampered base accepted")
	}
	if verifier.Verify("agent-1", base, sig, now, now.Add(10*time.Minute)) {
		t.Error("stale signature accepted (replay window)")
	}
	if verifier.Verify("agent-1", base, sig, now.Add(10*time.Minute), now) {
		t.Error("future-dated signature accepted")
	}
}
