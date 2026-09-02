// Command sim generates synthetic traffic against a target: human-like
// (jittery, bursty), polite-agent (declared identity, steady cadence), and
// scraper (fast, regular, deep pagination) personas. Personas are implemented
// in the traffic-sim phase; this skeleton validates config and lifecycle.
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/KlyneChrysler/datacat/pkg/obsx"
	"github.com/KlyneChrysler/datacat/services/traffic-sim/internal/config"
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

	log.Info("simulator skeleton ready", "target", cfg.TargetURL)
	<-ctx.Done()
	log.Info("shutting down")
	return nil
}
