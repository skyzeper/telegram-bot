package stats

import (
	"database/sql"
	"fmt"
	"time"
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

// GetOrderStats retrieves orders within a time range
func (r *PostgresRepository) GetOrderStats(start, end time.Time) ([]models.Order, error) {
	query := `
		SELECT id, user_id, category, subcategory, photos, video, date, time, phone, 
		       address, description, status, reason, cost, payment_method, payment_confirmed, 
		       created_at, updated_at, confirmed
		FROM orders
	`
	args := []interface{}{}
	if !start.IsZero() && !end.IsZero() {
		query += ` WHERE created_at >= $1 AND created_at < $2`
		args = append(args, start, end)
	} else if !start.IsZero() {
		query += ` WHERE created_at >= $1`
		args = append(args, start)
	} else if !end.IsZero() {
		query += ` WHERE created_at < $1`
		args = append(args, end)
	}

	rows, err := r.db.Conn().Query(query, args...)
	if err != nil {
		utils.LogError(err)
		return nil, fmt.Errorf("failed to get order stats: %v", err)
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		var date, time sql.NullTime
		if err := rows.Scan(
			&order.ID, &order.UserID, &order.Category, &order.Subcategory, &order.Photos,
			&order.Video, &date, &time, &order.Phone, &order.Address,
			&order.Description, &order.Status, &order.Reason, &order.Cost, &order.PaymentMethod,
			&order.PaymentConfirmed, &order.CreatedAt, &order.UpdatedAt, &order.Confirmed,
		); err != nil {
			utils.LogError(err)
			continue
		}
		if date.Valid {
			order.Date = date.Time
		}
		if time.Valid {
			order.Time = time.Time
		}
		orders = append(orders, order)
	}
	return orders, nil
}

// GetAccountingStats retrieves accounting records within a time range
func (r *PostgresRepository) GetAccountingStats(start, end time.Time) ([]models.AccountingRecord, error) {
	query := `
		SELECT id, order_id, user_id, type, amount, description, created_at
		FROM accounting_records
	`
	args := []interface{}{}
	if !start.IsZero() && !end.IsZero() {
		query += ` WHERE created_at >= $1 AND created_at < $2`
		args = append(args, start, end)
	} else if !start.IsZero() {
		query += ` WHERE created_at >= $1`
		args = append(args, start)
	} else if !end.IsZero() {
		query += ` WHERE created_at < $1`
		args = append(args, end)
	}

	rows, err := r.db.Conn().Query(query, args...)
	if err != nil {
		utils.LogError(err)
		return nil, fmt.Errorf("failed to get accounting stats: %v", err)
	}
	defer rows.Close()

	var records []models.AccountingRecord
	for rows.Next() {
		var record models.AccountingRecord
		var orderID sql.NullInt64
		if err := rows.Scan(
			&record.ID, &orderID, &record.UserID, &record.Type,
			&record.Amount, &record.Description, &record.CreatedAt,
		); err != nil {
			utils.LogError(err)
			continue
		}
		if orderID.Valid {
			record.OrderID = int(orderID.Int64)
		}
		records = append(records, record)
	}
	return records, nil
}