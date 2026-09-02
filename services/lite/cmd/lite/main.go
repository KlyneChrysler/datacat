// Command lite is the single binary bouncer, no Kafka, no Flink.
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/KlyneChrysler/datacat/pkg/edge/agentauth"
	"github.com/KlyneChrysler/datacat/pkg/edge/challenge"
	"github.com/KlyneChrysler/datacat/pkg/edge/gate"
	"github.com/KlyneChrysler/datacat/pkg/edge/proxy"
	"github.com/KlyneChrysler/datacat/pkg/guard"
	"github.com/KlyneChrysler/datacat/pkg/httpx"
	"github.com/KlyneChrysler/datacat/pkg/obsx"
	policy "github.com/KlyneChrysler/datacat/pkg/policy"
	"github.com/KlyneChrysler/datacat/services/lite/internal/adapters/watch"
	"github.com/KlyneChrysler/datacat/services/lite/internal/app"
	"github.com/KlyneChrysler/datacat/services/lite/internal/classify"
	"github.com/KlyneChrysler/datacat/services/lite/internal/config"
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
	log := obsx.NewLogger("datacat-lite")

	gatekeeper := guard.NewGatekeeper(cfg.GateTTL)
	limiter := guard.NewRateLimiter(cfg.RateLimitPerMinute, cfg.RateLimitBurst, cfg.GateTTL)
	challenger := guard.NewChallenger(cfg.ChallengeSecret, cfg.ChallengeDifficulty)
	analyzer := app.NewAnalyzer(classify.NewTracker(), policy.DefaultPolicy(), gatekeeper, log)
	traffic := gate.New(gatekeeper, limiter, challenger, log)
	server := httpx.NewServer(cfg.Port, newRouter(cfg, analyzer, traffic, challenger, log), cfg.ShutdownTimeout)

	log.Info("starting", "port", cfg.Port, "upstream", cfg.UpstreamURL.String())

	return server.ListenAndServe(ctx)
}

// newRouter wires watch then gate then proxy, challenge stays outside the gate.
func newRouter(cfg config.Config, analyzer *app.Analyzer, traffic *gate.Gate, challenger *guard.Challenger, log *slog.Logger) http.Handler {
	agentCheck := agentauth.Middleware(guard.NewAgentVerifier(cfg.AgentKeys))
	proxied := httpx.WithMiddleware(proxy.New(cfg.UpstreamURL, log), agentCheck, watch.Middleware(analyzer), traffic.Middleware())
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
	mux.Handle("/", proxied)

	return httpx.WithMiddleware(mux, httpx.RequestID(), httpx.Logging(log), httpx.Recover(log))
}
