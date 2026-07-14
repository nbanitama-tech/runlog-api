// Package main is the entry point of the RunLog API application. It initializes the application, sets up the database connection, configures routes and middleware, and starts the HTTP server. The application provides endpoints for user registration, login, profile management, and activity tracking. It also includes health check endpoints and Swagger documentation for API reference.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/nbanitama-tech/runlog-api/docs"
	"github.com/nbanitama-tech/runlog-api/internal/config"
	"github.com/nbanitama-tech/runlog-api/internal/handler"
	"github.com/nbanitama-tech/runlog-api/internal/repository"
	"github.com/nbanitama-tech/runlog-api/internal/usecase"
	logger "github.com/nbanitama-tech/runlog-api/pkg/observability/logging"
	"github.com/nbanitama-tech/runlog-api/pkg/transport/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	appmetrics "github.com/nbanitama-tech/runlog-api/pkg/observability/metrics"
)

// @title RunLog API
// @version 1.0
// @description Production-style running activity tracker API.
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	ctx := context.Background()
	startedAt := time.Now()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	appLog := logger.New()

	metricsRegistry := prometheus.NewRegistry()

	metricsRegistry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	httpMetrics := appmetrics.NewHTTPMetrics(metricsRegistry)
	databaseMetrics := appmetrics.NewDatabaseMetrics(metricsRegistry)
	databaseTracer := appmetrics.NewPGXTracer(databaseMetrics)

	db, err := config.NewPostgresPool(ctx, cfg.Database.URL, databaseTracer)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer db.Close()

	metricsRegistry.MustRegister(appmetrics.NewPGXPoolCollector(db))

	userRepo := repository.NewUserRepository(db)
	userUseCase := usecase.NewUserUseCase(userRepo, cfg.JWT.Secret, cfg.JWT.ExpiryHours)
	userHandler := handler.NewUserHandler(userUseCase)

	activityRepo := repository.NewActivityRepository(db)
	activityUseCase := usecase.NewActivityUseCase(activityRepo)
	activityHandler := handler.NewActivityHandler(activityUseCase)

	if err := db.Ping(ctx); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(middleware.RequestIDMiddleware())
	r.Use(middleware.CORSMiddleware(cfg.CORS.AllowOrigins))
	r.Use(middleware.LoggerMiddleware(appLog))
	r.Use(middleware.MetricsMiddleware(httpMetrics))

	healthHandler := handler.NewHealthHandler(db, cfg.App, startedAt)
	r.GET("/health", healthHandler.Check)
	r.GET("/health/live", healthHandler.Liveness)
	r.GET("/health/ready", healthHandler.Readiness)

	v1 := r.Group("/api/v1")
	{
		v1.POST("/users/register", userHandler.Register)
		v1.POST("/users/login", userHandler.Login)

		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		{
			protected.GET("/users/profile", userHandler.Profile)
			protected.POST("/activities", activityHandler.Create)
			protected.GET("/activities", activityHandler.List)
			protected.GET("/activities/:id", activityHandler.Detail)
			protected.PUT("/activities/:id", activityHandler.Update)
			protected.DELETE("/activities/:id", activityHandler.Delete)
		}
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	metricsHandler := promhttp.HandlerFor(
		metricsRegistry,
		promhttp.HandlerOpts{},
	)

	r.GET("/metrics", gin.WrapH(metricsHandler))

	srv := &http.Server{
		Addr:    ":" + cfg.App.Port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to run server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	appLog.Info("server exited gracefully")
}
