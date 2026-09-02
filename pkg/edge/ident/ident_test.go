package ident

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSessionIDPrefersCookie(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/products", nil)
	r.AddCookie(&http.Cookie{Name: SessionCookie, Value: "cookie-session"})

	if got := SessionID(r); got != "cookie-session" {
		t.Fatalf("SessionID = %q, want cookie value", got)
	}
}

func TestSessionIDFallbackIsStablePerClient(t *testing.T) {
	first := httptest.NewRequest(http.MethodGet, "/a", nil)
	first.Header.Set("User-Agent", "scraper/1.0")
	second := httptest.NewRequest(http.MethodGet, "/b", nil)
	second.Header.Set("User-Agent", "scraper/1.0")
	other := httptest.NewRequest(http.MethodGet, "/a", nil)
	other.Header.Set("User-Agent", "different/2.0")

	if SessionID(first) != SessionID(second) {
		t.Error("same client produced different session ids")
	}
	if SessionID(first) == SessionID(other) {
		t.Error("different clients produced the same session id")
	}
}
