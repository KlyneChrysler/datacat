// Command sim is the traffic sim composition root.
package main

import (
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
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

	personas := domain.DefaultPersonas(agentCredential(cfg, log))
	simulator := app.NewSimulator(httpsender.NewSender(cfg.TargetURL), log, personas)

	log.Info("starting", "target", cfg.TargetURL, "duration", cfg.Duration.String())

	return simulator.Run(ctx)
}

// agentCredential derives the signing key from the seed, nil when unset.
func agentCredential(cfg config.Config, log *slog.Logger) *domain.AgentCredential {
	if cfg.AgentKeySeed == "" {
		return nil
	}

	seed := sha256.Sum256([]byte(cfg.AgentKeySeed))
	key := ed25519.NewKeyFromSeed(seed[:])

	log.Info("agent signing enabled", "key_id", "sim-agent-key", "public_key", hex.EncodeToString(key.Public().(ed25519.PublicKey)))

	return &domain.AgentCredential{KeyID: "sim-agent-key", Key: key}
}

func withOptionalDeadline(ctx context.Context, cfg config.Config) (context.Context, context.CancelFunc) {
	if cfg.Duration <= 0 {
		return context.WithCancel(ctx)
	}

	return context.WithTimeout(ctx, cfg.Duration)
}
