// Command server is the enforcement composition root: wiring and lifecycle
// only. Concrete types meet here and nowhere else.
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/KlyneChrysler/datacat/pkg/httpx"
	"github.com/KlyneChrysler/datacat/pkg/obsx"
	"github.com/KlyneChrysler/datacat/services/enforcement/internal/adapters/httpapi"
	"github.com/KlyneChrysler/datacat/services/enforcement/internal/config"
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
	log := obsx.NewLogger("enforcement")

	router := httpapi.NewRouter(httpapi.New(), log)
	server := httpx.NewServer(cfg.Port, router, cfg.ShutdownTimeout)

	log.Info("starting", "port", cfg.Port)
	return server.ListenAndServe(ctx)
}
