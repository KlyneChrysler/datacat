// Package actions applies enforcement decisions. This phase logs them; the
// edge-proxy callback (actually blocking/challenging traffic) arrives when
// the verdict loop closes end to end.
package actions

import (
	"context"
	"log/slog"

	"github.com/KlyneChrysler/datacat/services/enforcement/internal/domain"
	"github.com/KlyneChrysler/datacat/services/enforcement/internal/ports"
)

type LogApplier struct {
	log *slog.Logger
}

var _ ports.ActionApplier = (*LogApplier)(nil)

func NewLogApplier(log *slog.Logger) *LogApplier {
	return &LogApplier{log: log}
}

func (a *LogApplier) Apply(ctx context.Context, d domain.Decision) error {
	a.log.InfoContext(ctx, "decision applied",
		"session_id", d.SessionID, "class", d.Class, "action", d.Action)
	return nil
}
