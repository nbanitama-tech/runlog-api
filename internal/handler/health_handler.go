package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nbanitama-tech/runlog-api/internal/config"
)

type HealthHandler struct {
	db        *pgxpool.Pool
	appConfig config.AppConfig
	startedAt time.Time
}

func NewHealthHandler(db *pgxpool.Pool, appConfig config.AppConfig, startedAt time.Time) *HealthHandler {
	return &HealthHandler{
		db:        db,
		appConfig: appConfig,
		startedAt: startedAt,
	}
}

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
