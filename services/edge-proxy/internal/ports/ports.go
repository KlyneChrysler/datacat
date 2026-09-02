// Package ports holds the interfaces edge-proxy consumes. Implementations
// live in adapters and meet these contracts only in the composition root.
package ports

import (
	"context"

	"github.com/KlyneChrysler/datacat/pkg/events"
)

type EventPublisher interface {
	PublishRequest(ctx context.Context, ev events.RequestEvent) error
}

type DecisionSource interface {
	// Consume blocks, invoking handle for each decision until ctx is done.
	Consume(ctx context.Context, handle func(events.Decision)) error
}
