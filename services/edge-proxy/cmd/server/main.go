// Command server is the edge-proxy composition root: wiring and lifecycle
// only.
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/KlyneChrysler/datacat/pkg/httpx"
	"github.com/KlyneChrysler/datacat/pkg/obsx"
	"github.com/KlyneChrysler/datacat/services/edge-proxy/internal/adapters/proxy"
	"github.com/KlyneChrysler/datacat/services/edge-proxy/internal/config"
)

func main() {
	if err := run(); err != nil {
		slog.Error("fatal", "err", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	log := obsx.NewLogger("edge-proxy")

	router := newRouter(cfg, log)
	server := httpx.NewServer(cfg.Port, router, cfg.ShutdownTimeout)

	log.Info("starting", "port", cfg.Port, "upstream", cfg.UpstreamURL.String())
	return server.ListenAndServe(ctx)
}

func newRouter(cfg config.Config, log *slog.Logger) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		httpx.JSON(w, http.StatusOK, map[string]string{"status": "alive"})
	})
	mux.Handle("/", proxy.New(cfg.UpstreamURL, log))
	return httpx.WithMiddleware(mux, httpx.RequestID(), httpx.Logging(log), httpx.Recover(log))
}
