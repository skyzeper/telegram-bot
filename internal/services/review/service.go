package review

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/skyzeper/telegram-bot/internal/models"
)

type Service struct {
	repo *Repository
}

func NewService(db *sql.DB) *Service {
	return &Service{
		repo: NewRepository(db),
	}
}

func (s *Service) RequestReview(ctx context.Context, orderID, chatID int64, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(chatID, "Спасибо за заказ! Будем очень благодарны за ваш отзыв! 🙌")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🌟 Оценить", fmt.Sprintf("review_rate_%d", orderID)),
			tgbotapi.NewInlineKeyboardButtonURL("📝 Оставить отзыв на Авито", "https://www.avito.ru/brands/i110181488/all/predlozheniya_uslug?src=search_seller_info&sellerId=c5142d6c5f2dbb8d7bda54b817575f76"),
		),
	)
	bot.Send(msg)
}

func (s *Service) SaveReview(ctx context.Context, orderID, chatID int64, rating int, comment string) error {
	review := &models.Review{
		OrderID: orderID,
		UserID:  chatID,
		Rating:  rating,
		Comment: comment,
	}
	return s.repo.CreateReview(ctx, review)
}
