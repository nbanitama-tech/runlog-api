package config

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

func Load() (*Config, error) {
	_ = godotenv.Load()

	jwtExpiryHours, err := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))
	if err != nil {
		return nil, errors.New("JWT_EXPIRY_HOURS must be a valid number")
	}

	cfg := &Config{
		App: AppConfig{
			Name:        getEnv("APP_NAME", "runlog-api"),
			Version:     getEnv("APP_VERSION", "1.0.0"),
			Environment: getEnv("APP_ENV", "local"),
			Port:        getEnv("APP_PORT", "8080"),
		},
		Database: DatabaseConfig{
			URL: getEnv("DATABASE_URL", ""),
		},
		JWT: JWTConfig{
			Secret:      getEnv("JWT_SECRET", ""),
			ExpiryHours: jwtExpiryHours,
		},
		CORS: CORSConfig{
			AllowOrigins: splitCSV(getEnv("CORS_ALLOW_ORIGINS", "http://localhost:3000,http://localhost:5173")),
		},
	}

	if err := validate(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func validate(cfg *Config) error {
	if cfg.Database.URL == "" {
		return errors.New("DATABASE_URL is required")
	}

	if cfg.JWT.Secret == "" {
		return errors.New("JWT_SECRET is required")
	}

	return nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func splitCSV(value string) []string {
	items := strings.Split(value, ",")
	result := make([]string, 0, len(items))

	for _, item := range items {
		item = strings.TrimSpace(item)
		if item != "" {
			result = append(result, item)
		}
	}

	return result
}
