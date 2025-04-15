package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/skyzeper/telegram-bot/internal/menus"
	"github.com/skyzeper/telegram-bot/internal/state"
)

func HandleCallback(ctx context.Context, bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery, db *sql.DB) {
	handler := NewHandler(db, state.NewStateManager())
	chatID := callback.Message.Chat.ID
	data := callback.Data

	// Удаляем предыдущее сообщение для чистоты UX
	if callback.Message != nil {
		bot.Send(tgbotapi.NewDeleteMessage(chatID, callback.Message.MessageID))
	}

	currentState := handler.states.Get(chatID)
	switch {
	case data == "action_construction_materials":
		handler.orders.StartOrder(ctx, chatID, "construction_materials")
		bot.Send(tgbotapi.NewMessage(chatID, "Шаг 1/9: Введите имя для заказа стройматериалов"))
	case data == "action_waste_removal":
		handler.orders.StartOrder(ctx, chatID, "waste_removal")
		msgConfig := tgbotapi.NewMessage(chatID, "Шаг 1/9: Выберите подкатегорию")
		msgConfig.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Мусор и хлам", "subcategory_trash"),
				tgbotapi.NewInlineKeyboardButtonData("Старая мебель", "subcategory_furniture"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Металл", "subcategory_metal"),
				tgbotapi.NewInlineKeyboardButtonData("Строительный мусор", "subcategory_construction_waste"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Покрышки", "subcategory_tires"),
				tgbotapi.NewInlineKeyboardButtonData("Пищевые отходы", "subcategory_food_waste"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Другое", "subcategory_other_waste"),
				tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "action_back"),
			),
		)
		bot.Send(msgConfig)
	case data == "action_demolition":
		handler.orders.StartOrder(ctx, chatID, "demolition")
		msgConfig := tgbotapi.NewMessage(chatID, "Шаг 1/7: Выберите подкатегорию")
		msgConfig.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Полы", "subcategory_floors"),
				tgbotapi.NewInlineKeyboardButtonData("Сантехника", "subcategory_plumbing"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Стены", "subcategory_walls"),
				tgbotapi.NewInlineKeyboardButtonData("Двери и окна", "subcategory_doors_windows"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Подготовка к ремонту", "subcategory_repair_prep"),
				tgbotapi.NewInlineKeyboardButtonData("Снос домов", "subcategory_house_demolition"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Другое", "subcategory_other_demolition"),
				tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "action_back"),
			),
		)
		bot.Send(msgConfig)
	case data == "action_contact_operator":
		msgConfig := tgbotapi.NewMessage(chatID, "📞 Как связаться?")
		msgConfig.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Сам позвоню", "chat_call_me"),
				tgbotapi.NewInlineKeyboardButtonData("Позвоните мне", "chat_call_you"),
				tgbotapi.NewInlineKeyboardButtonData("В чате", "chat_text"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "action_back"),
			),
		)
		bot.Send(msgConfig)
	case data == "action_referral":
		link := handler.referral.GenerateLink(chatID)
		msgConfig := tgbotapi.NewMessage(chatID, fmt.Sprintf("🔗 Приглашайте друзей и получите 500 руб. за заказ от 10,000 руб.!\nВаша ссылка: %s", link))
		msgConfig.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("📷 Создать QR-код", "referral_qr"),
				tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "action_back"),
			),
		)
		bot.Send(msgConfig)
	case strings.HasPrefix(data, "subcategory_"):
		if currentState.Module == "create_order" {
			currentState.Data["subcategory"] = strings.TrimPrefix(data, "subcategory_")
			currentState.Step = 2
			currentState.TotalSteps = 9
			if currentState.Data["category"] == "demolition" {
				currentState.TotalSteps = 7
			}
			handler.states.Set(chatID, currentState)
			bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Шаг 2/%d: Введите ваше имя", currentState.TotalSteps)))
		}
	case data == "photo_add":
		if currentState.Module == "create_order" {
			bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Шаг 3/%d: Загрузите ещё фото или видео", currentState.TotalSteps)))
		}
	case data == "photo_confirm":
		if currentState.Module == "create_order" {
			currentState.Step = 5
			currentState.TotalSteps = currentState.TotalSteps
			handler.states.Set(chatID, currentState)
			msgConfig := tgbotapi.NewMessage(chatID, fmt.Sprintf("Шаг 5/%d: Выберите дату", currentState.TotalSteps))
			msgConfig.ReplyMarkup = handler.orders.DateKeyboard()
			bot.Send(msgConfig)
		}
	case data == "photo_cancel":
		if currentState.Module == "create_order" {
			handler.states.Clear(chatID)
			msgConfig := tgbotapi.NewMessage(chatID, "Создание заказа отменено")
			msgConfig.ReplyMarkup = menus.UserMenu()
			bot.Send(msgConfig)
		}
	case data == "phone_contact" || data == "chat_contact":
		if currentState.Module == "create_order" || currentState.Module == "chat" {
			msgConfig := tgbotapi.NewMessage(chatID, "Нажмите кнопку для отправки номера")
			msgConfig.ReplyMarkup = tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButtonContact("📞 Отправить номер"),
				),
			)
			bot.Send(msgConfig)
		}
	case data == "chat_call_me":
		bot.Send(tgbotapi.NewMessage(chatID, "📞 Позвоните: +7(978)-959-70-77"))
	case data == "chat_call_you":
		handler.states.Set(chatID, state.State{
			Module: "chat",
			Data:   map[string]interface{}{},
		})
		msgConfig := tgbotapi.NewMessage(chatID, "Введите номер или нажмите 📞")
		msgConfig.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("📞 Отправить номер", "chat_contact"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "action_back"),
			),
		)
		bot.Send(msgConfig)
	case data == "chat_text":
		handler.states.Set(chatID, state.State{
			Module: "chat",
			Data:   map[string]interface{}{},
		})
		bot.Send(tgbotapi.NewMessage(chatID, "Напишите ваше сообщение оператору"))
	case data == "action_back":
		user, _ := handler.users.GetUser(ctx, chatID)
		var menu tgbotapi.InlineKeyboardMarkup
		switch user.Role {
		case "user":
			menu = menus.UserMenu()
		case "operator":
			menu = menus.OperatorMenu()
		case "main_operator":
			menu = menus.MainOperatorMenu()
		case "driver":
			menu = menus.DriverMenu()
		case "loader":
			menu = menus.LoaderMenu()
		case "owner":
			menu = menus.OwnerMenu()
		}
		handler.states.Clear(chatID)
		msgConfig := tgbotapi.NewMessage(chatID, "Выберите действие")
		msgConfig.ReplyMarkup = menu
		bot.Send(msgConfig)
	default:
		bot.Send(tgbotapi.NewMessage(chatID, "Неизвестное действие"))
	}

	// Подтверждаем callback
	bot.Send(tgbotapi.NewCallback(callback.ID, ""))
}
