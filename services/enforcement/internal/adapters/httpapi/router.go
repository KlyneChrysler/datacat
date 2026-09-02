// Package httpapi is the inbound HTTP adapter: thin handlers that parse,
// delegate, and respond. Business logic never lives here. Wire shapes in
// responses.go, conversion in mapper.go (file taxonomy).
package httpapi

import (
	"log/slog"
	"net/http"

	"github.com/KlyneChrysler/datacat/pkg/httpx"
)

// NewRouter wires routes and middleware. corsOrigin is optional: empty
// disables CORS (no browser clients in that deployment).
func NewRouter(h *Handlers, log *slog.Logger, corsOrigin string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", h.Health)
	mux.HandleFunc("GET /readyz", h.Ready)
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
