// Package config provides the configuration structures for the RunLog API application. It defines the AppConfig struct for application-level settings, including name, version, environment, and port. The Config struct aggregates various configuration sections, including AppConfig, DatabaseConfig, JWTConfig, and CORSConfig. These configurations are used throughout the application to manage settings related to the application itself, database connections, JWT authentication, and CORS policies.
package config

// AppConfig holds the application-level configuration settings.
type AppConfig struct {
	Name        string
	Version     string
	Environment string
	Port        string
}

// Config aggregates various configuration sections for the RunLog API application.
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
	CORS     CORSConfig
}
