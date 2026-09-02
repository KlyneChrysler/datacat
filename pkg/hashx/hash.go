// Package hashx holds the shared short hash helper.
package hashx

import (
	"crypto/sha256"
	"encoding/hex"
)

// Short returns a stable 16 character digest of s.
func Short(s string) string {
	sum := sha256.Sum256([]byte(s))

	return hex.EncodeToString(sum[:8])
}
