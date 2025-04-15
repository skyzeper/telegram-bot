package referral

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/skip2/go-qrcode"
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

func (s *Service) GenerateLink(chatID int64) string {
	return fmt.Sprintf("@vseVsimferopole?start=ref_%d", chatID)
}

func (s *Service) GenerateQRCode(chatID int64, bot *tgbotapi.BotAPI) error {
	link := s.GenerateLink(chatID)
	qr, err := qrcode.New(link, qrcode.Medium)
	if err != nil {
		return err
	}
	filename := fmt.Sprintf("qr_%d.png", chatID)
	if err := qr.WriteFile(256, filename); err != nil {
		return err
	}
	photo := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath(filename))
	bot.Send(photo)
	return nil
}

func (s *Service) RegisterReferral(ctx context.Context, inviterID, inviteeID int64) error {
	referral := &models.Referral{
		InviterID: inviterID,
		InviteeID: inviteeID,
	}
	return s.repo.CreateReferral(ctx, referral)
}
