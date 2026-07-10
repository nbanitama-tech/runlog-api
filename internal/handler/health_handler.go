package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nbanitama-tech/runlog-api/internal/config"
)

// HealthHandler is responsible for handling health check endpoints in the RunLog API application. It provides methods to check the overall health of the application, including database connectivity and service status. The handler responds with JSON-formatted health information, including service name, version, environment, uptime, and database status.
type HealthHandler struct {
	db        *pgxpool.Pool
	appConfig config.AppConfig
	startedAt time.Time
}

// NewHealthHandler creates a new instance of HealthHandler with the provided database connection pool, application configuration, and application start time. It initializes the handler with the necessary dependencies to perform health checks and respond to incoming HTTP requests.
func NewHealthHandler(db *pgxpool.Pool, appConfig config.AppConfig, startedAt time.Time) *HealthHandler {
	return &HealthHandler{
		db:        db,
		appConfig: appConfig,
		startedAt: startedAt,
	}
}

// Check performs a comprehensive health check of the application, including database connectivity. It returns a JSON response indicating the service status, version, environment, uptime, timestamp, and database status. The check endpoint is typically used by monitoring systems to assess the overall health of the application and its dependencies.
func (h *HealthHandler) Check(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	dbStatus := "UP"
	status := "UP"
	statusCode := http.StatusOK

	if err := h.db.Ping(ctx); err != nil {
		dbStatus = "DOWN"
		status = "DOWN"
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, gin.H{
		"status":      status,
		"service":     h.appConfig.Name,
		"version":     h.appConfig.Version,
		"environment": h.appConfig.Environment,
		"uptime_sec":  int(time.Since(h.startedAt).Seconds()),
		"timestamp":   time.Now().UTC().Format(time.RFC3339),
		"checks": gin.H{
			"database": dbStatus,
		},
	})
}

// Liveness checks if the application is running and responsive. It returns a JSON response indicating the service status, version, environment, uptime, and timestamp. The liveness endpoint is typically used by container orchestration platforms to determine if the application is alive and should be restarted if it becomes unresponsive.
func (h *HealthHandler) Liveness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":      "UP",
		"service":     h.appConfig.Name,
		"version":     h.appConfig.Version,
		"environment": h.appConfig.Environment,
		"uptime_sec":  int(time.Since(h.startedAt).Seconds()),
		"timestamp":   time.Now().UTC().Format(time.RFC3339),
	})
}

// Readiness checks if the application is ready to handle requests. It verifies the database connectivity and returns a JSON response indicating the service status, version, environment, timestamp, and database status. The readiness endpoint is typically used by container orchestration platforms to determine if the application is ready to receive traffic.
func (h *HealthHandler) Readiness(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	dbStatus := "UP"
	status := "UP"
	statusCode := http.StatusOK

	if err := h.db.Ping(ctx); err != nil {
		dbStatus = "DOWN"
		status = "DOWN"
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, gin.H{
		"status":      status,
		"service":     h.appConfig.Name,
		"version":     h.appConfig.Version,
		"environment": h.appConfig.Environment,
		"timestamp":   time.Now().UTC().Format(time.RFC3339),
		"checks": gin.H{
			"database": dbStatus,
		},
	})
}
