// Package httpx provides the shared HTTP toolkit for datacat services:
// lifecycle-managed server, middleware, and the canonical response envelope.
package httpx

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// Server wraps http.Server with context-driven graceful shutdown
// (twelve-factor IX: disposability).
type Server struct {
	inner           *http.Server
	shutdownTimeout time.Duration
}

func NewServer(port string, handler http.Handler, shutdownTimeout time.Duration) *Server {
	return &Server{
		inner: &http.Server{
			Addr:              ":" + port,
			Handler:           handler,
			ReadHeaderTimeout: 5 * time.Second,
		},
		shutdownTimeout: shutdownTimeout,
	}
}

// ListenAndServe blocks until ctx is cancelled, then drains in-flight
// requests within the shutdown timeout before returning.
func (s *Server) ListenAndServe(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() { errCh <- s.inner.ListenAndServe() }()

	select {
	case err := <-errCh:
		return fmt.Errorf("http server: %w", err)
	case <-ctx.Done():
		return s.shutdown()
	}
}

func (s *Server) shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()
	if err := s.inner.Shutdown(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("http shutdown: %w", err)
	}
	return nil
}
