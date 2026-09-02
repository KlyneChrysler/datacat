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
	"github.com/KlyneChrysler/datacat/services/edge-proxy/internal/adapters/gate"
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
	// Every proxy replica must see EVERY decision (broadcast), so each
	// instance gets its own consumer group - a shared group would split
	// decisions across replicas. A fresh instance replays the retained
	// topic, rebuilding its gate state on startup.
	consumer, err := kafkax.NewConsumer(cfg.KafkaBrokers, instanceGroup(cfg.DecisionsGroup), cfg.DecisionsTopic, log)
	if err != nil {
		return err
	}

	publisher := kafka.NewEventPublisher(producer, cfg.RequestsTopic)
	recorder := app.NewRecorder(publisher, log, cfg.EventBufferSize)
	gatekeeper := app.NewGatekeeper(cfg.GateTTL)
	decisions := kafka.NewDecisionSource(consumer)
	server := httpx.NewServer(cfg.Port, newRouter(cfg, recorder, gatekeeper, log), cfg.ShutdownTimeout)

	log.Info("starting", "port", cfg.Port, "upstream", cfg.UpstreamURL.String(), "topic", cfg.RequestsTopic)
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error { return recorder.Run(ctx) })
	g.Go(func() error { return decisions.Consume(ctx, gatekeeper.Update) })
	g.Go(func() error { return server.ListenAndServe(ctx) })
	return g.Wait()
}

func instanceGroup(base string) string {
	hostname, err := os.Hostname()
	if err != nil {
		return base
	}
	return base + "-" + hostname
}

// newRouter wires the traffic path: observe first (blocked requests are
// signal too), then the gate, then the upstream proxy. Health endpoints stay
// outside both.
func newRouter(cfg config.Config, recorder *app.Recorder, gatekeeper *app.Gatekeeper,
	log *slog.Logger) http.Handler {
	traffic := httpx.WithMiddleware(proxy.New(cfg.UpstreamURL, log),
		observe.Middleware(recorder), gate.Middleware(gatekeeper, log))

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		httpx.JSON(w, http.StatusOK, map[string]string{"status": "alive"})
	})
	mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, _ *http.Request) {
		httpx.JSON(w, http.StatusOK, map[string]string{"status": "ready"})
	})
	mux.Handle("/", traffic)
	return httpx.WithMiddleware(mux, httpx.RequestID(), httpx.Logging(log), httpx.Recover(log))
}
