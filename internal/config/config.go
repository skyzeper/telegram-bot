package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken    string
	DBHost           string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBName           string
	Env              string
	OwnerChatID      int64
	AccountingChatID int64
	GroupChatID      int64
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	cfg := &Config{
		TelegramToken:    os.Getenv("TELEGRAM_TOKEN"),
		DBHost:           os.Getenv("DB_HOST"),
		DBPort:           os.Getenv("DB_PORT"),
		DBUser:           os.Getenv("DB_USER"),
		DBPassword:       os.Getenv("DB_PASSWORD"),
		DBName:           os.Getenv("DB_NAME"),
		Env:              os.Getenv("ENV"),
		OwnerChatID:      parseInt64(os.Getenv("OWNER_CHAT_ID")),
		AccountingChatID: parseInt64(os.Getenv("ACCOUNTING_CHAT_ID")),
		GroupChatID:      parseInt64(os.Getenv("GROUP_CHAT_ID")),
	}

	if cfg.TelegramToken == "" || cfg.DBHost == "" || cfg.DBName == "" {
		return nil, fmt.Errorf("missing required environment variables")
	}

	return cfg, nil
}

func parseInt64(s string) int64 {
	n, _ := strconv.ParseInt(s, 10, 64)
	return n
}
