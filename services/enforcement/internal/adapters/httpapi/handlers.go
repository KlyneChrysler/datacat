package httpapi

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/KlyneChrysler/datacat/pkg/httpx"
	"github.com/KlyneChrysler/datacat/services/enforcement/internal/app"
	"github.com/KlyneChrysler/datacat/services/enforcement/internal/domain"
)

type Handlers struct {
	enforcer *app.Enforcer
	log      *slog.Logger
}

func New(enforcer *app.Enforcer, log *slog.Logger) *Handlers {
	return &Handlers{enforcer: enforcer, log: log}
}

// Health answers the liveness probe.
func (h *Handlers) Health(w http.ResponseWriter, _ *http.Request) {
	httpx.JSON(w, http.StatusOK, map[string]string{"status": "alive"})
}

// Ready answers the readiness probe.
func (h *Handlers) Ready(w http.ResponseWriter, _ *http.Request) {
	httpx.JSON(w, http.StatusOK, map[string]string{"status": "ready"})
}

// GetTrafficSummary serves classification counts for a recent window.
func (h *Handlers) GetTrafficSummary(w http.ResponseWriter, r *http.Request) {
	windowMinutes := parseWindowMinutes(r.URL.Query().Get("windowMinutes"))
	summary := h.enforcer.TrafficSummary(windowMinutes)

	httpx.JSON(w, http.StatusOK, toTrafficSummaryResponse(summary))
}

// GetDecision serves the standing decision for one session.
func (h *Handlers) GetDecision(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("sessionID")
	decision, err := h.enforcer.Lookup(r.Context(), sessionID)

	switch {
	case errors.Is(err, domain.ErrDecisionNotFound):
		httpx.Error(w, http.StatusNotFound, "no decision for session")
	case err != nil:
		httpx.InternalError(w, r, h.log, err)
	default:
		httpx.JSON(w, http.StatusOK, toDecisionResponse(decision))
	}
}

// parseWindowMinutes falls back to the default window on bad input.
func parseWindowMinutes(raw string) int {
	const defaultWindow = 15

	minutes, err := strconv.Atoi(raw)
	if err != nil {
		return defaultWindow
	}

	return minutes
}
