package callbacks

import (
	"fmt"
	"strings"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/skyzeper/telegram-bot/internal/models"
	"github.com/skyzeper/telegram-bot/internal/menus"
	"github.com/skyzeper/telegram-bot/internal/services/stats"
	"github.com/skyzeper/telegram-bot/internal/security"
	"github.com/skyzeper/telegram-bot/internal/utils"
)

// StatsHandler handles statistics-related callbacks
type StatsHandler struct {
	bot          *tgbotapi.BotAPI
	security     security.SecurityChecker
	menus        *menus.MenuGenerator
	statsService *stats.Service
}

// NewStatsHandler creates a new StatsHandler
func NewStatsHandler(
	bot *tgbotapi.BotAPI,
	security security.SecurityChecker,
	menus *menus.MenuGenerator,
	statsService *stats.Service,
) *StatsHandler {
	return &StatsHandler{
		bot:          bot,
		security:     security,
		menus:        menus,
		statsService: statsService,
	}
}

// HandleStatsCallback processes statistics-related callback queries
func (h *StatsHandler) HandleStatsCallback(callback *tgbotapi.CallbackQuery, user *models.User, data string) {
	if !h.security.HasRole(callback.Message.Chat.ID, "owner") {
		h.sendUnauthorized(callback.Message.Chat.ID)
		return
	}

	if data == "stats_date" {
		h.handleMonthSelection(callback)
	} else if strings.HasPrefix(data, "stats_month_") {
		h.handleWeekSelection(callback, data)
	} else if strings.HasPrefix(data, "stats_week_") {
		h.handleWeekStats(callback, data)
	} else {
		h.handleStats(callback, user, data)
	}
}

// handleStats shows statistics
func (h *StatsHandler) handleStats(callback *tgbotapi.CallbackQuery, user *models.User, data string) {
	var stats models.Stats
	var err error
	period := strings.TrimPrefix(data, "stats_")
	switch period {
	case "day":
		stats, err = h.statsService.GetStatsForDay()
	case "week":
		stats, err = h.statsService.GetStatsForWeek()
	case "month":
		stats, err = h.statsService.GetStatsForMonth()
	case "year":
		stats, err = h.statsService.GetStatsForYear()
	case "all":
		stats, err = h.statsService.GetStatsForAllTime()
	default:
		h.sendError(callback.Message.Chat.ID, "Неверный период.")
		return
	}

	if err != nil {
		h.sendError(callback.Message.Chat.ID, "Ошибка получения статистики.")
		return
	}

	reply := tgbotapi.NewEditMessageText(
		callback.Message.Chat.ID,
		callback.Message.MessageID,
		fmt.Sprintf(
			"📊 Статистика за %s:\n"+
				"- Всего заказов: %d\n"+
				"- Вывоз мусора: %d\n"+
				"- Демонтаж: %d\n"+
				"- Стройматериалы: %d\n"+
				"- Сумма: %.2f руб.\n"+
				"- Долги водителей: %.2f руб.\n\n📈 Держите руку на пульсе бизнеса!",
			period,
			stats.TotalOrders,
			stats.WasteRemovalOrders,
			stats.DemolitionOrders,
			stats.ConstructionOrders,
			stats.TotalAmount,
			stats.DriverDebts,
		),
	)
	reply.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "back_to_stats"),
		),
	)
	if _, err := h.bot.Send(reply); err != nil {
		utils.LogError(err)
	}
}

// handleMonthSelection shows month selection for stats
func (h *StatsHandler) handleMonthSelection(callback *tgbotapi.CallbackQuery) {
	reply := tgbotapi.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, "📅 Выберите месяц:")
	reply.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Январь", "stats_month_01"),
			tgbotapi.NewInlineKeyboardButtonData("Февраль", "stats_month_02"),
			tgbotapi.NewInlineKeyboardButtonData("Март", "stats_month_03"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Апрель", "stats_month_04"),
			tgbotapi.NewInlineKeyboardButtonData("Май", "stats_month_05"),
			tgbotapi.NewInlineKeyboardButtonData("Июнь", "stats_month_06"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Июль", "stats_month_07"),
			tgbotapi.NewInlineKeyboardButtonData("Август", "stats_month_08"),
			tgbotapi.NewInlineKeyboardButtonData("Сентябрь", "stats_month_09"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Октябрь", "stats_month_10"),
			tgbotapi.NewInlineKeyboardButtonData("Ноябрь", "stats_month_11"),
			tgbotapi.NewInlineKeyboardButtonData("Декабрь", "stats_month_12"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "back_to_stats"),
		),
	)
	if _, err := h.bot.Send(reply); err != nil {
		utils.LogError(err)
	}
}

// handleWeekSelection shows week selection for stats
func (h *StatsHandler) handleWeekSelection(callback *tgbotapi.CallbackQuery, data string) {
	month := strings.TrimPrefix(data, "stats_month_")
	reply := tgbotapi.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, fmt.Sprintf("📅 Выберите неделю для месяца %s:", month))
	reply.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("1 неделя", fmt.Sprintf("stats_week_%s_1", month)),
			tgbotapi.NewInlineKeyboardButtonData("2 неделя", fmt.Sprintf("stats_week_%s_2", month)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("3 неделя", fmt.Sprintf("stats_week_%s_3", month)),
			tgbotapi.NewInlineKeyboardButtonData("4 неделя", fmt.Sprintf("stats_week_%s_4", month)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "stats_date"),
		),
	)
	if _, err := h.bot.Send(reply); err != nil {
		utils.LogError(err)
	}
}

// handleWeekStats shows stats for selected week
func (h *StatsHandler) handleWeekStats(callback *tgbotapi.CallbackQuery, data string) {
	parts := strings.Split(strings.TrimPrefix(data, "stats_week_"), "_")
	if len(parts) != 2 {
		h.sendError(callback.Message.Chat.ID, "Неверный формат недели.")
		return
	}

	month := parts[0]
	// Placeholder: Fetch stats for specific week
	// In real implementation, query statsService with month and week
	reply := tgbotapi.NewEditMessageText(
		callback.Message.Chat.ID,
		callback.Message.MessageID,
		fmt.Sprintf("📊 Статистика за неделю %s/%s:\n(данные временно недоступны)", month, parts[1]),
	)
	reply.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", fmt.Sprintf("stats_month_%s", month)),
		),
	)
	if _, err := h.bot.Send(reply); err != nil {
		utils.LogError(err)
	}
}

// sendMainMenu sends the main menu
func (h *StatsHandler) sendMainMenu(chatID int64, user *models.User) {
	reply := tgbotapi.NewMessage(chatID, "Выберите действие:")
	reply.ReplyMarkup = h.menus.MainMenu(user)
	if _, err := h.bot.Send(reply); err != nil {
		utils.LogError(err)
	}
}

// sendError sends an error message
func (h *StatsHandler) sendError(chatID int64, text string) {
	reply := tgbotapi.NewMessage(chatID, text)
	reply.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🔙 Главное меню"),
		),
	)
	if _, err := h.bot.Send(reply); err != nil {
		utils.LogError(err)
	}
}

// sendUnauthorized sends unauthorized access message
func (h *StatsHandler) sendUnauthorized(chatID int64) {
	reply := tgbotapi.NewMessage(chatID, "🚫 Доступ запрещён.")
	reply.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🔙 Главное меню"),
		),
	)
	if _, err := h.bot.Send(reply); err != nil {
		utils.LogError(err)
	}
}