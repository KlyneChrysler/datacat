// Package app holds edge-proxy use cases.
package app

import (
	"context"
	"log/slog"
	"sync/atomic"

	"github.com/KlyneChrysler/datacat/pkg/events"
	"github.com/KlyneChrysler/datacat/services/edge-proxy/internal/ports"
)

// Recorder decouples the hot request path from Kafka: Record never blocks
// and never fails a request; events flow to the publisher on a background
// goroutine. When the buffer is full the event is dropped and counted —
// observability must never take down the traffic it observes.
type Recorder struct {
	publisher ports.EventPublisher
	log       *slog.Logger
	queue     chan events.RequestEvent
	dropped   atomic.Int64
}

func NewRecorder(publisher ports.EventPublisher, log *slog.Logger, bufferSize int) *Recorder {
	return &Recorder{
		publisher: publisher,
		log:       log,
		queue:     make(chan events.RequestEvent, bufferSize),
	}
}

// Record enqueues an event without blocking the caller.
func (r *Recorder) Record(ev events.RequestEvent) {
	select {
	case r.queue <- ev:
	default:
		r.countDrop(ev)
	}
}

// Run publishes queued events until ctx is cancelled.
func (r *Recorder) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case ev := <-r.queue:
			r.publish(ctx, ev)
		}
	}
}

func (r *Recorder) publish(ctx context.Context, ev events.RequestEvent) {
	if err := r.publisher.PublishRequest(ctx, ev); err != nil && ctx.Err() == nil {
		r.log.ErrorContext(ctx, "publish request event failed", "err", err, "session_id", ev.SessionID)
	}
}

func (r *Recorder) countDrop(ev events.RequestEvent) {
	total := r.dropped.Add(1)
	r.log.Warn("event queue full; event dropped", "session_id", ev.SessionID, "dropped_total", total)
}
