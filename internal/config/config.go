package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort          string
	DatabaseURL      string
	JWTSecret        string
	JWTExpiryHours   string
	CORSAllowOrigins string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		AppPort:          getEnv("APP_PORT", "8080"),
		DatabaseURL:      getEnv("DATABASE_URL", ""),
		JWTSecret:        getEnv("JWT_SECRET", "dev_secret"),
		JWTExpiryHours:   getEnv("JWT_EXPIRY_HOURS", "24"),
		CORSAllowOrigins: getEnv("CORS_ALLOW_ORIGINS", "http://localhost:3000,http://localhost:5173"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
