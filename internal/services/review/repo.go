package review

import (
	"context"
	"database/sql"

	"bot/internal/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateReview(ctx context.Context, review *models.Review) error {
	query := `
		INSERT INTO reviews (order_id, user_id, rating, comment)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`
	return r.db.QueryRowContext(ctx, query,
		review.OrderID, review.UserID, review.Rating, review.Comment,
	).Scan(&review.ID, &review.CreatedAt)
}
