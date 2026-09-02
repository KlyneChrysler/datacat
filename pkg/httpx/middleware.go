package httpx

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"time"
)

type Middleware func(http.Handler) http.Handler

// WithMiddleware wraps h so the first middleware runs outermost.
func WithMiddleware(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}

	return h
}

// RequestID tags every request with a fresh id.
func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := newRequestID()

			w.Header().Set("X-Request-ID", id)
			r.Header.Set("X-Request-ID", id)

			next.ServeHTTP(w, r)
		})
	}
}

// Logging writes one structured line per request.
func Logging(log *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			next.ServeHTTP(w, r)

			log.InfoContext(r.Context(), "request", "method", r.Method, "path", r.URL.Path, "duration_ms", time.Since(start).Milliseconds(), "request_id", r.Header.Get("X-Request-ID"))
		})
	}
}

// Recover turns panics into 500 responses instead of crashes.
func Recover(log *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					log.ErrorContext(r.Context(), "panic recovered", "panic", rec, "path", r.URL.Path)
					Error(w, http.StatusInternalServerError, "internal error")
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func newRequestID() string {
	buf := make([]byte, 8)
	_, _ = rand.Read(buf) // rand.Read never fails per its contract

	return hex.EncodeToString(buf)
}
