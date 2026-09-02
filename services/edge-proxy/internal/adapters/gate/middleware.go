// Package gate enforces standing decisions on the hot request path: a
// blocked session gets 403, a challenged one 429, a rate-limited one is
// throttled by the token bucket. Lookups are in-memory — the gate adds no
// I/O to the request.
package gate

import (
	"log/slog"
	"net/http"

	"github.com/KlyneChrysler/datacat/pkg/httpx"
	"github.com/KlyneChrysler/datacat/services/edge-proxy/internal/adapters/ident"
	"github.com/KlyneChrysler/datacat/services/edge-proxy/internal/app"
)

func Middleware(gatekeeper *app.Gatekeeper, limiter *app.RateLimiter, log *slog.Logger) httpx.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionID := ident.SessionID(r)
			action := gatekeeper.ActionFor(sessionID)
			if deny(w, r, action, limiter, log, sessionID) {
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// deny writes the refusal response for restricting actions and reports
// whether the request was stopped.
func deny(w http.ResponseWriter, r *http.Request, action string, limiter *app.RateLimiter,
	log *slog.Logger, sessionID string) bool {
	switch action {
	case "block":
		log.InfoContext(r.Context(), "request blocked", "session_id", sessionID, "path", r.URL.Path)
		httpx.Error(w, http.StatusForbidden, "session blocked")
		return true
	case "challenge":
		log.InfoContext(r.Context(), "request challenged", "session_id", sessionID, "path", r.URL.Path)
		httpx.Error(w, http.StatusTooManyRequests, "verification required")
		return true
	case "rate_limit":
		return denyOverLimit(w, r, limiter, log, sessionID)
	default:
		return false
	}
}

func denyOverLimit(w http.ResponseWriter, r *http.Request, limiter *app.RateLimiter,
	log *slog.Logger, sessionID string) bool {
	if limiter.Allow(sessionID) {
		return false
	}
	log.InfoContext(r.Context(), "request rate limited", "session_id", sessionID, "path", r.URL.Path)
	httpx.Error(w, http.StatusTooManyRequests, "rate limit exceeded")
	return true
}
