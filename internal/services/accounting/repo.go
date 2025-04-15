package accounting

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

func (r *Repository) CreateRecord(ctx context.Context, record *models.AccountingRecord) error {
	query := `
		INSERT INTO accounting_records (order_id, user_id, type, amount, description)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`
	return r.db.QueryRowContext(ctx, query,
		record.OrderID, record.UserID, record.Type, record.Amount, record.Description,
	).Scan(&record.ID, &record.CreatedAt)
}
