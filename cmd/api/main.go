package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/nbanitama-tech/runlog-api/internal/config"
	"github.com/nbanitama-tech/runlog-api/internal/handler"
	"github.com/nbanitama-tech/runlog-api/internal/repository"
	"github.com/nbanitama-tech/runlog-api/internal/usecase"
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

	userRepo := repository.NewUserRepository(db)
	userUseCase := usecase.NewUserUseCase(userRepo, cfg.JWTSecret, cfg.JWTExpiryHours)
	userHandler := handler.NewUserHandler(userUseCase)

	activityRepo := repository.NewActivityRepository(db)
	activityUseCase := usecase.NewActivityUseCase(activityRepo)
	activityHandler := handler.NewActivityHandler(activityUseCase)

	if err := db.Ping(ctx); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	r := gin.Default()

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
		}
	}

	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
