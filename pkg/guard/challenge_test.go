package guard

import (
	"strconv"
	"testing"
	"time"
)

const testDifficulty = 4 // tiny so tests solve in microseconds

func solvedNonce(t *testing.T, c *Challenger, token, sessionID string, now time.Time) string {
	t.Helper()
	for n := 0; n < 1_000_000; n++ {
		nonce := strconv.Itoa(n)
		if c.VerifySolution(token, nonce, sessionID, now) {
			return nonce
		}
	}
	t.Fatal("no nonce found at test difficulty")
	return ""
}

func TestSolutionRoundTrip(t *testing.T) {
	c := NewChallenger("secret", testDifficulty)
	now := time.Unix(8_000_000, 0)
	token := c.MintChallenge("s-1", now)

	nonce := solvedNonce(t, c, token, "s-1", now)

	if !c.VerifySolution(token, nonce, "s-1", now) {
		t.Error("valid solution rejected")
	}
	if c.VerifySolution(token, nonce, "someone-else", now) {
		t.Error("solution accepted for a different session")
	}
	if c.VerifySolution(token, nonce, "s-1", now.Add(10*time.Minute)) {
		t.Error("solution accepted after token expiry")
	}
	if c.VerifySolution(token, "999999999", "s-1", now) {
		t.Error("arbitrary nonce accepted") // astronomically unlikely to solve
	}
}

func TestClearanceRoundTrip(t *testing.T) {
	c := NewChallenger("secret", testDifficulty)
	now := time.Unix(8_000_000, 0)

	clearance := c.MintClearance("s-1", now)

	if !c.ValidClearance(clearance, "s-1", now) {
		t.Error("valid clearance rejected")
	}
	if c.ValidClearance(clearance, "someone-else", now) {
		t.Error("clearance accepted for a different session")
	}
	if c.ValidClearance(clearance, "s-1", now.Add(2*time.Hour)) {
		t.Error("clearance accepted after expiry")
	}
}

func TestTamperedAndCrossLabelTokensRejected(t *testing.T) {
	c := NewChallenger("secret", testDifficulty)
	other := NewChallenger("other-secret", testDifficulty)
	now := time.Unix(8_000_000, 0)

	if c.ValidClearance(c.MintChallenge("s-1", now), "s-1", now) {
		t.Error("challenge token accepted as clearance (label confusion)")
	}
	if c.ValidClearance(other.MintClearance("s-1", now), "s-1", now) {
		t.Error("clearance signed with a different secret accepted")
	}
	if c.ValidClearance("garbage", "s-1", now) || c.ValidClearance("a.b", "s-1", now) {
		t.Error("malformed clearance accepted")
	}
}
