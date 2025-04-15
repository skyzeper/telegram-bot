package executor

import (
	"context"
	"database/sql"
	"fmt"

	"bot/internal/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateExecutor(ctx context.Context, executor *models.Executor) error {
	query := `
		INSERT INTO executors (order_id, user_id, role)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`
	return r.db.QueryRowContext(ctx, query,
		executor.OrderID, executor.UserID, executor.Role,
	).Scan(&executor.ID, &executor.CreatedAt)
}

func (r *Repository) GetExecutor(ctx context.Context, orderID, userID int64) (*models.Executor, error) {
	executor := &models.Executor{}
	query := `
		SELECT id, order_id, user_id, role, confirmed, notified, created_at
		FROM executors WHERE order_id = $1 AND user_id = $2`
	err := r.db.QueryRowContext(ctx, query, orderID, userID).Scan(
		&executor.ID, &executor.OrderID, &executor.UserID, &executor.Role,
		&executor.Confirmed, &executor.Notified, &executor.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("executor not found")
	}
	return executor, err
}

func (r *Repository) UpdateExecutor(ctx context.Context, executor *models.Executor) error {
	query := `
		UPDATE executors
		SET confirmed = $1, notified = $2
		WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query,
		executor.Confirmed, executor.Notified, executor.ID,
	)
	return err
}
