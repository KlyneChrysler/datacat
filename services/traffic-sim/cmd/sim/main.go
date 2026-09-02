// Command sim is the traffic-sim composition root: wiring and lifecycle
// only. It drives the persona registry against the target until SIGTERM or
// the configured duration.
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/KlyneChrysler/datacat/pkg/obsx"
	"github.com/KlyneChrysler/datacat/services/traffic-sim/internal/adapters/httpsender"
	"github.com/KlyneChrysler/datacat/services/traffic-sim/internal/app"
	"github.com/KlyneChrysler/datacat/services/traffic-sim/internal/config"
	"github.com/KlyneChrysler/datacat/services/traffic-sim/internal/domain"
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
	log := obsx.NewLogger("traffic-sim")

	ctx, cancel := withOptionalDeadline(ctx, cfg)
	defer cancel()

	simulator := app.NewSimulator(httpsender.NewSender(cfg.TargetURL), log, domain.DefaultPersonas())
	log.Info("starting", "target", cfg.TargetURL, "duration", cfg.Duration.String())
	return simulator.Run(ctx)
}

func withOptionalDeadline(ctx context.Context, cfg config.Config) (context.Context, context.CancelFunc) {
	if cfg.Duration <= 0 {
		return context.WithCancel(ctx)
	}
	return context.WithTimeout(ctx, cfg.Duration)
}
