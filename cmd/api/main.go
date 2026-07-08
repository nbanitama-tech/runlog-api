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
	"github.com/nbanitama-tech/runlog-api/internal/config"
	"github.com/nbanitama-tech/runlog-api/internal/handler"
	"github.com/nbanitama-tech/runlog-api/internal/repository"
	"github.com/nbanitama-tech/runlog-api/internal/usecase"
	logger "github.com/nbanitama-tech/runlog-api/pkg/logging"
	"github.com/nbanitama-tech/runlog-api/pkg/middleware"
)

func main() {
	ctx := context.Background()

	cfg := config.Load()

	db, err := config.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer db.Close()

	appLog := logger.New()

	userRepo := repository.NewUserRepository(db)
	userUseCase := usecase.NewUserUseCase(userRepo, cfg.JWTSecret, cfg.JWTExpiryHours)
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
	r.Use(middleware.CORSMiddleware(cfg.CORSAllowOrigins))
	r.Use(middleware.LoggerMiddleware(appLog))

	healthHandler := handler.NewHealthHandler(db)
	r.GET("/health", healthHandler.Check)

	v1 := r.Group("/api/v1")
	{
		v1.POST("/users/register", userHandler.Register)
		v1.POST("/users/login", userHandler.Login)

		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			protected.GET("/users/profile", userHandler.Profile)
			protected.POST("/activities", activityHandler.Create)
			protected.GET("/activities", activityHandler.List)
			protected.GET("/activities/:id", activityHandler.Detail)
			protected.PUT("/activities/:id", activityHandler.Update)
			protected.DELETE("/activities/:id", activityHandler.Delete)
		}
	}

	srv := &http.Server{
		Addr:    ":" + cfg.AppPort,
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
