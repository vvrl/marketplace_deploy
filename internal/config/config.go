package config

import (
	"marketplace/internal/logger"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Logger   LoggerConfig
	JWT      JwtConfig
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

type DatabaseConfig struct {
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	User        string `yaml:"user"`
	Password    string `yaml:"password"`
	Dbname      string `yaml:"dbname"`
	MaxAttempts int    `yaml:"maxAttempts"`
}

type LoggerConfig struct {
	Level string `yaml:"level"`
}

type JwtConfig struct {
	Key      string `yaml:"key"`
	Lifetime int    `yaml:"lifetime"`
}

func NewConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	if err := viper.ReadInConfig(); err != nil {
		logger.Logger.Fatalf("config reading error: %v", err)
	}

	var cfg Config

	err := viper.Unmarshal(&cfg)
	if err != nil {
		logger.Logger.Fatalf("failed config parsing into the structure: %v", err)
	}

	return &cfg
}
