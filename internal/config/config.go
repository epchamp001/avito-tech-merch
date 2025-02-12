package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
)

type Config struct {
	Env          string            `yaml:"env" env-default:"dev"`
	Application  ApplicationConfig `yaml:"application"`
	PublicServer ServerConfig      `yaml:"public_server"`
	Storage      StorageConfig     `yaml:"storage"`
	JWT          JWTConfig         `yaml:"jwt"`
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
