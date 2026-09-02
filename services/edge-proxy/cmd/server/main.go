// Command server is the edge proxy composition root.
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/KlyneChrysler/datacat/pkg/edge/agentauth"
	"github.com/KlyneChrysler/datacat/pkg/edge/challenge"
	"github.com/KlyneChrysler/datacat/pkg/edge/gate"
	"github.com/KlyneChrysler/datacat/pkg/edge/proxy"
	"github.com/KlyneChrysler/datacat/pkg/events"
	"github.com/KlyneChrysler/datacat/pkg/guard"
	"github.com/KlyneChrysler/datacat/pkg/httpx"
	"github.com/KlyneChrysler/datacat/pkg/kafkax"
	"github.com/KlyneChrysler/datacat/pkg/obsx"
	"github.com/KlyneChrysler/datacat/services/edge-proxy/internal/adapters/kafka"
	"github.com/KlyneChrysler/datacat/services/edge-proxy/internal/adapters/observe"
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

	recorder := app.NewRecorder(kafka.NewEventPublisher(producer, cfg.RequestsTopic), log, cfg.EventBufferSize)
	gatekeeper := guard.NewGatekeeper(cfg.GateTTL)
	limiter := guard.NewRateLimiter(cfg.RateLimitPerMinute, cfg.RateLimitBurst, cfg.GateTTL)
	challenger := guard.NewChallenger(cfg.ChallengeSecret, cfg.ChallengeDifficulty)
	traffic := gate.New(gatekeeper, limiter, challenger, log)
	decisions := kafka.NewDecisionSource(consumer)
	server := httpx.NewServer(cfg.Port, newRouter(cfg, recorder, traffic, challenger, log), cfg.ShutdownTimeout)

	log.Info("starting", "port", cfg.Port, "upstream", cfg.UpstreamURL.String(), "topic", cfg.RequestsTopic)

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error { return recorder.Run(ctx) })
	g.Go(func() error {
		return decisions.Consume(ctx, func(d events.Decision) { gatekeeper.Update(d.SessionID, d.Action) })
	})
	g.Go(func() error { return server.ListenAndServe(ctx) })

	return g.Wait()
}

// newRouter wires the traffic path, health and challenge stay outside the gate.
func newRouter(cfg config.Config, recorder *app.Recorder, traffic *gate.Gate, challenger *guard.Challenger, log *slog.Logger) http.Handler {
	agentCheck := agentauth.Middleware(guard.NewAgentVerifier(cfg.AgentKeys))
	proxied := httpx.WithMiddleware(proxy.New(cfg.UpstreamURL, log), agentCheck, observe.Middleware(recorder), traffic.Middleware())
	verification := challenge.New(challenger, log)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", httpx.Alive)
	mux.HandleFunc("GET /readyz", httpx.Ready)
	verification.Mount(mux)
	mux.Handle("/", proxied)

	return httpx.WithMiddleware(mux, httpx.RequestID(), httpx.Logging(log), httpx.Recover(log))
}
