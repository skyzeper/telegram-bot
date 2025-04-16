package main

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/skyzeper/telegram-bot/internal/db"
	"github.com/skyzeper/telegram-bot/internal/handlers"
	"github.com/skyzeper/telegram-bot/internal/handlers/callbacks"
	"github.com/skyzeper/telegram-bot/internal/menus"
	"github.com/skyzeper/telegram-bot/internal/security"
	"github.com/skyzeper/telegram-bot/internal/services/accounting"
	"github.com/skyzeper/telegram-bot/internal/services/chat"
	"github.com/skyzeper/telegram-bot/internal/services/executor"
	"github.com/skyzeper/telegram-bot/internal/services/notification"
	"github.com/skyzeper/telegram-bot/internal/services/order"
	"github.com/skyzeper/telegram-bot/internal/services/payment"
	"github.com/skyzeper/telegram-bot/internal/services/referral"
	"github.com/skyzeper/telegram-bot/internal/services/review"
	"github.com/skyzeper/telegram-bot/internal/services/stats"
	"github.com/skyzeper/telegram-bot/internal/services/user"
	"github.com/skyzeper/telegram-bot/internal/state"
	"github.com/skyzeper/telegram-bot/internal/utils"
)

func main() {
	// Load configuration
	cfg, err := utils.LoadConfig()
	if err != nil {
		utils.LogError(fmt.Errorf("failed to load config: %v", err))
		return
	}

	// Initialize database
	dbCfg := db.Config{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
	}
	dbConn, err := db.NewDB(dbCfg)
	if err != nil {
		utils.LogError(fmt.Errorf("failed to initialize database: %v", err))
		return
	}
	defer dbConn.Close()

	// Initialize Telegram bot
	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		utils.LogError(fmt.Errorf("failed to initialize bot: %v", err))
		return
	}
	bot.Debug = cfg.BotDebug

	// Initialize state
	stateManager := state.NewState()

	// Initialize services
	userService := user.NewService(user.NewPostgresRepository(dbConn))
	orderService := order.NewService(order.NewPostgresRepository(dbConn))
	executorService := executor.NewService(executor.NewPostgresRepository(dbConn))
	paymentService := payment.NewService(payment.NewPostgresRepository(dbConn))
	reviewService := review.NewService(review.NewPostgresRepository(dbConn))
	referralService := referral.NewService(referral.NewPostgresRepository(dbConn))
	chatService := chat.NewService(bot, chat.NewPostgresRepository(dbConn))
	notificationService := notification.NewService(bot, notification.NewPostgresRepository(dbConn))
	accountingService := accounting.NewService(accounting.NewPostgresRepository(dbConn))
	statsService := stats.NewService(stats.NewPostgresRepository(dbConn))

	// Initialize security
	securityChecker := security.NewSecurityChecker(userService)

	// Initialize menus
	menuGenerator := menus.NewMenuGenerator()

	// Initialize callback handlers
	ordersHandler := callbacks.NewOrdersHandler(
		bot, securityChecker, menuGenerator, userService, orderService,
		chatService, executorService, paymentService, reviewService, stateManager,
	)
	staffHandler := callbacks.NewStaffHandler(bot, securityChecker, menuGenerator, userService, stateManager)
	contactHandler := callbacks.NewContactHandler(bot, securityChecker, menuGenerator, chatService, stateManager)
	referralsHandler := callbacks.NewReferralsHandler(bot, securityChecker, menuGenerator, referralService, userService)
	reviewsHandler := callbacks.NewReviewsHandler(bot, securityChecker, menuGenerator, reviewService, stateManager)
	statsHandler := callbacks.NewStatsHandler(bot, securityChecker, menuGenerator, statsService)
	callbackHandler := callbacks.NewCallbackHandler(
		bot, securityChecker, menuGenerator, userService, stateManager,
		ordersHandler, staffHandler, contactHandler, referralsHandler, reviewsHandler, statsHandler,
	)

	// Initialize main handler
	mainHandler := handlers.NewHandler(
		bot, securityChecker, menuGenerator, userService, orderService,
		chatService, stateManager, callbackHandler, notificationService,
	)

	// Set up Telegram updates
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	// Process updates
	for update := range updates {
		mainHandler.HandleUpdate(&update)
	}
}