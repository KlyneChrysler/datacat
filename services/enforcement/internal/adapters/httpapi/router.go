// Package httpapi holds the inbound HTTP adapter.
package httpapi

import (
	"log/slog"
	"net/http"

	"github.com/KlyneChrysler/datacat/pkg/httpx"
	"github.com/KlyneChrysler/datacat/pkg/obsx"
)

// NewRouter wires routes and middleware, empty corsOrigin disables cors.
func NewRouter(h *Handlers, log *slog.Logger, corsOrigin string, metrics *obsx.Metrics) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", httpx.Alive)
	mux.HandleFunc("GET /readyz", httpx.Ready)
	mux.Handle("GET /metrics", metrics.Handler())
	mux.HandleFunc("GET /v1/decisions/{sessionID}", h.GetDecision)
	mux.HandleFunc("GET /v1/traffic/summary", h.GetTrafficSummary)

	return httpx.WithMiddleware(mux, middlewares(log, corsOrigin)...)
}

func middlewares(log *slog.Logger, corsOrigin string) []httpx.Middleware {
	chain := []httpx.Middleware{httpx.RequestID(), httpx.Logging(log), httpx.Recover(log)}
	if corsOrigin != "" {
		chain = append(chain, httpx.CORS(corsOrigin))
	}

	return chain
}
