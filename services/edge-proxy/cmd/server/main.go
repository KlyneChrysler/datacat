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
	"github.com/KlyneChrysler/datacat/services/edge-proxy/internal/adapters/agentauth"
	"github.com/KlyneChrysler/datacat/services/edge-proxy/internal/adapters/challenge"
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
	consumer, err := kafkax.NewConsumer(cfg.KafkaBrokers, cfg.InstanceDecisionsGroup(), cfg.DecisionsTopic, log)
	if err != nil {
		return err
	}

	publisher := kafka.NewEventPublisher(producer, cfg.RequestsTopic)
	recorder := app.NewRecorder(publisher, log, cfg.EventBufferSize)
	gatekeeper := app.NewGatekeeper(cfg.GateTTL)
	limiter := app.NewRateLimiter(cfg.RateLimitPerMinute, cfg.RateLimitBurst, cfg.GateTTL)
	challenger := app.NewChallenger(cfg.ChallengeSecret, cfg.ChallengeDifficulty)
	decisions := kafka.NewDecisionSource(consumer)
	router := newRouter(cfg, recorder, gatekeeper, limiter, challenger, log)
	server := httpx.NewServer(cfg.Port, router, cfg.ShutdownTimeout)

	log.Info("starting", "port", cfg.Port, "upstream", cfg.UpstreamURL.String(), "topic", cfg.RequestsTopic)
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error { return recorder.Run(ctx) })
	g.Go(func() error { return decisions.Consume(ctx, gatekeeper.Update) })
	g.Go(func() error { return server.ListenAndServe(ctx) })
	return g.Wait()
}

// newRouter wires the traffic path: observe first (blocked requests are
// signal too), then the gate, then the upstream proxy. Health endpoints and
// the challenge flow stay outside the gate — a challenged session must be
// able to reach the verification endpoints.
func newRouter(cfg config.Config, recorder *app.Recorder, gatekeeper *app.Gatekeeper,
	limiter *app.RateLimiter, challenger *app.Challenger, log *slog.Logger) http.Handler {
	// Order: verify agent identity first (observe stamps the result on the
	// event), then observe, then the gate.
	traffic := httpx.WithMiddleware(proxy.New(cfg.UpstreamURL, log),
		agentauth.Middleware(app.NewAgentVerifier(cfg.AgentKeys)),
		observe.Middleware(recorder), gate.Middleware(gatekeeper, limiter, challenger, log))
	verification := challenge.New(challenger, log)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		httpx.JSON(w, http.StatusOK, map[string]string{"status": "alive"})
	})
	mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, _ *http.Request) {
		httpx.JSON(w, http.StatusOK, map[string]string{"status": "ready"})
	})
	mux.HandleFunc("GET "+challenge.PagePath, verification.Page)
	mux.HandleFunc("POST "+challenge.PagePath+"/verify", verification.Verify)
	mux.Handle("/", traffic)
	return httpx.WithMiddleware(mux, httpx.RequestID(), httpx.Logging(log), httpx.Recover(log))
}
