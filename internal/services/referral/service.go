package referral

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
	"github.com/skyzeper/telegram-bot/internal/models"
	"github.com/skip2/go-qrcode"
)

// Service handles referral-related business logic
type Service struct {
	repo Repository
}

// Repository defines the interface for referral data access
type Repository interface {
	CreateReferral(referral *models.Referral) error
	GetReferralByInvitee(inviteeID int64) (*models.Referral, error)
	GetReferralsByInviter(inviterID int64) ([]models.Referral, error)
	UpdateReferral(referral *models.Referral) error
}

// NewService creates a new referral service
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// CreateReferral creates a new referral
func (s *Service) CreateReferral(inviterID, inviteeID, orderID int64) error {
	if inviterID <= 0 || inviteeID <= 0 {
		return errors.New("invalid inviter or invitee ID")
	}

	referral := &models.Referral{
		InviterID:      inviterID,
		InviteeID:      inviteeID,
		OrderID:        int(orderID),
		PayoutRequested: false,
		CreatedAt:      time.Now(),
	}

	return s.repo.CreateReferral(referral)
}

// GenerateQRCode generates a QR code for a referral link
func (s *Service) GenerateQRCode(link string, userID int64) (string, error) {
	if link == "" {
		return "", errors.New("empty referral link")
	}

	// Create a temporary file path
	filename := fmt.Sprintf("referral_qr_%d.png", userID)
	filepath := filepath.Join(os.TempDir(), filename)

	// Generate QR code
	err := qrcode.WriteFile(link, qrcode.Medium, 256, filepath)
	if err != nil {
		return "", fmt.Errorf("failed to generate QR code: %v", err)
	}

	return filepath, nil
}

// RequestPayout marks a referral payout as requested
func (s *Service) RequestPayout(inviterID, inviteeID int64) error {
	if inviterID <= 0 || inviteeID <= 0 {
		return errors.New("invalid inviter or invitee ID")
	}

	referral, err := s.repo.GetReferralByInvitee(inviteeID)
	if err != nil {
		return fmt.Errorf("failed to find referral: %v", err)
	}
	if referral.InviterID != inviterID {
		return errors.New("inviter does not match referral")
	}
	if referral.PayoutRequested {
		return errors.New("payout already requested")
	}

	referral.PayoutRequested = true
	return s.repo.UpdateReferral(referral)
}

// GetReferralsByInviter retrieves all referrals for an inviter
func (s *Service) GetReferralsByInviter(inviterID int64) ([]models.Referral, error) {
	if inviterID <= 0 {
		return nil, errors.New("invalid inviter ID")
	}
	return s.repo.GetReferralsByInviter(inviterID)
}