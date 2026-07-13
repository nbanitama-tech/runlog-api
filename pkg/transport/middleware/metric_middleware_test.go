package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nbanitama-tech/runlog-api/pkg/observability/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestMetricsMiddleware_RecordsRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	registry := prometheus.NewRegistry()
	httpMetrics := metrics.NewHTTPMetrics(registry)

	router := gin.New()
	router.Use(MetricsMiddleware(httpMetrics))
	router.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	count := testutil.ToFloat64(
		httpMetrics.RequestsTotal.WithLabelValues(
			http.MethodGet,
			"/health",
			"200",
		),
	)

	if count != 1 {
		t.Fatalf("expected request count 1, got %v", count)
	}
}
