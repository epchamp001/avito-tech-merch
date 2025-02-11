package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"time"
)

type Config struct {
	Env          string            `yaml:"env" env-default:"dev"`
	Application  ApplicationConfig `yaml:"application"`
	PublicServer ServerConfig      `yaml:"public_server"`
	Storage      StorageConfig     `yaml:"storage"`
	Redis        RedisConfig       `yaml:"redis"`
	JWT          JWTConfig         `yaml:"jwt"`
}

type ApplicationConfig struct {
	App string `yaml:"app"`
}

type ServerConfig struct {
	Enable   bool   `yaml:"enable"`
	Endpoint string `yaml:"endpoint"`
	Port     int    `yaml:"port" env:"PORT"`
}

type StorageConfig struct {
	Hosts                 []string      `yaml:"hosts"`
	Port                  int           `yaml:"port"`
	Database              string        `yaml:"database"`
	Username              string        `yaml:"username"`
	Password              string        `env:"DB_PASSWORD"`
	SSLMode               string        `yaml:"ssl_mode"`
	ConnectionAttempts    int           `yaml:"connection_attempts"`
	InitializationTimeout time.Duration `yaml:"initialization_timeout"`
}

type RedisConfig struct {
	Address  string `yaml:"address"`
	Password string `env:"REDIS_PASSWORD"`
}

type JWTConfig struct {
	SecretKey   string `yaml:"secret_key" env:"JWT_SECRET_KEY"`
	TokenExpiry int    `yaml:"token_expiry"`
}

func LoadConfig(configPath, envPath string) (*Config, error) {
	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf("Ошибка загрузки .env файла: %v\n", err)
	}

	cfg := &Config{}

	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		return nil, fmt.Errorf("Ошибка загрузки config.yaml: %v\n", err)
	}

	return cfg, nil
}
