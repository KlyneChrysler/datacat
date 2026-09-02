package app

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/KlyneChrysler/datacat/pkg/events"
	"github.com/KlyneChrysler/datacat/pkg/obsx"
)

type fakePublisher struct {
	mu        sync.Mutex
	published []events.RequestEvent
}

func (f *fakePublisher) PublishRequest(_ context.Context, ev events.RequestEvent) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.published = append(f.published, ev)
	return nil
}

func (f *fakePublisher) count() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.published)
}

func TestRecorderPublishesQueuedEvents(t *testing.T) {
	publisher := &fakePublisher{}
	recorder := NewRecorder(publisher, obsx.NewLogger("test"), 8)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() { _ = recorder.Run(ctx) }()

	recorder.Record(events.RequestEvent{SessionID: "s-1"})
	recorder.Record(events.RequestEvent{SessionID: "s-2"})

	waitFor(t, func() bool { return publisher.count() == 2 })
}

func TestRecorderDropsWhenQueueFull(t *testing.T) {
	recorder := NewRecorder(&fakePublisher{}, obsx.NewLogger("test"), 1)
	// Run is never started, so the queue never drains.
	recorder.Record(events.RequestEvent{SessionID: "s-1"})
	recorder.Record(events.RequestEvent{SessionID: "s-2"}) // must not block

	if got := recorder.dropped.Load(); got != 1 {
		t.Fatalf("dropped = %d, want 1", got)
	}
}

func waitFor(t *testing.T, condition func() bool) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if condition() {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatal("condition not met within deadline")
}
