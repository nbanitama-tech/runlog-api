package metrics

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
)

// PGXPoolCollector is a Prometheus collector that gathers metrics from a pgxpool.Pool instance. It collects information about the number of acquired, idle, total, constructing, and maximum connections in the connection pool, allowing for monitoring and analysis of database connection usage and performance in a production environment.
type PGXPoolCollector struct {
	pool *pgxpool.Pool

	acquiredConnections     *prometheus.Desc
	idleConnections         *prometheus.Desc
	totalConnections        *prometheus.Desc
	constructingConnections *prometheus.Desc
	maxConnections          *prometheus.Desc
}

// NewPGXPoolCollector creates a new PGXPoolCollector instance for the provided pgxpool.Pool. It initializes the Prometheus metric descriptors for acquired, idle, total, constructing, and maximum connections in the pool, allowing for monitoring and analysis of database connection usage and performance in a production environment.
func NewPGXPoolCollector(pool *pgxpool.Pool) *PGXPoolCollector {
	return &PGXPoolCollector{
		pool: pool,

		acquiredConnections: prometheus.NewDesc(
			"runlog_database_pool_connections_acquired",
			"Number of connections currently acquired from the pool.",
			nil,
			nil,
		),
		idleConnections: prometheus.NewDesc(
			"runlog_database_pool_connections_idle",
			"Number of currently idle connections in the pool.",
			nil,
			nil,
		),
		totalConnections: prometheus.NewDesc(
			"runlog_database_pool_connections_total",
			"Total number of connections currently in the pool.",
			nil,
			nil,
		),
		constructingConnections: prometheus.NewDesc(
			"runlog_database_pool_connections_constructing",
			"Number of connections currently being established.",
			nil,
			nil,
		),
		maxConnections: prometheus.NewDesc(
			"runlog_database_pool_connections_max",
			"Maximum number of connections allowed in the pool.",
			nil,
			nil,
		),
	}
}

// Describe sends the descriptors of the metrics collected by the PGXPoolCollector to the provided channel. It allows Prometheus to discover and understand the metrics available for collection from the pgxpool.Pool instance.
func (c *PGXPoolCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.acquiredConnections
	ch <- c.idleConnections
	ch <- c.totalConnections
	ch <- c.constructingConnections
	ch <- c.maxConnections
}

// Collect gathers the current values of the metrics from the pgxpool.Pool instance and sends them to the provided channel. It collects information about the number of acquired, idle, total, constructing, and maximum connections in the connection pool, allowing for monitoring and analysis of database connection usage and performance in a production environment.
func (c *PGXPoolCollector) Collect(ch chan<- prometheus.Metric) {
	stat := c.pool.Stat()

	ch <- prometheus.MustNewConstMetric(
		c.acquiredConnections,
		prometheus.GaugeValue,
		float64(stat.AcquiredConns()),
	)

	ch <- prometheus.MustNewConstMetric(
		c.idleConnections,
		prometheus.GaugeValue,
		float64(stat.IdleConns()),
	)

	ch <- prometheus.MustNewConstMetric(
		c.totalConnections,
		prometheus.GaugeValue,
		float64(stat.TotalConns()),
	)

	ch <- prometheus.MustNewConstMetric(
		c.constructingConnections,
		prometheus.GaugeValue,
		float64(stat.ConstructingConns()),
	)

	ch <- prometheus.MustNewConstMetric(
		c.maxConnections,
		prometheus.GaugeValue,
		float64(stat.MaxConns()),
	)
}
