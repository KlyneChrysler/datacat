// Package observe turns passing requests into RequestEvents. It is an
// adapter on the inbound side: it reads the request and hands a fact to the
// app layer, adding zero latency and never failing the request.
package observe

import (
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/KlyneChrysler/datacat/pkg/events"
	"github.com/KlyneChrysler/datacat/pkg/httpx"
	"github.com/KlyneChrysler/datacat/services/edge-proxy/internal/adapters/ident"
	"github.com/KlyneChrysler/datacat/services/edge-proxy/internal/app"
)

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
		SessionID:      ident.SessionID(r),
		Timestamp:      time.Now().UTC(),
		Method:         r.Method,
		Path:           r.URL.Path,
		ClientIP:       ident.ClientIP(r),
		UserAgent:      r.UserAgent(),
		HeaderOrder:    headerOrderHash(r),
		TLSFingerprint: "", // requires raw ClientHello capture; later phase
	}
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
	return ident.ShortHash(strings.Join(names, ","))
}
