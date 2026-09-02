package httpx

import "net/http"

// Alive answers liveness probes.
func Alive(w http.ResponseWriter, _ *http.Request) {
	JSON(w, http.StatusOK, map[string]string{"status": "alive"})
}

// Ready answers readiness probes.
func Ready(w http.ResponseWriter, _ *http.Request) {
	JSON(w, http.StatusOK, map[string]string{"status": "ready"})
}
