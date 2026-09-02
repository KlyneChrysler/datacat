// Package ports holds the interfaces traffic-sim consumes.
package ports

import (
	"context"

	"github.com/KlyneChrysler/datacat/services/traffic-sim/internal/domain"
)

type Sender interface {
	// Send performs one request and returns the HTTP status received.
	Send(ctx context.Context, req domain.Request) (int, error)
}
