package order

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/skyzeper/telegram-bot/internal/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateOrder(ctx context.Context, order *models.Order) error {
	query := `
		INSERT INTO orders (user_id, category, subcategory, photos, video, date, time, phone, address, description, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at`
	return r.db.QueryRowContext(ctx, query,
		order.UserID, order.Category, order.Subcategory, order.Photos, order.Video,
		order.Date, order.Time, order.Phone, order.Address, order.Description, order.Status,
	).Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)
}

func (r *Repository) GetOrder(ctx context.Context, id int64) (*models.Order, error) {
	order := &models.Order{}
	query := `
		SELECT id, user_id, category, subcategory, photos, video, date, time, phone, address,
			description, status, reason, cost, payment_method, payment_confirmed, created_at, updated_at
		FROM orders WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&order.ID, &order.UserID, &order.Category, &order.Subcategory, &order.Photos, &order.Video,
		&order.Date, &order.Time, &order.Phone, &order.Address, &order.Description, &order.Status,
		&order.Reason, &order.Cost, &order.PaymentMethod, &order.PaymentConfirmed, &order.CreatedAt, &order.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("order not found")
	}
	return order, err
}
