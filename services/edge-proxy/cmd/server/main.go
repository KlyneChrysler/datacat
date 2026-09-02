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

	"golang.org/x/sync/errgroup"

	"github.com/KlyneChrysler/datacat/pkg/httpx"
	"github.com/KlyneChrysler/datacat/pkg/kafkax"
	"github.com/KlyneChrysler/datacat/pkg/obsx"
	"github.com/KlyneChrysler/datacat/services/edge-proxy/internal/adapters/kafka"
	"github.com/KlyneChrysler/datacat/services/edge-proxy/internal/adapters/observe"
	"github.com/KlyneChrysler/datacat/services/edge-proxy/internal/adapters/proxy"
	"github.com/KlyneChrysler/datacat/services/edge-proxy/internal/app"
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

	producer, err := kafkax.NewProducer(cfg.KafkaBrokers)
	if err != nil {
		return err
	}
	defer producer.Close()

	publisher := kafka.NewEventPublisher(producer, cfg.RequestsTopic)
	recorder := app.NewRecorder(publisher, log, cfg.EventBufferSize)
	server := httpx.NewServer(cfg.Port, newRouter(cfg, recorder, log), cfg.ShutdownTimeout)

	log.Info("starting", "port", cfg.Port, "upstream", cfg.UpstreamURL.String(), "topic", cfg.RequestsTopic)
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error { return recorder.Run(ctx) })
	g.Go(func() error { return server.ListenAndServe(ctx) })
	return g.Wait()
}

func newRouter(cfg config.Config, recorder *app.Recorder, log *slog.Logger) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		httpx.JSON(w, http.StatusOK, map[string]string{"status": "alive"})
	})
	mux.Handle("/", proxy.New(cfg.UpstreamURL, log))
	return httpx.WithMiddleware(mux,
		httpx.RequestID(), httpx.Logging(log), httpx.Recover(log), observe.Middleware(recorder))
}
