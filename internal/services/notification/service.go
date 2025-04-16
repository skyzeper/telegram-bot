package notification

import (
	"fmt"
	"time"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/skyzeper/telegram-bot/internal/models"
)

// Service handles notification-related business logic
type Service struct {
	bot  *tgbotapi.BotAPI
	repo Repository
}

// Repository defines the interface for notification data access
type Repository interface {
	CreateNotification(notification *models.Notification) error
	GetPendingNotifications() ([]models.Notification, error)
	MarkNotificationSent(notificationID int) error
}

// NewService creates a new notification service
func NewService(bot *tgbotapi.BotAPI, repo Repository) *Service {
	return &Service{
		bot:  bot,
		repo: repo,
	}
}

// SendOrderNotification sends a notification about an order event
func (s *Service) SendOrderNotification(userID int64, order *models.Order, event string) error {
	if userID <= 0 || order == nil {
		return fmt.Errorf("invalid user ID or order")
	}

	var message string
	switch event {
	case "created":
		message = fmt.Sprintf(
			"> **Заказ #%d создан!** 🚛\n"+
				"> Категория: %s (%s)\n"+
				"> Адрес: %s\n"+
				"> Мы свяжемся для подтверждения стоимости! 😊",
			order.ID, order.Category, order.Subcategory, order.Address,
		)
	case "confirmed":
		message = fmt.Sprintf(
			"> **Заказ #%d подтверждён!** ✅\n"+
				"> Категория: %s (%s)\n"+
				"> Стоимость: %.2f руб.\n"+
				"> Спасибо за выбор нас! 🙌",
			order.ID, order.Category, order.Subcategory, order.Cost,
		)
	case "assigned":
		message = fmt.Sprintf(
			"> **Заказ #%d назначен!** 👷\n"+
				"> Категория: %s (%s)\n"+
				"> Адрес: %s\n"+
				"> Дата: %s\n"+
				"> Подтвердите выполнение после завершения.",
			order.ID, order.Category, order.Subcategory, order.Address, order.Date.Format("2 January 2006"),
		)
	default:
		return fmt.Errorf("unknown order event: %s", event)
	}

	msg := tgbotapi.NewMessage(userID, message)
	msg.ParseMode = "Markdown"
	if _, err := s.bot.Send(msg); err != nil {
		return fmt.Errorf("failed to send order notification: %v", err)
	}

	notification := &models.Notification{
		UserID:    userID,
		Type:      "order_" + event,
		Message:   message,
		SentAt:    time.Now(),
		CreatedAt: time.Now(),
	}
	return s.repo.CreateNotification(notification)
}

// SendReferralNotification sends a notification about a referral event
func (s *Service) SendReferralNotification(userID, inviteeID int64, event string) error {
	if userID <= 0 || inviteeID <= 0 {
		return fmt.Errorf("invalid user or invitee ID")
	}

	var message string
	switch event {
	case "joined":
		message = fmt.Sprintf(
			"> 🎉 Ваш друг (ID: %d) присоединился по вашей ссылке!\n"+
				"> Получите 500 рублей за его заказ от 10,000 рублей!",
			inviteeID,
		)
	case "payout_requested":
		message = fmt.Sprintf(
			"> 💸 Запрос на выплату за реферала (ID: %d) отправлен!\n"+
				"> Мы свяжемся с вами для перевода 500 рублей.",
			inviteeID,
		)
	default:
		return fmt.Errorf("unknown referral event: %s", event)
	}

	msg := tgbotapi.NewMessage(userID, message)
	msg.ParseMode = "Markdown"
	if _, err := s.bot.Send(msg); err != nil {
		return fmt.Errorf("failed to send referral notification: %v", err)
	}

	notification := &models.Notification{
		UserID:    userID,
		Type:      "referral_" + event,
		Message:   message,
		SentAt:    time.Now(),
		CreatedAt: time.Now(),
	}
	return s.repo.CreateNotification(notification)
}

// SendReviewNotification sends a notification about a review event
func (s *Service) SendReviewNotification(userID int64, review *models.Review) error {
	if userID <= 0 || review == nil {
		return fmt.Errorf("invalid user ID or review")
	}

	message := fmt.Sprintf(
		"> 🌟 Спасибо за отзыв на заказ #%d!\n"+
			"> Оценка: %d/5\n"+
			"> Ваш голос помогает нам стать лучше! 🙌",
		review.OrderID, review.Rating,
	)
	if review.Comment != "" {
		message += fmt.Sprintf("> Комментарий: %s", review.Comment)
	}

	msg := tgbotapi.NewMessage(userID, message)
	msg.ParseMode = "Markdown"
	if _, err := s.bot.Send(msg); err != nil {
		return fmt.Errorf("failed to send review notification: %v", err)
	}

	notification := &models.Notification{
		UserID:    userID,
		Type:      "review_submitted",
		Message:   message,
		SentAt:    time.Now(),
		CreatedAt: time.Now(),
	}
	return s.repo.CreateNotification(notification)
}

// SendOperatorNotification sends a notification to operators
func (s *Service) SendOperatorNotification(operatorID int64, message string) error {
	if operatorID <= 0 || message == "" {
		return fmt.Errorf("invalid operator ID or message")
	}

	msg := tgbotapi.NewMessage(operatorID, message)
	msg.ParseMode = "Markdown"
	if _, err := s.bot.Send(msg); err != nil {
		return fmt.Errorf("failed to send operator notification: %v", err)
	}

	notification := &models.Notification{
		UserID:    operatorID,
		Type:      "operator_alert",
		Message:   message,
		SentAt:    time.Now(),
		CreatedAt: time.Now(),
	}
	return s.repo.CreateNotification(notification)
}

// ProcessPendingNotifications sends pending notifications
func (s *Service) ProcessPendingNotifications() error {
	notifications, err := s.repo.GetPendingNotifications()
	if err != nil {
		return fmt.Errorf("failed to get pending notifications: %v", err)
	}

	for _, notification := range notifications {
		msg := tgbotapi.NewMessage(notification.UserID, notification.Message)
		msg.ParseMode = "Markdown"
		if _, err := s.bot.Send(msg); err != nil {
			continue // Log error but continue processing
		}
		if err := s.repo.MarkNotificationSent(notification.ID); err != nil {
			continue // Log error but continue processing
		}
	}
	return nil
}