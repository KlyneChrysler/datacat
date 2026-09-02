// Package ident resolves a request to a session identity.
package ident

import (
	"net"
	"net/http"

	"github.com/KlyneChrysler/datacat/pkg/events"
	"github.com/KlyneChrysler/datacat/pkg/hashx"
)

// SessionCookie re-exports the wire contract cookie name.
const SessionCookie = events.SessionCookie

// SessionID uses the session cookie, else a fingerprint of ip and agent.
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
