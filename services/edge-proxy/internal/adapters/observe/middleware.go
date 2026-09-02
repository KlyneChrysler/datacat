// Package observe turns passing requests into RequestEvents. It is an
// adapter on the inbound side: it reads the request and hands a fact to the
// app layer, adding zero latency and never failing the request.
// Request → event conversion lives in event_mapper.go (file taxonomy).
package observe

import (
	"net/http"

	"github.com/KlyneChrysler/datacat/pkg/httpx"
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
