// Package observe turns passing requests into RequestEvents. It is an
// adapter on the inbound side: it reads the request and hands a fact to the
// app layer, adding zero latency and never failing the request.
package observe

import (
	"crypto/sha256"
	"encoding/hex"
	"net"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/KlyneChrysler/datacat/pkg/events"
	"github.com/KlyneChrysler/datacat/pkg/httpx"
	"github.com/KlyneChrysler/datacat/services/edge-proxy/internal/app"
)

const sessionCookie = "dc_session"

func Middleware(recorder *app.Recorder) httpx.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			recorder.Record(eventFrom(r))
			next.ServeHTTP(w, r)
		})
	}
}

func eventFrom(r *http.Request) events.RequestEvent {
	return events.RequestEvent{
		SessionID:      sessionID(r),
		Timestamp:      time.Now().UTC(),
		Method:         r.Method,
		Path:           r.URL.Path,
		ClientIP:       clientIP(r),
		UserAgent:      r.UserAgent(),
		HeaderOrder:    headerOrderHash(r),
		TLSFingerprint: "", // requires raw ClientHello capture; later phase
	}
}

// sessionID prefers the session cookie; without one it falls back to a
// fingerprint of IP + User-Agent so anonymous traffic still groups.
func sessionID(r *http.Request) string {
	if cookie, err := r.Cookie(sessionCookie); err == nil && cookie.Value != "" {
		return cookie.Value
	}
	return shortHash(clientIP(r) + "|" + r.UserAgent())
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// headerOrderHash approximates a header fingerprint. net/http stores headers
// in a map, losing wire order, so this hashes the sorted name set; true
// wire-order capture needs a lower-level listener (later phase).
func headerOrderHash(r *http.Request) string {
	names := make([]string, 0, len(r.Header))
	for name := range r.Header {
		names = append(names, name)
	}
	slices.Sort(names)
	return shortHash(strings.Join(names, ","))
}

func shortHash(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:8])
}
