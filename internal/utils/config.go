package utils

import (
	"fmt"
	"os"
)

// Config holds application configuration
type Config struct {
	BotToken   string
	BotDebug   bool
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	cfg := &Config{
		BotToken:   os.Getenv("BOT_TOKEN"),
		BotDebug:   os.Getenv("BOT_DEBUG") == "true",
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
	}

	if cfg.BotToken == "" {
		return nil, fmt.Errorf("BOT_TOKEN is required")
	}
	if cfg.DBHost == "" {
		cfg.DBHost = "localhost"
	}
	if cfg.DBPort == "" {
		cfg.DBPort = "5432"
	}
	if cfg.DBUser == "" {
		cfg.DBUser = "postgres"
	}
	if cfg.DBPassword == "" {
		cfg.DBPassword = "Ho5049707"
	}
	if cfg.DBName == "" {
		cfg.DBName = "telegram_bot"
	}

	return cfg, nil
}