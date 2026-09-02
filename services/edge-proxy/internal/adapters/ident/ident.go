// Package ident resolves a request to a session identity. Both the observe
// (event emission) and gate (enforcement) adapters use it, so a session is
// identified identically on the way in and on the way out.
package ident

import (
	"net"
	"net/http"

	"github.com/KlyneChrysler/datacat/pkg/hashx"
)

const SessionCookie = "dc_session"

// SessionID prefers the session cookie; without one it falls back to a
// fingerprint of IP + User-Agent so anonymous traffic still groups. O(1).
func SessionID(r *http.Request) string {
	if cookie, err := r.Cookie(SessionCookie); err == nil && cookie.Value != "" {
		return cookie.Value
	}
	return hashx.Short(ClientIP(r) + "|" + r.UserAgent())
}

func ClientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
