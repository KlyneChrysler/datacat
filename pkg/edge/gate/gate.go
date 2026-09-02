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
	"github.com/KlyneChrysler/datacat/pkg/policy"
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
	case string(policy.Block):
		g.refuse(w, r, sessionID, refusal{event: "request blocked", status: http.StatusForbidden, message: "session blocked"})
		return true
	case string(policy.Challenge):
		return g.denyUncleared(w, r, sessionID)
	case string(policy.RateLimit):
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

	g.refuse(w, r, sessionID, refusal{event: "request challenged", status: http.StatusTooManyRequests, message: "verification required"})
	return true
}

// denyOverLimit throttles rate limited sessions past their bucket.
func (g *Gate) denyOverLimit(w http.ResponseWriter, r *http.Request, sessionID string) bool {
	if g.limiter.Allow(sessionID) {
		return false
	}

	g.refuse(w, r, sessionID, refusal{event: "request rate limited", status: http.StatusTooManyRequests, message: "rate limit exceeded"})
	return true
}

func (g *Gate) refuse(w http.ResponseWriter, r *http.Request, sessionID string, ref refusal) {
	g.log.InfoContext(r.Context(), ref.event, "session_id", sessionID, "path", r.URL.Path)

	httpx.Error(w, ref.status, ref.message)
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
