package user

import (
	"errors"
	"time"
	"github.com/skyzeper/telegram-bot/internal/models"
)

// Service handles user-related business logic
type Service struct {
	repo Repository
}

// Repository defines the interface for user data access
type Repository interface {
	CreateUser(user *models.User) error
	GetUser(chatID int64) (*models.User, error)
	GetUserByRole(chatID int64, role string) (*models.User, error)
	ListUsersByRole(role string) ([]models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(chatID int64) error
}

// NewService creates a new user service
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// CreateUser creates a new user
func (s *Service) CreateUser(user *models.User) error {
	if user.ChatID == 0 || user.Role == "" {
		return errors.New("missing required user fields")
	}

	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}
	if user.UpdatedAt.IsZero() {
		user.UpdatedAt = time.Now()
	}

	return s.repo.CreateUser(user)
}

// GetUser retrieves a user by ChatID
func (s *Service) GetUser(chatID int64) (*models.User, error) {
	if chatID <= 0 {
		return nil, errors.New("invalid chat ID")
	}
	return s.repo.GetUser(chatID)
}

// GetUserByRole retrieves a user by ChatID and role
func (s *Service) GetUserByRole(chatID int64, role string) (*models.User, error) {
	if chatID <= 0 || role == "" {
		return nil, errors.New("invalid chat ID or role")
	}
	return s.repo.GetUserByRole(chatID, role)
}

// ListUsersByRole retrieves users by role
func (s *Service) ListUsersByRole(role string) ([]models.User, error) {
	if role == "" {
		return nil, errors.New("role cannot be empty")
	}
	return s.repo.ListUsersByRole(role)
}

// UpdateUser updates an existing user
func (s *Service) UpdateUser(user *models.User) error {
	if user.ChatID <= 0 {
		return errors.New("invalid chat ID")
	}
	user.UpdatedAt = time.Now()
	return s.repo.UpdateUser(user)
}

// DeleteUser deletes a user
func (s *Service) DeleteUser(chatID int64) error {
	if chatID <= 0 {
		return errors.New("invalid chat ID")
	}
	return s.repo.DeleteUser(chatID)
}