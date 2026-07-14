package metrics

import (
	"context"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type queryTraceContextKey struct{}

type queryTraceData struct {
	startedAt time.Time
	operation string
}

// PGXTracer implements the pgx.QueryTracer interface to provide tracing for database queries. It tracks the start and end of queries, recording metrics such as total queries, query duration, and in-flight queries. The tracer uses Prometheus metrics to monitor database performance and behavior in a production environment.
type PGXTracer struct {
	metrics *DatabaseMetrics
}

// NewPGXTracer creates a new PGXTracer instance with the provided DatabaseMetrics. It initializes the tracer with the necessary metrics for monitoring database queries and returns the configured PGXTracer.
func NewPGXTracer(dbMetrics *DatabaseMetrics) *PGXTracer {
	return &PGXTracer{
		metrics: dbMetrics,
	}
}

// TraceQueryStart is called when a database query starts. It increments the in-flight queries metric and records the start time and operation type of the query. The function returns a new context containing the trace data, which will be used later in TraceQueryEnd to calculate query duration and update metrics.
func (t *PGXTracer) TraceQueryStart(
	ctx context.Context,
	_ *pgx.Conn,
	data pgx.TraceQueryStartData,
) context.Context {
	t.metrics.QueriesInFlight.Inc()

	traceData := queryTraceData{
		startedAt: time.Now(),
		operation: queryOperation(data.SQL),
	}

	return context.WithValue(ctx, queryTraceContextKey{}, traceData)
}

// TraceQueryEnd is called when a database query ends. It decrements the in-flight queries metric, retrieves the trace data from the context, and updates the total queries and query duration metrics based on the operation type and status (success or error). The function calculates the duration of the query and records it in the appropriate Prometheus histogram.
func (t *PGXTracer) TraceQueryEnd(
	ctx context.Context,
	_ *pgx.Conn,
	data pgx.TraceQueryEndData,
) {
	t.metrics.QueriesInFlight.Dec()

	traceData, ok := ctx.Value(queryTraceContextKey{}).(queryTraceData)
	if !ok {
		return
	}

	status := "success"
	if data.Err != nil {
		status = "error"
	}

	t.metrics.QueriesTotal.WithLabelValues(
		traceData.operation,
		status,
	).Inc()

	t.metrics.QueryDuration.WithLabelValues(
		traceData.operation,
		status,
	).Observe(time.Since(traceData.startedAt).Seconds())
}

func queryOperation(sql string) string {
	fields := strings.Fields(strings.TrimSpace(sql))
	if len(fields) == 0 {
		return "UNKNOWN"
	}

	switch strings.ToUpper(fields[0]) {
	case "SELECT", "INSERT", "UPDATE", "DELETE":
		return strings.ToUpper(fields[0])
	default:
		return "OTHER"
	}
}
