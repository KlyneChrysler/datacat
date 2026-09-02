// Package proxy is the reverse-proxy adapter: forwards traffic to the
// protected upstream. Event emission (request observation → Kafka) hooks in
// here in the Kafka phase.
package proxy

import (
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/KlyneChrysler/datacat/pkg/httpx"
)

func New(upstream *url.URL, log *slog.Logger) http.Handler {
	reverse := httputil.NewSingleHostReverseProxy(upstream)
	reverse.ErrorHandler = errorHandler(log)
	return reverse
}

func errorHandler(log *slog.Logger) func(http.ResponseWriter, *http.Request, error) {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		log.ErrorContext(r.Context(), "upstream unreachable", "err", err, "path", r.URL.Path)
		httpx.Error(w, http.StatusBadGateway, "upstream unavailable")
	}
}
