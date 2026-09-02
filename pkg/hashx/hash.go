// Package hashx provides the shared short-hash used for fingerprints and
// fallback identities across services.
package hashx

import (
	"crypto/sha256"
	"encoding/hex"
)

// Short returns a stable 16-hex-char digest of s. O(len(s)).
func Short(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:8])
}
