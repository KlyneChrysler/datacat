package app

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/KlyneChrysler/datacat/pkg/obsx"
	"github.com/KlyneChrysler/datacat/services/traffic-sim/internal/domain"
)

type fakeSender struct {
	mu       sync.Mutex
	byActor  map[string]int
	statuses []int
}

func (f *fakeSender) Send(_ context.Context, req domain.Request) (int, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.byActor[req.SessionID]++
	if len(f.statuses) == 0 {
		return 200, nil
	}
	status := f.statuses[0]
	f.statuses = f.statuses[1:]
	return status, nil
}

func (f *fakeSender) count(sessionID string) int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.byActor[sessionID]
}

func fastPersona(name string) domain.Persona {
	return domain.Persona{
		Name:      name,
		SessionID: name,
		UserAgent: "test/1.0",
		BaseDelay: time.Millisecond,
		Paths:     domain.NewCrawlingPaths("/p-"),
	}
}

func TestSimulatorDrivesEveryPersonaUntilContextEnds(t *testing.T) {
	sender := &fakeSender{byActor: make(map[string]int)}
	sim := NewSimulator(sender, obsx.NewLogger("test"),
		[]domain.Persona{fastPersona("a"), fastPersona("b")})
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	if err := sim.Run(ctx); err != nil {
		t.Fatalf("run returned %v, want nil on context end", err)
	}
	if sender.count("a") == 0 || sender.count("b") == 0 {
		t.Errorf("persona counts = a:%d b:%d, want both > 0", sender.count("a"), sender.count("b"))
	}
}
