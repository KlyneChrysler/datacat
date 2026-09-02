// Package watch feeds passing requests into the analyzer.
package watch

import (
	"net/http"
	"time"

	"github.com/KlyneChrysler/datacat/pkg/edge/agentauth"
	"github.com/KlyneChrysler/datacat/pkg/edge/ident"
	"github.com/KlyneChrysler/datacat/pkg/httpx"
	"github.com/KlyneChrysler/datacat/services/lite/internal/app"
	"github.com/KlyneChrysler/datacat/services/lite/internal/classify"
)

// Middleware observes every request inline, bounded cost per request.
func Middleware(analyzer *app.Analyzer) httpx.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			analyzer.Observe(observationFrom(r))

			next.ServeHTTP(w, r)
		})
	}
}

// observationFrom converts one request into its observation.
func observationFrom(r *http.Request) classify.Observation {
	return classify.Observation{
		SessionID: ident.SessionID(r),
		Path:      r.URL.Path,
		UserAgent: r.UserAgent(),
		Verified:  agentauth.Verified(r.Context()),
		At:        time.Now(),
	}
}
