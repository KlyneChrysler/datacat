package obsx

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics holds the datacat business counters, labelled by class or action.
type Metrics struct {
	registry *prometheus.Registry
	verdicts *prometheus.CounterVec
	actions  *prometheus.CounterVec
}

// NewMetrics builds the counters, registering them on the given registry.
func NewMetrics(service string, registry *prometheus.Registry) *Metrics {
	labels := prometheus.Labels{"service": service}

	verdicts := prometheus.NewCounterVec(prometheus.CounterOpts{Name: "datacat_verdicts_total", Help: "Verdicts by classification", ConstLabels: labels}, []string{"classification"})
	actions := prometheus.NewCounterVec(prometheus.CounterOpts{Name: "datacat_actions_total", Help: "Enforcement actions applied", ConstLabels: labels}, []string{"action"})

	registry.MustRegister(verdicts, actions)

	return &Metrics{registry: registry, verdicts: verdicts, actions: actions}
}

// NewTestMetrics builds metrics on a throwaway registry for tests.
func NewTestMetrics() *Metrics {
	return NewMetrics("test", prometheus.NewRegistry())
}

// CountVerdict records one verdict of a classification.
func (m *Metrics) CountVerdict(classification string) {
	m.verdicts.WithLabelValues(classification).Inc()
}

// CountAction records one applied enforcement action.
func (m *Metrics) CountAction(action string) {
	m.actions.WithLabelValues(action).Inc()
}

// Handler serves the Prometheus scrape endpoint for this registry.
func (m *Metrics) Handler() http.Handler {
	return promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{})
}
