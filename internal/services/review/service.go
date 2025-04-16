package review

import (
	"errors"
	"time"
	"github.com/skyzeper/telegram-bot/internal/models"
)

// Service handles review-related business logic
type Service struct {
	repo Repository
}

// Repository defines the interface for review data access
type Repository interface {
	CreateReview(review *models.Review) error
	GetReview(orderID int) (*models.Review, error)
	GetReviewsByUser(userID int64) ([]models.Review, error)
}

// NewService creates a new review service
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// SubmitReview creates or updates a review for an order
func (s *Service) SubmitReview(orderID int, userID int64, rating int, comment string) error {
	if orderID <= 0 || userID <= 0 {
		return errors.New("invalid order or user ID")
	}
	if rating < 1 || rating > 5 {
		return errors.New("rating must be between 1 and 5")
	}

	review := &models.Review{
		OrderID:   orderID,
		UserID:    userID,
		Rating:    rating,
		Comment:   comment,
		CreatedAt: time.Now(),
	}

	return s.repo.CreateReview(review)
}

// GetReview retrieves a review by order ID
func (s *Service) GetReview(orderID int) (*models.Review, error) {
	if orderID <= 0 {
		return nil, errors.New("invalid order ID")
	}
	return s.repo.GetReview(orderID)
}

// GetReviewsByUser retrieves all reviews by a user
func (s *Service) GetReviewsByUser(userID int64) ([]models.Review, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user ID")
	}
	return s.repo.GetReviewsByUser(userID)
}