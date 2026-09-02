// Package ports holds the interfaces edge proxy consumes.
package ports

import (
	"context"

	"github.com/KlyneChrysler/datacat/pkg/events"
)

type EventPublisher interface {
	PublishRequest(ctx context.Context, ev events.RequestEvent) error
}

type DecisionSource interface {
	// Consume handles each decision until ctx ends.
	Consume(ctx context.Context, handle func(events.Decision)) error
}
