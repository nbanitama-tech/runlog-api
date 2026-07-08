package config

type AppConfig struct {
	Name        string
	Version     string
	Environment string
	Port        string
}

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
	CORS     CORSConfig
}
