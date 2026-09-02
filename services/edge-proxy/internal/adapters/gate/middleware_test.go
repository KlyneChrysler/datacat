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

func gatedRequest(t *testing.T, action string) *httptest.ResponseRecorder {
	t.Helper()
	gatekeeper := app.NewGatekeeper(time.Hour)
	if action != "" {
		gatekeeper.Update(events.Decision{SessionID: "s-1", Action: action})
	}
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := Middleware(gatekeeper, obsx.NewLogger("test"))(next)

	r := httptest.NewRequest(http.MethodGet, "/products", nil)
	r.AddCookie(&http.Cookie{Name: ident.SessionCookie, Value: "s-1"})
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	return w
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
		{name: "rate limited session passes until limiter phase", action: "rate_limit", wantStatus: http.StatusOK},
		{name: "unknown session passes", action: "", wantStatus: http.StatusOK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := gatedRequest(t, tt.action).Code; got != tt.wantStatus {
				t.Errorf("status = %d, want %d", got, tt.wantStatus)
			}
		})
	}
}
