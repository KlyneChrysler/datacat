package gate

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/KlyneChrysler/datacat/pkg/events"
	"github.com/KlyneChrysler/datacat/pkg/obsx"
	"github.com/KlyneChrysler/datacat/services/edge-proxy/internal/adapters/ident"
	"github.com/KlyneChrysler/datacat/services/edge-proxy/internal/app"
)

func gatedHandler(t *testing.T, action string, burst int) http.Handler {
	t.Helper()
	gatekeeper := app.NewGatekeeper(time.Hour)
	if action != "" {
		gatekeeper.Update(events.Decision{SessionID: "s-1", Action: action})
	}
	limiter := app.NewRateLimiter(60, burst, time.Hour)
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	return Middleware(gatekeeper, limiter, obsx.NewLogger("test"))(next)
}

func send(handler http.Handler) int {
	r := httptest.NewRequest(http.MethodGet, "/products", nil)
	r.AddCookie(&http.Cookie{Name: ident.SessionCookie, Value: "s-1"})
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
		{name: "challenged session gets 429", action: "challenge", wantStatus: http.StatusTooManyRequests},
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
