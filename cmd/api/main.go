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
	logger "github.com/nbanitama-tech/runlog-api/pkg/logging"
	"github.com/nbanitama-tech/runlog-api/pkg/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	db, err := config.NewPostgresPool(ctx, cfg.Database.URL)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer db.Close()

	appLog := logger.New()

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
