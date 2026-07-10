package config

// JWTConfig holds the configuration settings for JSON Web Token (JWT) authentication in the RunLog API application. It defines the secret key used for signing and verifying JWT tokens, as well as the token expiration duration in hours. These settings are essential for implementing secure authentication and authorization mechanisms in the application.
type JWTConfig struct {
	Secret      string
	ExpiryHours int
}
