package security

import (
	"github.com/skyzeper/telegram-bot/internal/db"
	"github.com/skyzeper/telegram-bot/internal/models"
	"github.com/skyzeper/telegram-bot/internal/services/user"
)

// SecurityChecker handles user access control
type SecurityChecker struct {
	userService *user.Service
}

// NewSecurityChecker creates a new SecurityChecker
func NewSecurityChecker(dbConn *db.DB) *SecurityChecker {
	userRepo := user.NewPostgresRepository(dbConn)
	userService := user.NewService(userRepo)
	return &SecurityChecker{userService: userService}
}

// HasAccess checks if a user has access to a specific module
func (s *SecurityChecker) HasAccess(chatID int64, module string) (bool, error) {
	user, err := s.userService.GetUser(chatID)
	if err != nil {
		return false, err
	}

	if user.IsBlocked {
		return false, nil
	}

	switch module {
	case "orders":
		return user.Role == "client" || user.Role == "operator" || user.Role == "main_operator" || user.Role == "owner", nil
	case "staff":
		return user.Role == "main_operator" || user.Role == "owner", nil
	case "contact":
		return user.Role == "client" || user.Role == "operator" || user.Role == "main_operator" || user.Role == "owner", nil
	case "referrals":
		return user.Role == "client" || user.Role == "operator" || user.Role == "main_operator" || user.Role == "owner", nil
	case "reviews":
		return user.Role == "client" || user.Role == "operator" || user.Role == "main_operator" || user.Role == "owner", nil
	case "stats":
		return user.Role == "owner", nil
	default:
		return false, nil
	}
}

// IsBlocked checks if a user is blocked
func (s *SecurityChecker) IsBlocked(chatID int64) (bool, error) {
	user, err := s.userService.GetUser(chatID)
	if err != nil {
		return false, err
	}
	return user.IsBlocked, nil
}

// GetUserRole retrieves the user's role
func (s *SecurityChecker) GetUserRole(chatID int64) (string, error) {
	user, err := s.userService.GetUser(chatID)
	if err != nil {
		return "", err
	}
	return user.Role, nil
}

// CreateUser creates a new user if they don't exist
func (s *SecurityChecker) CreateUser(chatID int64, role, firstName, lastName, nickname, phone string) error {
	user, err := s.userService.GetUser(chatID)
	if err == nil && user != nil {
		return nil // User already exists
	}

	newUser := &models.User{
		ChatID:    chatID,
		Role:      role,
		FirstName: firstName,
		LastName:  lastName,
		Nickname:  nickname,
		Phone:     phone,
		IsBlocked: false,
	}
	return s.userService.CreateUser(newUser)
}