package gate

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/KlyneChrysler/datacat/pkg/edge/ident"
	"github.com/KlyneChrysler/datacat/pkg/guard"
	"github.com/KlyneChrysler/datacat/pkg/obsx"
)

func testChallenger() *guard.Challenger {
	return guard.NewChallenger("test-secret", 4)
}

func gatedHandler(t *testing.T, action string, burst int) http.Handler {
	t.Helper()

	gatekeeper := guard.NewGatekeeper(time.Hour)
	if action != "" {
		gatekeeper.Update("s-1", action)
	}
	gate := New(gatekeeper, guard.NewRateLimiter(60, burst, time.Hour), testChallenger(), obsx.NewLogger("test"))
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })

	return gate.Middleware()(next)
}

func send(handler http.Handler, extras ...*http.Cookie) int {
	return sendWith(handler, "", extras...)
}

func sendWith(handler http.Handler, accept string, extras ...*http.Cookie) int {
	r := httptest.NewRequest(http.MethodGet, "/products", nil)
	r.AddCookie(&http.Cookie{Name: ident.SessionCookie, Value: "s-1"})
	for _, cookie := range extras {
		r.AddCookie(cookie)
	}
	if accept != "" {
		r.Header.Set("Accept", accept)
	}

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	return w.Code
}

func TestGate(t *testing.T) {
	tests := []struct {
		name       string
		action     string
		wantStatus int
	}{
		{name: "blocked session gets 403", action: "block", wantStatus: http.StatusForbidden},
		{name: "challenged api client gets 429", action: "challenge", wantStatus: http.StatusTooManyRequests},
		{name: "allowed session passes", action: "allow", wantStatus: http.StatusOK},
		{name: "rate limited session passes within burst", action: "rate_limit", wantStatus: http.StatusOK},
		{name: "unknown session passes", action: "", wantStatus: http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := send(gatedHandler(t, tt.action, 3)); got != tt.wantStatus {
				t.Errorf("status = %d, want %d", got, tt.wantStatus)
			}
		})
	}
}

func TestGateThrottlesRateLimitedSessionBeyondBurst(t *testing.T) {
	handler := gatedHandler(t, "rate_limit", 2)

	first, second, third := send(handler), send(handler), send(handler)

	if first != http.StatusOK || second != http.StatusOK {
		t.Fatalf("requests within burst = %d, %d, want 200, 200", first, second)
	}
	if third != http.StatusTooManyRequests {
		t.Errorf("request beyond burst = %d, want 429", third)
	}
}

func TestGateRedirectsChallengedBrowsersToInterstitial(t *testing.T) {
	handler := gatedHandler(t, "challenge", 3)

	if got := sendWith(handler, "text/html,application/xhtml+xml"); got != http.StatusFound {
		t.Errorf("status = %d, want 302 redirect to challenge page", got)
	}
}

func TestGatePassesChallengedSessionWithValidClearance(t *testing.T) {
	gatekeeper := guard.NewGatekeeper(time.Hour)
	gatekeeper.Update("s-1", "challenge")
	challenger := testChallenger()
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	handler := New(gatekeeper, guard.NewRateLimiter(60, 3, time.Hour), challenger, obsx.NewLogger("test")).Middleware()(next)

	clearance := &http.Cookie{Name: guard.ClearanceCookie, Value: challenger.MintClearance("s-1", time.Now())}
	if got := send(handler, clearance); got != http.StatusOK {
		t.Errorf("status = %d, want 200 for cleared session", got)
	}

	forged := &http.Cookie{Name: guard.ClearanceCookie, Value: "forged"}
	if got := send(handler, forged); got != http.StatusTooManyRequests {
		t.Errorf("status = %d, want 429 for forged clearance", got)
	}
}
