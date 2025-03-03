package config

type PublicServerConfig struct {
	Enable          bool   `mapstructure:"enabled"`
	Endpoint        string `mapstructure:"endpoint"`
	Port            int    `mapstructure:"port"`
	ShutdownTimeout int    `mapstructure:"shutdown_timeout"`
}
