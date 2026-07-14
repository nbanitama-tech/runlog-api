package metrics

import "github.com/prometheus/client_golang/prometheus"

// DatabaseMetrics defines the Prometheus metrics for database operations, including total queries, query duration, and in-flight queries. It provides a structured way to monitor and analyze database performance and behavior in a production environment.
type DatabaseMetrics struct {
	QueriesTotal    *prometheus.CounterVec
	QueryDuration   *prometheus.HistogramVec
	QueriesInFlight prometheus.Gauge
}

// NewDatabaseMetrics creates and registers Prometheus metrics for database operations. It initializes counters for total queries, histograms for query duration, and gauges for in-flight queries. The metrics are registered with the provided Prometheus registry, allowing them to be collected and monitored in a production environment.
func NewDatabaseMetrics(registry prometheus.Registerer) *DatabaseMetrics {
	dbMetrics := &DatabaseMetrics{
		QueriesTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "runlog",
				Subsystem: "database",
				Name:      "queries_total",
				Help:      "Total number of database queries.",
			},
			[]string{"operation", "status"},
		),
		QueryDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "runlog",
				Subsystem: "database",
				Name:      "query_duration_seconds",
				Help:      "Duration of database queries in seconds.",
				Buckets: []float64{
					0.001,
					0.0025,
					0.005,
					0.01,
					0.025,
					0.05,
					0.1,
					0.25,
					0.5,
					1,
					2.5,
				},
			},
			[]string{"operation", "status"},
		),
		QueriesInFlight: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: "runlog",
				Subsystem: "database",
				Name:      "queries_in_flight",
				Help:      "Current number of database queries being executed.",
			},
		),
	}

	registry.MustRegister(
		dbMetrics.QueriesTotal,
		dbMetrics.QueryDuration,
		dbMetrics.QueriesInFlight,
	)

	return dbMetrics
}
