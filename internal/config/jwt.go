package config

type JWTConfig struct {
	SecretKey   string `yaml:"secret_key" env:"JWT_SECRET_KEY"`
	TokenExpiry int    `yaml:"token_expiry"`
}
