package httpapi

import (
	"net/http"

	"github.com/KlyneChrysler/datacat/pkg/httpx"
)

type Handlers struct{}

func New() *Handlers {
	return &Handlers{}
}

// Health reports the process is alive (liveness probe).
func (h *Handlers) Health(w http.ResponseWriter, _ *http.Request) {
	httpx.JSON(w, http.StatusOK, map[string]string{"status": "alive"})
}

// Ready reports dependencies are reachable (readiness probe). Dependency
// checks are added as adapters arrive (Kafka, DynamoDB).
func (h *Handlers) Ready(w http.ResponseWriter, _ *http.Request) {
	httpx.JSON(w, http.StatusOK, map[string]string{"status": "ready"})
}
