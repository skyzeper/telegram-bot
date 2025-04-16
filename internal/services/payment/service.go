package payment

import (
	"errors"
	"time"
	"github.com/skyzeper/telegram-bot/internal/models"
)

// Service handles payment-related business logic
type Service struct {
	repo Repository
}

// Repository defines the interface for payment data access
type Repository interface {
	CreatePayment(payment *models.Payment) error
	GetPendingPayments(orderID int) ([]models.Payment, error)
	ConfirmPayment(orderID int, driverID int64) error
	GetPayment(orderID int, driverID int64) (*models.Payment, error)
}

// NewService creates a new payment service
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// CreatePayment creates a new payment
func (s *Service) CreatePayment(payment *models.Payment) error {
	if payment.OrderID <= 0 || payment.UserID <= 0 || payment.Amount <= 0 || payment.Method == "" {
		return errors.New("missing required payment fields")
	}

	if payment.CreatedAt.IsZero() {
		payment.CreatedAt = time.Now()
	}

	return s.repo.CreatePayment(payment)
}

// GetPendingPayments retrieves pending payments for an order
func (s *Service) GetPendingPayments(orderID int) ([]models.Payment, error) {
	if orderID <= 0 {
		return nil, errors.New("invalid order ID")
	}
	return s.repo.GetPendingPayments(orderID)
}

// ConfirmPayment confirms a payment
func (s *Service) ConfirmPayment(orderID int, driverID int64) error {
	if orderID <= 0 || driverID <= 0 {
		return errors.New("invalid order or driver ID")
	}
	return s.repo.ConfirmPayment(orderID, driverID)
}

// GetPayment retrieves a specific payment
func (s *Service) GetPayment(orderID int, driverID int64) (*models.Payment, error) {
	if orderID <= 0 || driverID <= 0 {
		return nil, errors.New("invalid order or driver ID")
	}
	return s.repo.GetPayment(orderID, driverID)
}