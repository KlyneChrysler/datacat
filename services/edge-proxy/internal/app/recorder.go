// Package app holds edge proxy use cases.
package app

import (
	"context"
	"log/slog"
	"sync/atomic"

	"github.com/KlyneChrysler/datacat/pkg/events"
	"github.com/KlyneChrysler/datacat/pkg/obsx"
	"github.com/KlyneChrysler/datacat/services/edge-proxy/internal/ports"
)

// Recorder queues events without ever blocking or failing a request.
type Recorder struct {
	publisher ports.EventPublisher
	log       *slog.Logger
	metrics   *obsx.Metrics
	queue     chan events.RequestEvent
	dropped   atomic.Int64
}

func NewRecorder(publisher ports.EventPublisher, log *slog.Logger, metrics *obsx.Metrics, bufferSize int) *Recorder {
	return &Recorder{publisher: publisher, log: log, metrics: metrics, queue: make(chan events.RequestEvent, bufferSize)}
}

// Record enqueues an event, dropping it when the queue is full.
func (r *Recorder) Record(ev events.RequestEvent) {
	r.metrics.CountRequest(ev.VerifiedAgent)

	select {
	case r.queue <- ev:
	default:
		r.countDrop(ev)
	}
}

// Run publishes queued events until ctx ends.
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

	r.log.Warn("event queue full, event dropped", "session_id", ev.SessionID, "dropped_total", total)
}
