package executor

import (
	"errors"
	"github.com/skyzeper/telegram-bot/internal/models"
)

// Service handles executor-related business logic
type Service struct {
	repo Repository
}

// Repository defines the interface for executor data access
type Repository interface {
	AssignExecutor(orderID int, userID int64, role string) error
	RemoveExecutor(orderID int, userID int64) error
	GetExecutors(orderID int) ([]models.Executor, error)
	ConfirmExecutor(orderID int, userID int64) error
}

// NewService creates a new executor service
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// AssignExecutor assigns an executor to an order
func (s *Service) AssignExecutor(orderID int, userID int64, role string) error {
	if orderID <= 0 || userID <= 0 {
		return errors.New("invalid order or user ID")
	}
	if role != "driver" && role != "loader" {
		return errors.New("invalid role")
	}

	return s.repo.AssignExecutor(orderID, userID, role)
}

// RemoveExecutor removes an executor from an order
func (s *Service) RemoveExecutor(orderID int, userID int64) error {
	if orderID <= 0 || userID <= 0 {
		return errors.New("invalid order or user ID")
	}
	return s.repo.RemoveExecutor(orderID, userID)
}

// GetExecutors retrieves all executors for an order
func (s *Service) GetExecutors(orderID int) ([]models.Executor, error) {
	if orderID <= 0 {
		return nil, errors.New("invalid order ID")
	}
	return s.repo.GetExecutors(orderID)
}

// ConfirmExecutor confirms an executor's completion
func (s *Service) ConfirmExecutor(orderID int, userID int64) error {
	if orderID <= 0 || userID <= 0 {
		return errors.New("invalid order or user ID")
	}
	return s.repo.ConfirmExecutor(orderID, userID)
}