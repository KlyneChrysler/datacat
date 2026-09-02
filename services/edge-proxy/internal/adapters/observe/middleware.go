// Package observe turns passing requests into request events.
package observe

import (
	"net/http"

	"github.com/KlyneChrysler/datacat/pkg/httpx"
	"github.com/KlyneChrysler/datacat/services/edge-proxy/internal/app"
)

// Middleware records every request without adding latency.
func Middleware(recorder *app.Recorder) httpx.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			recorder.Record(eventFrom(r))

			next.ServeHTTP(w, r)
		})
	}
}
