package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nbanitama-tech/runlog-api/pkg/observability/metrics"
)

// MetricsMiddleware is a Gin middleware that collects Prometheus metrics for HTTP requests. It tracks the number of in-flight requests, request durations, and total requests by method, route, and status code. The collected metrics are exposed to Prometheus for monitoring and analysis.
func MetricsMiddleware(httpMetrics *metrics.HTTPMetrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		httpMetrics.RequestsInFlight.Inc()
		defer httpMetrics.RequestsInFlight.Dec()

		startedAt := time.Now()

		c.Next()

		route := c.FullPath()
		if route == "" {
			route = "unmatched"
		}

		if c.FullPath() == "/metrics" {
			c.Next()
			return
		}

		status := strconv.Itoa(c.Writer.Status())

		httpMetrics.RequestsTotal.WithLabelValues(
			c.Request.Method,
			route,
			status,
		).Inc()

		httpMetrics.RequestDuration.WithLabelValues(
			c.Request.Method,
			route,
			status,
		).Observe(time.Since(startedAt).Seconds())
	}
}
