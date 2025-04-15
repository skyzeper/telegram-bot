package menus

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func UserMenu() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🛠️ Стройматериалы"),
			tgbotapi.NewKeyboardButton("🗑️ Вывоз мусора"),
			tgbotapi.NewKeyboardButton("🔨 Демонтаж"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📞 Связаться с оператором"),
			tgbotapi.NewKeyboardButton("🔗 Пригласить друга"),
		),
	)
}
