package review

import (
	"database/sql"
	"fmt"
	"github.com/skyzeper/telegram-bot/internal/db"
	"github.com/skyzeper/telegram-bot/internal/models"
	"github.com/skyzeper/telegram-bot/internal/utils"
)

// PostgresRepository implements the Repository interface for PostgreSQL
type PostgresRepository struct {
	db *db.DB
}

// NewPostgresRepository creates a new PostgreSQL repository
func NewPostgresRepository(db *db.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// CreateReview creates a new review
func (r *PostgresRepository) CreateReview(review *models.Review) error {
	query := `
		INSERT INTO reviews (order_id, user_id, rating, comment, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	err := r.db.Conn().QueryRow(
		query,
		review.OrderID, review.UserID, review.Rating, review.Comment, review.CreatedAt,
	).Scan(&review.ID)
	if err != nil {
		utils.LogError(err)
		return fmt.Errorf("failed to create review: %v", err)
	}
	return nil
}

// GetReview retrieves a review by order ID
func (r *PostgresRepository) GetReview(orderID int) (*models.Review, error) {
	query := `
		SELECT id, order_id, user_id, rating, comment, created_at
		FROM reviews
		WHERE order_id = $1
	`
	review := &models.Review{}
	err := r.db.Conn().QueryRow(query, orderID).Scan(
		&review.ID, &review.OrderID, &review.UserID, &review.Rating, &review.Comment, &review.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("review not found")
	}
	if err != nil {
		utils.LogError(err)
		return nil, fmt.Errorf("failed to get review: %v", err)
	}
	return review, nil
}

// GetReviewsByUser retrieves all reviews by a user
func (r *PostgresRepository) GetReviewsByUser(userID int64) ([]models.Review, error) {
	query := `
		SELECT id, order_id, user_id, rating, comment, created_at
		FROM reviews
		WHERE user_id = $1
	`
	rows, err := r.db.Conn().Query(query, userID)
	if err != nil {
		utils.LogError(err)
		return nil, fmt.Errorf("failed to get reviews by user: %v", err)
	}
	defer rows.Close()

	var reviews []models.Review
	for rows.Next() {
		var review models.Review
		if err := rows.Scan(
			&review.ID, &review.OrderID, &review.UserID, &review.Rating, &review.Comment, &review.CreatedAt,
		); err != nil {
			utils.LogError(err)
			continue
		}
		reviews = append(reviews, review)
	}
	return reviews, nil
}