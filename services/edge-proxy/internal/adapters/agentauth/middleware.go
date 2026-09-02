package agentauth

import (
	"encoding/base64"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/KlyneChrysler/datacat/pkg/httpx"
	"github.com/KlyneChrysler/datacat/services/edge-proxy/internal/adapters/ident"
	"github.com/KlyneChrysler/datacat/services/edge-proxy/internal/app"
)

// Signature headers (simplified Web-Bot-Auth-style profile; the signature
// base is method, path, session, and timestamp, newline-joined).
const (
	HeaderKey       = "X-Agent-Key"
	HeaderSignature = "X-Agent-Signature"
	HeaderTimestamp = "X-Agent-Timestamp"
)

// Middleware marks requests bearing a valid trusted-agent signature. An
// invalid or absent signature is not an error — the request simply stays
// unverified and the classifier judges it on behavior.
func Middleware(verifier *app.AgentVerifier) httpx.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if verifiedRequest(r, verifier) {
				r = r.WithContext(withVerified(r.Context()))
			}
			next.ServeHTTP(w, r)
		})
	}
}

func verifiedRequest(r *http.Request, verifier *app.AgentVerifier) bool {
	keyID := r.Header.Get(HeaderKey)
	if keyID == "" {
		return false
	}
	sig, err := base64.StdEncoding.DecodeString(r.Header.Get(HeaderSignature))
	if err != nil {
		return false
	}
	issued, err := strconv.ParseInt(r.Header.Get(HeaderTimestamp), 10, 64)
	if err != nil {
		return false
	}
	base := SignatureBase(r.Method, r.URL.Path, ident.SessionID(r), issued)
	return verifier.Verify(keyID, base, sig, time.Unix(issued, 0), time.Now())
}

// SignatureBase is the canonical string both signer and verifier commit to.
func SignatureBase(method, path, sessionID string, issuedUnix int64) string {
	return strings.Join([]string{method, path, sessionID, strconv.FormatInt(issuedUnix, 10)}, "\n")
}
