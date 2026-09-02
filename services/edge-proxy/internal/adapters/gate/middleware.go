// Package gate enforces standing decisions on the hot request path: a
// blocked session gets 403, a challenged one 429, everything else passes.
// Lookup is an in-memory read — the gate adds no I/O to the request.
package gate

import (
	"log/slog"
	"net/http"

	"github.com/KlyneChrysler/datacat/pkg/httpx"
	"github.com/KlyneChrysler/datacat/services/edge-proxy/internal/adapters/ident"
	"github.com/KlyneChrysler/datacat/services/edge-proxy/internal/app"
)

func Middleware(gatekeeper *app.Gatekeeper, log *slog.Logger) httpx.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionID := ident.SessionID(r)
			action := gatekeeper.ActionFor(sessionID)
			if deny(w, r, action, log, sessionID) {
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// deny writes the refusal response for blocking actions and reports whether
// the request was stopped. rate_limit passes through until the limiter phase.
func deny(w http.ResponseWriter, r *http.Request, action string, log *slog.Logger, sessionID string) bool {
	switch action {
	case "block":
		log.InfoContext(r.Context(), "request blocked", "session_id", sessionID, "path", r.URL.Path)
		httpx.Error(w, http.StatusForbidden, "session blocked")
		return true
	case "challenge":
		log.InfoContext(r.Context(), "request challenged", "session_id", sessionID, "path", r.URL.Path)
		httpx.Error(w, http.StatusTooManyRequests, "verification required")
		return true
	default:
		return false
	}
}
