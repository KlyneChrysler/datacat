package observe

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

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
