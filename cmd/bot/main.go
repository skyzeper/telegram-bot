package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/skyzeper/telegram-bot/internal/config"
	"github.com/skyzeper/telegram-bot/internal/db"
	"github.com/skyzeper/telegram-bot/internal/handlers"
	"github.com/skyzeper/telegram-bot/internal/security"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	dbConn, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := dbConn.Close(); err != nil {
			log.Printf("Failed to close DB: %v", err)
		}
	}()

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatalf("Failed to initialize bot: invalid token: %v", err)
	}

	bot.Debug = cfg.Env == "dev"

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := bot.GetUpdatesChan(updateConfig)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for update := range updates {
				if update.Message != nil {
					if security.CheckAccess(ctx, update.Message.From.ID, dbConn) {
						handlers.HandleMessage(ctx, bot, update.Message, dbConn)
					} else {
						bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "🚫 Доступ запрещён."))
					}
				} else if update.CallbackQuery != nil {
					if security.CheckAccess(ctx, update.CallbackQuery.From.ID, dbConn) {
						handlers.HandleCallback(ctx, bot, update.CallbackQuery, dbConn)
					}
				}
			}
		}()
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	cancel()
	wg.Wait()
	log.Println("Bot stopped")
}
