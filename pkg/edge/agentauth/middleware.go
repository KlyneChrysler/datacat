package agentauth

import (
	"encoding/base64"
	"net/http"
	"strconv"
	"time"

	"github.com/KlyneChrysler/datacat/pkg/agentsig"
	"github.com/KlyneChrysler/datacat/pkg/edge/ident"
	"github.com/KlyneChrysler/datacat/pkg/guard"
	"github.com/KlyneChrysler/datacat/pkg/httpx"
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
	keyID := r.Header.Get(agentsig.HeaderKey)
	if keyID == "" {
		return false
	}

	sig, err := base64.StdEncoding.DecodeString(r.Header.Get(agentsig.HeaderSignature))
	if err != nil {
		return false
	}
	issued, err := strconv.ParseInt(r.Header.Get(agentsig.HeaderTimestamp), 10, 64)
	if err != nil {
		return false
	}

	base := agentsig.Base(r.Method, r.URL.Path, ident.SessionID(r), issued)
	return verifier.Verify(keyID, base, sig, time.Unix(issued, 0), time.Now())
}
