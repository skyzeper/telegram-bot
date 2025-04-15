package main

import (
	"bot/internal/state"
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"bot/internal/config"
	"bot/internal/db"
	"bot/internal/handlers"
	"bot/internal/security"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	defer dbConn.Close()

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatalf("Failed to initialize bot: %v", err)
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
						handler := handlers.NewHandler(dbConn, state.NewStateManager())
						handler.HandleCallback(ctx, bot, update.CallbackQuery, dbConn)
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
