package config

import (
	"fmt"
	"time"
)

type StorageConfig struct {
	Hosts                 []string      `yaml:"hosts" env:"DB_HOST"`
	Port                  int           `yaml:"port"`
	Database              string        `yaml:"database"`
	Username              string        `yaml:"username"`
	Password              string        `env:"DB_PASSWORD"`
	SSLMode               string        `yaml:"ssl_mode"`
	ConnectionAttempts    int           `yaml:"connection_attempts"`
	InitializationTimeout time.Duration `yaml:"initialization_timeout"`
}

func (s *StorageConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		s.Hosts[0], s.Port, s.Username, s.Password, s.Database, s.SSLMode,
	)
}
