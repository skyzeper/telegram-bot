package payment

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

func (r *Repository) CreatePayment(ctx context.Context, payment *models.Payment) error {
	query := `
		INSERT INTO payments (order_id, user_id, amount, method, driver_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`
	return r.db.QueryRowContext(ctx, query,
		payment.OrderID, payment.UserID, payment.Amount, payment.Method, payment.DriverID,
	).Scan(&payment.ID, &payment.CreatedAt)
}

func (r *Repository) GetPayment(ctx context.Context, id int64) (*models.Payment, error) {
	payment := &models.Payment{}
	query := `
		SELECT id, order_id, user_id, amount, method, driver_id, confirmed, created_at
		FROM payments WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&payment.ID, &payment.OrderID, &payment.UserID, &payment.Amount, &payment.Method,
		&payment.DriverID, &payment.Confirmed, &payment.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("payment not found")
	}
	return payment, err
}

func (r *Repository) UpdatePayment(ctx context.Context, payment *models.Payment) error {
	query := `
		UPDATE payments
		SET confirmed = $1
		WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query,
		payment.Confirmed, payment.ID,
	)
	return err
}
