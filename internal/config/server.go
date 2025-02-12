package config

type ServerConfig struct {
	Enable   bool   `yaml:"enable"`
	Endpoint string `yaml:"endpoint"`
	Port     int    `yaml:"port" env:"PORT"`
}
