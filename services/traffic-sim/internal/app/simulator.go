// Package app holds traffic-sim use cases.
package app

import (
	"context"
	"log/slog"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/KlyneChrysler/datacat/services/traffic-sim/internal/domain"
	"github.com/KlyneChrysler/datacat/services/traffic-sim/internal/ports"
)

// Simulator drives every persona concurrently until the context ends.
// Status transitions are logged (200 → 429 → 403 is the demo's story);
// send errors are logged and the persona keeps going — a refusing target
// is a result, not a failure.
type Simulator struct {
	sender   ports.Sender
	log      *slog.Logger
	personas []domain.Persona
}

func NewSimulator(sender ports.Sender, log *slog.Logger, personas []domain.Persona) *Simulator {
	return &Simulator{sender: sender, log: log, personas: personas}
}

// Run blocks until ctx is done; one goroutine per persona (bounded by the
// registry size).
func (s *Simulator) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)
	for _, persona := range s.personas {
		g.Go(func() error { return s.runPersona(ctx, persona) })
	}
	return g.Wait()
}

func (s *Simulator) runPersona(ctx context.Context, persona domain.Persona) error {
	lastStatus := 0
	for {
		if ctx.Err() != nil {
			return nil
		}
		lastStatus = s.sendOnce(ctx, persona, lastStatus)
		if err := sleepCtx(ctx, persona.NextDelay()); err != nil {
			return nil
		}
	}
}

// sendOnce performs one request and returns the status for transition
// tracking (unchanged status is not re-logged).
func (s *Simulator) sendOnce(ctx context.Context, persona domain.Persona, lastStatus int) int {
	status, err := s.sender.Send(ctx, persona.NextRequest())
	if err != nil {
		if ctx.Err() == nil {
			s.log.WarnContext(ctx, "send failed", "persona", persona.Name, "err", err)
		}
		return lastStatus
	}
	if status != lastStatus {
		s.log.InfoContext(ctx, "status transition", "persona", persona.Name, "status", status)
	}
	return status
}

func sleepCtx(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
