package config

type JWTConfig struct {
	Secret      string
	ExpiryHours int
}
