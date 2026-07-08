package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort        string
	DatabaseURL    string
	JWTSecret      string
	JWTExpiryHours string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		AppPort:        getEnv("APP_PORT", "8080"),
		DatabaseURL:    getEnv("DATABASE_URL", ""),
		JWTSecret:      getEnv("JWT_SECRET", "dev_secret"),
		JWTExpiryHours: getEnv("JWT_EXPIRY_HOURS", "24"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
