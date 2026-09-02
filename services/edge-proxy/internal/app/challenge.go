package app

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"strconv"
	"strings"
	"time"
)

// ClearanceCookie carries a solved challenge. Proxy-internal: no other
// service reads it.
const ClearanceCookie = "dc_clearance"

const (
	labelChallenge = "challenge"
	labelClearance = "clearance"

	challengeTTL = 5 * time.Minute
	clearanceTTL = time.Hour
)

// Challenger mints and verifies challenge tokens and clearances. Tokens are
// HMAC-signed session+expiry pairs; a solution is a nonce whose
// SHA-256(token:nonce) carries the required leading zero bits (proof of
// work). All operations are O(1) per call.
type Challenger struct {
	secret     []byte
	difficulty int
}

func NewChallenger(secret string, difficultyBits int) *Challenger {
	return &Challenger{secret: []byte(secret), difficulty: difficultyBits}
}

func (c *Challenger) Difficulty() int {
	return c.difficulty
}

func (c *Challenger) MintChallenge(sessionID string, now time.Time) string {
	return c.sign(labelChallenge, sessionID, now.Add(challengeTTL))
}

// VerifySolution accepts a nonce only for a valid, unexpired token bound to
// this session, whose digest clears the difficulty.
func (c *Challenger) VerifySolution(token, nonce, sessionID string, now time.Time) bool {
	if !c.valid(labelChallenge, token, sessionID, now) {
		return false
	}
	digest := sha256.Sum256([]byte(token + ":" + nonce))
	return leadingZeroBits(digest[:]) >= c.difficulty
}

func (c *Challenger) MintClearance(sessionID string, now time.Time) string {
	return c.sign(labelClearance, sessionID, now.Add(clearanceTTL))
}

func (c *Challenger) ValidClearance(value, sessionID string, now time.Time) bool {
	return c.valid(labelClearance, value, sessionID, now)
}

func (c *Challenger) sign(label, sessionID string, expiry time.Time) string {
	payload := sessionID + "|" + strconv.FormatInt(expiry.Unix(), 10)
	return base64.RawURLEncoding.EncodeToString([]byte(payload)) +
		"." + base64.RawURLEncoding.EncodeToString(c.mac(label, payload))
}

func (c *Challenger) valid(label, token, sessionID string, now time.Time) bool {
	payload, sig, ok := decodeToken(token)
	if !ok || subtle.ConstantTimeCompare(sig, c.mac(label, payload)) != 1 {
		return false
	}
	boundSession, expiry, ok := splitPayload(payload)
	return ok && boundSession == sessionID && now.Unix() < expiry
}

func (c *Challenger) mac(label, payload string) []byte {
	m := hmac.New(sha256.New, c.secret)
	m.Write([]byte(label + "|" + payload))
	return m.Sum(nil)
}

func decodeToken(token string) (payload string, sig []byte, ok bool) {
	encPayload, encSig, found := strings.Cut(token, ".")
	if !found {
		return "", nil, false
	}
	rawPayload, err := base64.RawURLEncoding.DecodeString(encPayload)
	if err != nil {
		return "", nil, false
	}
	rawSig, err := base64.RawURLEncoding.DecodeString(encSig)
	if err != nil {
		return "", nil, false
	}
	return string(rawPayload), rawSig, true
}

func splitPayload(payload string) (sessionID string, expiry int64, ok bool) {
	sessionID, rawExpiry, found := strings.Cut(payload, "|")
	if !found {
		return "", 0, false
	}
	expiry, err := strconv.ParseInt(rawExpiry, 10, 64)
	return sessionID, expiry, err == nil
}

func leadingZeroBits(digest []byte) int {
	bits := 0
	for _, b := range digest {
		if b == 0 {
			bits += 8
			continue
		}
		for mask := byte(0x80); mask > 0 && b&mask == 0; mask >>= 1 {
			bits++
		}
		break
	}
	return bits
}
