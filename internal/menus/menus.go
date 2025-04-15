package menus

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func UserMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🛠️ Стройматериалы", "action_construction_materials"),
			tgbotapi.NewInlineKeyboardButtonData("🗑️ Вывоз мусора", "action_waste_removal"),
			tgbotapi.NewInlineKeyboardButtonData("🔨 Демонтаж", "action_demolition"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📞 Связаться с оператором", "action_contact_operator"),
			tgbotapi.NewInlineKeyboardButtonData("🔗 Пригласить друга", "action_referral"),
		),
	)
}

func OperatorMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📋 Заказы", "action_orders"),
			tgbotapi.NewInlineKeyboardButtonData("🚫 Блокировка", "action_block_user"),
		),
	)
}

func MainOperatorMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📋 Заказы", "action_orders"),
			tgbotapi.NewInlineKeyboardButtonData("🚫 Блокировка", "action_block_user"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🧑‍💼 Штат", "action_staff"),
			tgbotapi.NewInlineKeyboardButtonData("🔗 Управление рефералами", "action_manage_referrals"),
		),
	)
}

func OwnerMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📋 Заказы", "action_orders"),
			tgbotapi.NewInlineKeyboardButtonData("🚫 Блокировка", "action_block_user"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🧑‍💼 Штат", "action_staff"),
			tgbotapi.NewInlineKeyboardButtonData("🔗 Управление рефералами", "action_manage_referrals"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📊 Статистика", "action_stats"),
		),
	)
}

func DriverMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📋 Мои заказы", "action_my_orders"),
			tgbotapi.NewInlineKeyboardButtonData("⛽ Расходы", "action_expenses"),
		),
	)
}

func LoaderMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📋 Мои заказы", "action_my_orders"),
		),
	)
}
