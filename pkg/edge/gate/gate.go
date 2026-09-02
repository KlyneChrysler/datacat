// Package gate enforces standing decisions on the request path.
package gate

import (
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/KlyneChrysler/datacat/pkg/edge/challenge"
	"github.com/KlyneChrysler/datacat/pkg/edge/ident"
	"github.com/KlyneChrysler/datacat/pkg/guard"
	"github.com/KlyneChrysler/datacat/pkg/httpx"
)

// Gate holds the collaborators every gating step needs.
type Gate struct {
	gatekeeper *guard.Gatekeeper
	limiter    *guard.RateLimiter
	challenger *guard.Challenger
	log        *slog.Logger
}

func New(gatekeeper *guard.Gatekeeper, limiter *guard.RateLimiter, challenger *guard.Challenger, log *slog.Logger) *Gate {
	return &Gate{gatekeeper: gatekeeper, limiter: limiter, challenger: challenger, log: log}
}

// Middleware stops requests whose standing action forbids them.
func (g *Gate) Middleware() httpx.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionID := ident.SessionID(r)
			action := g.gatekeeper.ActionFor(sessionID)

			if g.deny(w, r, action, sessionID) {
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// deny refuses the request when its action requires that.
func (g *Gate) deny(w http.ResponseWriter, r *http.Request, action, sessionID string) bool {
	switch action {
	case "block":
		g.refuse(w, r, sessionID, "request blocked", http.StatusForbidden, "session blocked")
		return true
	case "challenge":
		return g.denyUncleared(w, r, sessionID)
	case "rate_limit":
		return g.denyOverLimit(w, r, sessionID)
	default:
		return false
	}
}

// denyUncleared sends browsers to the interstitial and refuses api clients.
func (g *Gate) denyUncleared(w http.ResponseWriter, r *http.Request, sessionID string) bool {
	if g.hasClearance(r, sessionID) {
		return false
	}

	if wantsHTML(r) {
		g.log.InfoContext(r.Context(), "request challenged", "session_id", sessionID, "path", r.URL.Path)
		http.Redirect(w, r, challengeURL(r), http.StatusFound)
		return true
	}

	g.refuse(w, r, sessionID, "request challenged", http.StatusTooManyRequests, "verification required")
	return true
}

// denyOverLimit throttles rate limited sessions past their bucket.
func (g *Gate) denyOverLimit(w http.ResponseWriter, r *http.Request, sessionID string) bool {
	if g.limiter.Allow(sessionID) {
		return false
	}

	g.refuse(w, r, sessionID, "request rate limited", http.StatusTooManyRequests, "rate limit exceeded")
	return true
}

func (g *Gate) refuse(w http.ResponseWriter, r *http.Request, sessionID, event string, status int, message string) {
	g.log.InfoContext(r.Context(), event, "session_id", sessionID, "path", r.URL.Path)

	httpx.Error(w, status, message)
}

func (g *Gate) hasClearance(r *http.Request, sessionID string) bool {
	cookie, err := r.Cookie(guard.ClearanceCookie)
	if err != nil {
		return false
	}

	return g.challenger.ValidClearance(cookie.Value, sessionID, time.Now())
}

func wantsHTML(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Accept"), "text/html")
}

func challengeURL(r *http.Request) string {
	return challenge.PagePath + "?return=" + url.QueryEscape(r.URL.RequestURI())
}
