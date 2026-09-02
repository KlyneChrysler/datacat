// Command server is the enforcement composition root: wiring and lifecycle
// only. Concrete types meet here and nowhere else.
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/KlyneChrysler/datacat/pkg/httpx"
	"github.com/KlyneChrysler/datacat/pkg/kafkax"
	"github.com/KlyneChrysler/datacat/pkg/obsx"
	"github.com/KlyneChrysler/datacat/services/enforcement/internal/adapters/actions"
	"github.com/KlyneChrysler/datacat/services/enforcement/internal/adapters/httpapi"
	"github.com/KlyneChrysler/datacat/services/enforcement/internal/adapters/kafka"
	"github.com/KlyneChrysler/datacat/services/enforcement/internal/adapters/memory"
	"github.com/KlyneChrysler/datacat/services/enforcement/internal/app"
	"github.com/KlyneChrysler/datacat/services/enforcement/internal/config"
	"github.com/KlyneChrysler/datacat/services/enforcement/internal/domain"
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

	consumer, err := kafkax.NewConsumer(cfg.KafkaBrokers, cfg.ConsumerGroup, cfg.VerdictsTopic, log)
	if err != nil {
		return err
	}
	producer, err := kafkax.NewProducer(cfg.KafkaBrokers)
	if err != nil {
		return err
	}
	defer producer.Close()

	publisher := kafka.NewDecisionPublisher(producer, cfg.DecisionsTopic)
	enforcer := app.NewEnforcer(domain.DefaultPolicy(), memory.NewDecisionStore(), publisher, actions.NewLogApplier(log))
	source := kafka.NewVerdictSource(consumer)
	server := httpx.NewServer(cfg.Port, httpapi.NewRouter(httpapi.New(enforcer, log), log), cfg.ShutdownTimeout)

	log.Info("starting", "port", cfg.Port, "topic", cfg.VerdictsTopic, "group", cfg.ConsumerGroup)
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error { return source.Consume(ctx, enforcer.HandleVerdict) })
	g.Go(func() error { return server.ListenAndServe(ctx) })
	return g.Wait()
}
