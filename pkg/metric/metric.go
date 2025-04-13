package metric

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metric struct {
	registry *prometheus.Registry

	cacheMiss *prometheus.CounterVec
	cacheHit  *prometheus.CounterVec
	errors    *prometheus.CounterVec

	httpRequestsTotal *prometheus.CounterVec
	httpStatusCodes   *prometheus.CounterVec
}

func New() *Metric {
	r := prometheus.NewRegistry()

	m := &Metric{
		registry: r,

		cacheHit: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "cache",
				Subsystem: "hit",
				Name:      "total",
				Help:      "Total number of cache hits",
			},
			[]string{"from"},
		),

		cacheMiss: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "cache",
				Subsystem: "miss",
				Name:      "total",
				Help:      "Total number of missing cache hits",
			},
			[]string{"from"},
		),

		errors: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "app",
				Subsystem: "errors",
				Name:      "total",
				Help:      "Total number of errors",
			},
			[]string{"service", "action"},
		),

		httpRequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "http",
				Subsystem: "requests",
				Name:      "total",
				Help:      "Total number of HTTP requests",
			},
			[]string{"method", "path", "status"},
		),

		httpStatusCodes: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "http",
				Subsystem: "status_codes",
				Name:      "total",
				Help:      "Total number of HTTP status codes",
			},
			[]string{"status"},
		),
	}

	r.MustRegister(
		m.cacheMiss,
		m.cacheHit,
		m.errors,
		m.httpRequestsTotal,
		m.httpStatusCodes,
	)

	return m
}

func (m *Metric) RecordCacheHit(from string) {
	m.cacheHit.WithLabelValues(from).Inc()
}

func (m *Metric) RecordCacheMiss(from string) {
	m.cacheMiss.WithLabelValues(from).Inc()
}

func (m *Metric) RecordError(service, action string) {
	m.errors.WithLabelValues(service, action).Inc()
}

func (m *Metric) Registry() *prometheus.Registry {
	return m.registry
}

func (m *Metric) RecordHTTPRequest(method, path, status string) {
	m.httpRequestsTotal.WithLabelValues(method, path, status).Inc()
	m.httpStatusCodes.WithLabelValues(status).Inc()
}
