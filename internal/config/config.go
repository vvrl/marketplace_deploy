package config

import (
	"fmt"
	"marketplace/internal/logger"
	"os"

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

	cfg := &Config{
		Server: ServerConfig{
			Port: getEnvOrDefault("PORT", viper.GetString("server.port")),
		},
		Database: DatabaseConfig{
			User:     getEnvOrDefault("DB_USER", viper.GetString("database.user")),
			Password: getEnvOrDefault("DB_PASSWORD", viper.GetString("database.password")),
			Host:     getEnvOrDefault("DB_HOST", viper.GetString("database.host")),
			Port:     getEnvOrDefaultInt("DB_PORT", viper.GetInt("database.port")),
			Dbname:   getEnvOrDefault("DB_NAME", viper.GetString("database.dbname")),
		},
		JWT: JwtConfig{
			Key:      getEnvOrDefault("SECRET_KEY", viper.GetString("jwt.key")),
			Lifetime: viper.GetInt("jwt.lifetime"),
		},
	}

	err := viper.Unmarshal(&cfg)
	if err != nil {
		logger.Logger.Fatalf("failed config parsing into the structure: %v", err)
	}

	return cfg
}

func getEnvOrDefault(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvOrDefaultInt(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		var i int
		fmt.Sscanf(val, "%d", &i)
		return i
	}
	return fallback
}
