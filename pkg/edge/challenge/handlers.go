// Package challenge serves the verification page and issues clearances.
package challenge

import (
	_ "embed"
	"encoding/json"
	"html/template"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/KlyneChrysler/datacat/pkg/edge/ident"
	"github.com/KlyneChrysler/datacat/pkg/guard"
	"github.com/KlyneChrysler/datacat/pkg/httpx"
)

// PagePath is where the gate redirects challenged browsers.
const PagePath = "/__datacat/challenge"

//go:embed page.html
var pageHTML string

var pageTemplate = template.Must(template.New("challenge").Parse(pageHTML))

type Handlers struct {
	challenger *guard.Challenger
	log        *slog.Logger
}

func New(challenger *guard.Challenger, log *slog.Logger) *Handlers {
	return &Handlers{challenger: challenger, log: log}
}

// Page serves the interstitial with a fresh session bound token.
func (h *Handlers) Page(w http.ResponseWriter, r *http.Request) {
	data := pageData{
		Token:      h.challenger.MintChallenge(ident.SessionID(r), time.Now()),
		Difficulty: h.challenger.Difficulty(),
		ReturnTo:   sanitizeReturn(r.URL.Query().Get("return")),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")

	if err := pageTemplate.Execute(w, data); err != nil {
		h.log.ErrorContext(r.Context(), "render challenge page", "err", err)
	}
}

// Verify accepts a solved proof and sets the clearance cookie.
func (h *Handlers) Verify(w http.ResponseWriter, r *http.Request) {
	var req verifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "malformed request")
		return
	}

	sessionID := ident.SessionID(r)
	now := time.Now()
	if !h.challenger.VerifySolution(req.Token, req.Nonce, sessionID, now) {
		h.log.InfoContext(r.Context(), "challenge solution rejected", "session_id", sessionID)
		httpx.Error(w, http.StatusForbidden, "invalid solution")
		return
	}

	http.SetCookie(w, clearanceCookie(h.challenger.MintClearance(sessionID, now)))
	h.log.InfoContext(r.Context(), "challenge passed", "session_id", sessionID)
	w.WriteHeader(http.StatusNoContent)
}

func clearanceCookie(value string) *http.Cookie {
	return &http.Cookie{Name: guard.ClearanceCookie, Value: value, Path: "/", MaxAge: 3600, HttpOnly: true, SameSite: http.SameSiteLaxMode}
}

// sanitizeReturn allows same origin paths only, everything else becomes /.
func sanitizeReturn(raw string) string {
	if strings.ContainsAny(raw, "\\\r\n\t") {
		return "/"
	}

	parsed, err := url.Parse(raw)
	if err != nil || parsed.Scheme != "" || parsed.Host != "" || !strings.HasPrefix(parsed.Path, "/") || strings.HasPrefix(parsed.Path, "//") {
		return "/"
	}

	return parsed.RequestURI()
}
