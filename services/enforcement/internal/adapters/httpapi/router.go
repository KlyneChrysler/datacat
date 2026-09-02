// Package httpapi is the inbound HTTP adapter: thin handlers that parse,
// delegate, and respond. Business logic never lives here.
package httpapi

import (
	"log/slog"
	"net/http"

	"github.com/KlyneChrysler/datacat/pkg/httpx"
)

func NewRouter(h *Handlers, log *slog.Logger) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", h.Health)
	mux.HandleFunc("GET /readyz", h.Ready)
	return httpx.WithMiddleware(mux, httpx.RequestID(), httpx.Logging(log), httpx.Recover(log))
}
