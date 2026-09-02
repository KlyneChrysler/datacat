package httpx

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func JSON(w http.ResponseWriter, status int, data any) {
	write(w, status, envelope{OK: true, Data: data})
}

func Error(w http.ResponseWriter, status int, message string) {
	write(w, status, envelope{OK: false, Error: message})
}

// InternalError logs the detail server-side and returns a generic message,
// never leaking internals to clients.
func InternalError(w http.ResponseWriter, r *http.Request, log *slog.Logger, err error) {
	log.ErrorContext(r.Context(), "internal error", "err", err, "path", r.URL.Path)
	Error(w, http.StatusInternalServerError, "internal error")
}

func write(w http.ResponseWriter, status int, body envelope) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body) // headers already sent; nothing left to do
}
