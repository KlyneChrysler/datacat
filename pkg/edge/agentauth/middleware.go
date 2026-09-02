package agentauth

import (
	"encoding/base64"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/KlyneChrysler/datacat/pkg/edge/ident"
	"github.com/KlyneChrysler/datacat/pkg/guard"
	"github.com/KlyneChrysler/datacat/pkg/httpx"
)

// Signature headers of the trusted agent profile.
const (
	HeaderKey       = "X-Agent-Key"
	HeaderSignature = "X-Agent-Signature"
	HeaderTimestamp = "X-Agent-Timestamp"
)

// Middleware marks signed requests, unsigned ones just stay unverified.
func Middleware(verifier *guard.AgentVerifier) httpx.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if verifiedRequest(r, verifier) {
				r = r.WithContext(withVerified(r.Context()))
			}

			next.ServeHTTP(w, r)
		})
	}
}

func verifiedRequest(r *http.Request, verifier *guard.AgentVerifier) bool {
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
