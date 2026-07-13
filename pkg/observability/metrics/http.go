// Package metrics provides a set of Prometheus metrics for monitoring HTTP requests in the Runlog API.
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// HTTPMetrics defines a set of Prometheus metrics for monitoring HTTP requests. It includes counters for total requests, histograms for request durations, and gauges for in-flight requests. These metrics can be used to track the performance and health of the Runlog API's HTTP endpoints.
type HTTPMetrics struct {
	RequestsTotal    *prometheus.CounterVec
	RequestDuration  *prometheus.HistogramVec
	RequestsInFlight prometheus.Gauge
}

// NewHTTPMetrics creates a new instance of HTTPMetrics and registers the metrics with the provided Prometheus registry. It initializes counters, histograms, and gauges for monitoring HTTP requests, allowing for tracking of request counts, durations, and in-flight requests in the Runlog API.
func NewHTTPMetrics(registry prometheus.Registerer) *HTTPMetrics {
	httpMetrics := &HTTPMetrics{
		RequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "runlog",
				Subsystem: "http",
				Name:      "requests_total",
				Help:      "Total number of HTTP requests.",
			},
			[]string{"method", "route", "status"},
		),
		RequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "runlog",
				Subsystem: "http",
				Name:      "request_duration_seconds",
				Help:      "Duration of HTTP requests in seconds.",
				Buckets: []float64{
					0.005,
					0.01,
					0.025,
					0.05,
					0.1,
					0.25,
					0.5,
					1,
					2.5,
					5,
				},
			},
			[]string{"method", "route", "status"},
		),
		RequestsInFlight: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: "runlog",
				Subsystem: "http",
				Name:      "requests_in_flight",
				Help:      "Current number of HTTP requests being processed.",
			},
		),
	}

	registry.MustRegister(
		httpMetrics.RequestsTotal,
		httpMetrics.RequestDuration,
		httpMetrics.RequestsInFlight,
	)

	return httpMetrics
}
