package obsx

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMetricsHandlerExposesCounts(t *testing.T) {
	metrics := NewTestMetrics()

	metrics.CountVerdict("abusive")
	metrics.CountVerdict("abusive")
	metrics.CountAction("block")

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()
	metrics.Handler().ServeHTTP(rec, req)

	body := rec.Body.String()
	if !strings.Contains(body, `datacat_verdicts_total{classification="abusive",service="test"} 2`) {
		t.Errorf("verdict count missing from scrape:\n%s", body)
	}
	if !strings.Contains(body, `datacat_actions_total{action="block",service="test"} 1`) {
		t.Errorf("action count missing from scrape:\n%s", body)
	}
}
