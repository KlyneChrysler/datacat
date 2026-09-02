// Command server is the enforcement composition root.
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
	policy "github.com/KlyneChrysler/datacat/pkg/policy"
	"github.com/KlyneChrysler/datacat/services/enforcement/internal/adapters/actions"
	"github.com/KlyneChrysler/datacat/services/enforcement/internal/adapters/dynamo"
	"github.com/KlyneChrysler/datacat/services/enforcement/internal/adapters/httpapi"
	"github.com/KlyneChrysler/datacat/services/enforcement/internal/adapters/kafka"
	"github.com/KlyneChrysler/datacat/services/enforcement/internal/adapters/memory"
	"github.com/KlyneChrysler/datacat/services/enforcement/internal/app"
	"github.com/KlyneChrysler/datacat/services/enforcement/internal/config"
	"github.com/KlyneChrysler/datacat/services/enforcement/internal/ports"
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
	store, err := newDecisionStore(ctx, cfg, log)
	if err != nil {
		return err
	}

	publisher := kafka.NewDecisionPublisher(producer, cfg.DecisionsTopic)
	enforcer := app.NewEnforcer(policy.DefaultPolicy(), store, publisher, actions.NewLogApplier(log), app.NewTally())
	source := kafka.NewVerdictSource(consumer)
	router := httpapi.NewRouter(httpapi.New(enforcer, log), log, cfg.CORSOrigin)
	server := httpx.NewServer(cfg.Port, router, cfg.ShutdownTimeout)

	log.Info("starting", "port", cfg.Port, "topic", cfg.VerdictsTopic, "group", cfg.ConsumerGroup)

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error { return source.Consume(ctx, enforcer.HandleVerdict) })
	g.Go(func() error { return server.ListenAndServe(ctx) })

	return g.Wait()
}

// newDecisionStore selects DynamoDB when a table is configured.
func newDecisionStore(ctx context.Context, cfg config.Config, log *slog.Logger) (ports.DecisionStore, error) {
	if cfg.DecisionsTable == "" {
		log.Info("decision store: in-memory (single replica only)")
		return memory.NewDecisionStore(), nil
	}

	client, err := dynamo.NewClient(ctx, cfg.DynamoEndpoint)
	if err != nil {
		return nil, err
	}

	log.Info("decision store: dynamodb", "table", cfg.DecisionsTable)
	return dynamo.NewDecisionStore(client, cfg.DecisionsTable, cfg.DecisionTTL), nil
}
