package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Env          string             `mapstructure:"env"`
	Application  ApplicationConfig  `mapstructure:"application"`
	PublicServer PublicServerConfig `mapstructure:"public_server"`
	Storage      StorageConfig      `mapstructure:"storage"`
	JWT          JWTConfig          `mapstructure:"jwt"`
	Metrics      MetricsConfig      `mapstructure:"metrics"`
}

func LoadConfig(configPath, envPath string) (*Config, error) {
	err := godotenv.Load(envPath)
	if err != nil {
		fmt.Printf("WARNING: error loading .env file from %s: %v\n", envPath, err)
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)

	err = viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("error reading config file from %s: %v", configPath, err)
	}

	viper.AutomaticEnv()

	viper.BindEnv("storage.postgres.hosts", "DB_HOST")
	viper.BindEnv("storage.postgres.password", "DB_PASSWORD")
	viper.BindEnv("storage.redis.host", "REDIS_HOST")
	viper.BindEnv("storage.redis.password", "REDIS_PASSWORD")
	viper.BindEnv("jwt.secret_key", "JWT_SECRET_KEY")

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %v", err)
	}

	return &config, nil
}
