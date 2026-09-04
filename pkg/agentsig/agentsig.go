// Package agentsig is the trusted-agent signing wire contract shared by the
// signer (traffic-sim) and the verifier (edge-proxy).
package agentsig

import (
	"strconv"
	"strings"
)

// Signature headers carried on a signed request.
const (
	HeaderKey       = "X-Agent-Key"
	HeaderSignature = "X-Agent-Signature"
	HeaderTimestamp = "X-Agent-Timestamp"
)

// Base is the canonical string both signer and verifier commit to.
func Base(method, path, sessionID string, issuedUnix int64) string {
	return strings.Join([]string{method, path, sessionID, strconv.FormatInt(issuedUnix, 10)}, "\n")
}
