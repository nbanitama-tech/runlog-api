package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/nbanitama-tech/runlog-api/internal/config"
	"github.com/nbanitama-tech/runlog-api/internal/handler"
)

func main() {
	ctx := context.Background()

	cfg := config.Load()

	db, err := config.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(ctx); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	r := gin.Default()

	healthHandler := handler.NewHealthHandler(db)
	r.GET("/health", healthHandler.Check)

	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
