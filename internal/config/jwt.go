package config

type JWTConfig struct {
	SecretKey   string `mapstructure:"secret_key"`
	TokenExpiry int    `mapstructure:"token_expiry"`
}
