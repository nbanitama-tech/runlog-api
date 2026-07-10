package config

// CORSConfig holds the configuration settings for Cross-Origin Resource Sharing (CORS) in the RunLog API application. It defines the allowed origins that can access the API, enabling secure cross-origin requests from specified domains.
type CORSConfig struct {
	AllowOrigins []string
}
