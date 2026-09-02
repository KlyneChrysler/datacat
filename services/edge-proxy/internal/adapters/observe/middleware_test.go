package observe

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSessionIDPrefersCookie(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/products", nil)
	r.AddCookie(&http.Cookie{Name: sessionCookie, Value: "cookie-session"})

	if got := sessionID(r); got != "cookie-session" {
		t.Fatalf("sessionID = %q, want cookie value", got)
	}
}

func TestSessionIDFallbackIsStablePerClient(t *testing.T) {
	first := httptest.NewRequest(http.MethodGet, "/a", nil)
	first.Header.Set("User-Agent", "scraper/1.0")
	second := httptest.NewRequest(http.MethodGet, "/b", nil)
	second.Header.Set("User-Agent", "scraper/1.0")
	other := httptest.NewRequest(http.MethodGet, "/a", nil)
	other.Header.Set("User-Agent", "different/2.0")

	if sessionID(first) != sessionID(second) {
		t.Error("same client produced different session ids")
	}
	if sessionID(first) == sessionID(other) {
		t.Error("different clients produced the same session id")
	}
}

func TestEventFromCapturesRequestFacts(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/checkout", nil)
	r.Header.Set("User-Agent", "agent/1.0")

	ev := eventFrom(r)

	if ev.Method != http.MethodPost || ev.Path != "/checkout" {
		t.Errorf("method/path = %s %s, want POST /checkout", ev.Method, ev.Path)
	}
	if ev.UserAgent != "agent/1.0" {
		t.Errorf("user agent = %q", ev.UserAgent)
	}
	if ev.SessionID == "" || ev.HeaderOrder == "" || ev.Timestamp.IsZero() {
		t.Error("session id, header order, and timestamp must be populated")
	}
}
