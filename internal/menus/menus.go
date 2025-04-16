package menus

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/skyzeper/telegram-bot/internal/models"
	"time"
)

// MenuGenerator generates Telegram menus
type MenuGenerator struct{}

// NewMenuGenerator creates a new MenuGenerator
func NewMenuGenerator() *MenuGenerator {
	return &MenuGenerator{}
}

// MainMenu generates the main menu
func (m *MenuGenerator) MainMenu(user *models.User) tgbotapi.ReplyKeyboardMarkup {
	buttons := [][]tgbotapi.KeyboardButton{
		{tgbotapi.NewKeyboardButton("🗑️ Заказать услугу")},
		{tgbotapi.NewKeyboardButton("📞 Связаться с оператором")},
		{tgbotapi.NewKeyboardButton("🔗 Приглашайте друзей и зарабатывайте!")},
	}
	if user.Role == "operator" || user.Role == "main_operator" || user.Role == "owner" {
		buttons = append(buttons, []tgbotapi.KeyboardButton{tgbotapi.NewKeyboardButton("📋 Заказы")})
	}
	if user.Role == "main_operator" || user.Role == "owner" {
		buttons = append(buttons, []tgbotapi.KeyboardButton{tgbotapi.NewKeyboardButton("🧑‍💼 Управление штатом")})
	}
	if user.Role == "owner" {
		buttons = append(buttons, []tgbotapi.KeyboardButton{tgbotapi.NewKeyboardButton("📊 Статистика")})
	}
	return tgbotapi.NewReplyKeyboard(buttons...)
}

// CategoryMenu generates the category selection menu
func (m *MenuGenerator) CategoryMenu() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🗑️ Вывоз мусора"),
			tgbotapi.NewKeyboardButton("🔨 Демонтаж"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🏗️ Стройматериалы"),
			tgbotapi.NewKeyboardButton("🔙 Главное меню"),
		),
	)
}

// SubcategoryMenu generates the subcategory selection menu
func (m *MenuGenerator) SubcategoryMenu(category string) tgbotapi.ReplyKeyboardMarkup {
	var buttons []tgbotapi.KeyboardButton
	switch category {
	case "вывоз мусора":
		buttons = []tgbotapi.KeyboardButton{
			tgbotapi.NewKeyboardButton("Строительный мусор"),
			tgbotapi.NewKeyboardButton("Бытовой мусор"),
			tgbotapi.NewKeyboardButton("Мебель"),
		}
	case "демонтаж":
		buttons = []tgbotapi.KeyboardButton{
			tgbotapi.NewKeyboardButton("Стены"),
			tgbotapi.NewKeyboardButton("Полы"),
			tgbotapi.NewKeyboardButton("Потолки"),
		}
	case "стройматериалы":
		buttons = []tgbotapi.KeyboardButton{
			tgbotapi.NewKeyboardButton("Песок"),
			tgbotapi.NewKeyboardButton("Цемент"),
			tgbotapi.NewKeyboardButton("Кирпич"),
		}
	}
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(buttons...),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("🔙 Главное меню")),
	)
}

// DateMenu generates the date selection menu
func (m *MenuGenerator) DateMenu() tgbotapi.InlineKeyboardMarkup {
	var buttons []tgbotapi.InlineKeyboardButton
	now := time.Now()
	for i := 0; i < 7; i++ {
		date := now.AddDate(0, 0, i)
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(
			date.Format("02.01.2006"),
			fmt.Sprintf("date_%s", date.Format("2006-01-02")),
		))
	}
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(buttons[0:4]...),
		tgbotapi.NewInlineKeyboardRow(buttons[4:7]...),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Пред. неделя", "prev_week"),
			tgbotapi.NewInlineKeyboardButtonData("➡️ След. неделя", "next_week"),
		),
	)
}

// TimeMenu generates the time selection menu
func (m *MenuGenerator) TimeMenu() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("09:00"),
			tgbotapi.NewKeyboardButton("12:00"),
			tgbotapi.NewKeyboardButton("15:00"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("18:00"),
			tgbotapi.NewKeyboardButton("🔙 Главное меню"),
		),
	)
}

// PhotoMenu generates the photo upload menu
func (m *MenuGenerator) PhotoMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📸 Добавить фото", "photo_add"),
			tgbotapi.NewInlineKeyboardButtonData("➡️ Пропустить", "photo_skip"),
		),
	)
}

// VideoMenu generates the video upload menu
func (m *MenuGenerator) VideoMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🎥 Добавить видео", "video_add"),
			tgbotapi.NewInlineKeyboardButtonData("➡️ Пропустить", "video_skip"),
		),
	)
}

// PhoneMenu generates the phone input menu
func (m *MenuGenerator) PhoneMenu() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButtonContact("📞 Отправить номер"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🔙 Главное меню"),
		),
	)
}

// SkipMenu generates the skip option menu
func (m *MenuGenerator) SkipMenu() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("➡️ Пропустить"),
			tgbotapi.NewKeyboardButton("🔙 Главное меню"),
		),
	)
}

// PaymentMenu generates the payment method selection menu
func (m *MenuGenerator) PaymentMenu() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("💵 Наличные"),
			tgbotapi.NewKeyboardButton("💳 Карта"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🔙 Главное меню"),
		),
	)
}

// ConfirmMenu generates the order confirmation menu
func (m *MenuGenerator) ConfirmMenu() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("✅ Подтвердить"),
			tgbotapi.NewKeyboardButton("🔙 Главное меню"),
		),
	)
}

// NewOrdersMenu generates the new orders menu
func (m *MenuGenerator) NewOrdersMenu(orders []models.Order) tgbotapi.InlineKeyboardMarkup {
	var buttons []tgbotapi.InlineKeyboardButton
	for _, order := range orders {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("Заказ #%d", order.ID),
			fmt.Sprintf("order_%d", order.ID),
		))
	}
	rows := make([][]tgbotapi.InlineKeyboardButton, 0)
	for i := 0; i < len(buttons); i += 2 {
		end := i + 2
		if end > len(buttons) {
			end = len(buttons)
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(buttons[i:end]...))
	}
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// InProgressOrdersMenu generates the in-progress orders menu
func (m *MenuGenerator) InProgressOrdersMenu(category string, orders []models.Order) tgbotapi.InlineKeyboardMarkup {
	var buttons []tgbotapi.InlineKeyboardButton
	buttons = append(buttons,
		tgbotapi.NewInlineKeyboardButtonData("🗑️ Вывоз мусора", "in_progress_waste_removal"),
		tgbotapi.NewInlineKeyboardButtonData("🔨 Демонтаж", "in_progress_demolition"),
		tgbotapi.NewInlineKeyboardButtonData("🏗️ Стройматериалы", "in_progress_construction_materials"),
		tgbotapi.NewInlineKeyboardButtonData("📋 Все заказы", "in_progress_all"),
	)
	for _, order := range orders {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("Заказ #%d", order.ID),
			fmt.Sprintf("order_%d", order.ID),
		))
	}
	rows := make([][]tgbotapi.InlineKeyboardButton, 0)
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(buttons[0:4]...))
	for i := 4; i < len(buttons); i += 2 {
		end := i + 2
		if end > len(buttons) {
			end = len(buttons)
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(buttons[i:end]...))
	}
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// OrderActionsMenu generates the order actions menu
func (m *MenuGenerator) OrderActionsMenu(orderID int) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Принять", fmt.Sprintf("accept_order_%d", orderID)),
			tgbotapi.NewInlineKeyboardButtonData("❌ Отменить", fmt.Sprintf("cancel_order_%d", orderID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📞 Связаться с клиентом", fmt.Sprintf("contact_client_%d", orderID)),
			tgbotapi.NewInlineKeyboardButtonData("🚫 Заблокировать клиента", fmt.Sprintf("block_client_%d", orderID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🚫 Заблокировать и отменить", fmt.Sprintf("block_cancel_%d", orderID)),
		),
	)
}

// InProgressOrderActionsMenu generates the in-progress order actions menu
func (m *MenuGenerator) InProgressOrderActionsMenu(orderID int) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👷 Назначить исполнителей", fmt.Sprintf("assign_executor_%d", orderID)),
			tgbotapi.NewInlineKeyboardButtonData("✅ Подтвердить выполнение", fmt.Sprintf("confirm_order_%d", orderID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💸 Учесть наличные", fmt.Sprintf("cash_order_%d", orderID)),
			tgbotapi.NewInlineKeyboardButtonData("📞 Связаться с клиентом", fmt.Sprintf("contact_client_%d", orderID)),
		),
	)
}

// AssignExecutorMenu generates the executor assignment menu
func (m *MenuGenerator) AssignExecutorMenu(orderID int) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🛻 Водители", fmt.Sprintf("assign_drivers_%d", orderID)),
			tgbotapi.NewInlineKeyboardButtonData("💪 Грузчики", fmt.Sprintf("assign_loaders_%d", orderID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Подтвердить", fmt.Sprintf("confirm_executors_%d", orderID)),
			tgbotapi.NewInlineKeyboardButtonData("❌ Отменить", fmt.Sprintf("cancel_executors_%d", orderID)),
		),
	)
}

// StaffMenu generates the staff management menu
func (m *MenuGenerator) StaffMenu(chatID int64) tgbotapi.InlineKeyboardMarkup {
	buttons := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("🧑‍💼 Добавить сотрудника", "staff_add"),
		tgbotapi.NewInlineKeyboardButtonData("🗑️ Удалить сотрудника", "staff_delete"),
	}
	if true { // Replace with owner role check
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData("📋 Список сотрудников", "staff_list"))
	}
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(buttons...),
	)
}

// ContactOperatorMenu generates the contact operator menu
func (m *MenuGenerator) ContactOperatorMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📞 Сам позвоню", "contact_call"),
			tgbotapi.NewInlineKeyboardButtonData("📲 Позвоните мне", "contact_request_call"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💬 В чате", "contact_chat"),
		),
	)
}

// ReferralMenu generates the referral program menu
func (m *MenuGenerator) ReferralMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔗 Получить ссылку", "referral_link"),
			tgbotapi.NewInlineKeyboardButtonData("📷 Получить QR-код", "referral_qr"),
		),
	)
}
