// Package httpx holds the shared HTTP toolkit.
package httpx

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// Server is an http.Server that drains cleanly on context cancel.
type Server struct {
	inner           *http.Server
	shutdownTimeout time.Duration
}

func NewServer(port string, handler http.Handler, shutdownTimeout time.Duration) *Server {
	inner := &http.Server{Addr: ":" + port, Handler: handler, ReadHeaderTimeout: 5 * time.Second}

	return &Server{inner: inner, shutdownTimeout: shutdownTimeout}
}

// ListenAndServe serves until ctx is cancelled, then drains and returns.
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
