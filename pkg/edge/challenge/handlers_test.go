package challenge

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/KlyneChrysler/datacat/pkg/edge/ident"
	"github.com/KlyneChrysler/datacat/pkg/guard"
	"github.com/KlyneChrysler/datacat/pkg/obsx"
)

func testHandlers() (*Handlers, *guard.Challenger) {
	challenger := guard.NewChallenger("test-secret", 4)
	return New(challenger, obsx.NewLogger("test")), challenger
}

func sessionRequest(method, target, body string) *http.Request {
	r := httptest.NewRequest(method, target, strings.NewReader(body))
	r.AddCookie(&http.Cookie{Name: ident.SessionCookie, Value: "s-1"})
	return r
}

func TestPageServesTokenAndSolver(t *testing.T) {
	handlers, _ := testHandlers()
	w := httptest.NewRecorder()

	handlers.Page(w, sessionRequest(http.MethodGet, PagePath+"?return=/products", ""))

	body := w.Body.String()
	if w.Code != http.StatusOK || !strings.Contains(body, "crypto.subtle.digest") {
		t.Fatalf("status = %d, solver present = %v", w.Code, strings.Contains(body, "crypto.subtle"))
	}
	if !strings.Contains(body, "/products") {
		t.Error("return path missing from page")
	}
}

func TestPageRejectsForeignReturnURLs(t *testing.T) {
	handlers, _ := testHandlers()
	w := httptest.NewRecorder()

	handlers.Page(w, sessionRequest(http.MethodGet, PagePath+"?return=//evil.example", ""))

	if strings.Contains(w.Body.String(), "evil.example") {
		t.Error("protocol-relative return URL survived sanitization (open redirect)")
	}
}

func TestVerifyIssuesClearanceForSolvedChallenge(t *testing.T) {
	handlers, challenger := testHandlers()
	token := challenger.MintChallenge("s-1", time.Now())
	nonce := solve(t, challenger, token)
	w := httptest.NewRecorder()

	body := fmt.Sprintf(`{"token":%q,"nonce":%q}`, token, nonce)
	handlers.Verify(w, sessionRequest(http.MethodPost, PagePath+"/verify", body))

	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204", w.Code)
	}
	if !hasValidClearance(w.Result().Cookies(), challenger) {
		t.Error("no valid clearance cookie issued")
	}
}

func TestVerifyRejectsWrongNonce(t *testing.T) {
	handlers, challenger := testHandlers()
	token := challenger.MintChallenge("s-1", time.Now())
	w := httptest.NewRecorder()

	body := fmt.Sprintf(`{"token":%q,"nonce":"999999999"}`, token)
	handlers.Verify(w, sessionRequest(http.MethodPost, PagePath+"/verify", body))

	if w.Code != http.StatusForbidden {
		t.Errorf("status = %d, want 403", w.Code)
	}
}

func solve(t *testing.T, challenger *guard.Challenger, token string) string {
	t.Helper()
	now := time.Now()
	for n := 0; n < 1_000_000; n++ {
		nonce := strconv.Itoa(n)
		if challenger.VerifySolution(token, nonce, "s-1", now) {
			return nonce
		}
	}
	t.Fatal("no nonce found at test difficulty")
	return ""
}

func hasValidClearance(cookies []*http.Cookie, challenger *guard.Challenger) bool {
	for _, cookie := range cookies {
		if cookie.Name == guard.ClearanceCookie &&
			challenger.ValidClearance(cookie.Value, "s-1", time.Now()) {
			return true
		}
	}
	return false
}
