package accounting

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

// CreateRecord creates a new accounting record
func (r *PostgresRepository) CreateRecord(record *models.AccountingRecord) error {
	query := `
		INSERT INTO accounting_records (order_id, user_id, type, amount, description, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	var orderID sql.NullInt64
	if record.OrderID != 0 {
		orderID.Valid = true
		orderID.Int64 = int64(record.OrderID)
	}
	err := r.db.Conn().QueryRow(
		query,
		orderID, record.UserID, record.Type, record.Amount, record.Description, record.CreatedAt,
	).Scan(&record.ID)
	if err != nil {
		utils.LogError(err)
		return fmt.Errorf("failed to create accounting record: %v", err)
	}
	return nil
}

// GetRecordsByOrder retrieves accounting records for an order
func (r *PostgresRepository) GetRecordsByOrder(orderID int) ([]models.AccountingRecord, error) {
	query := `
		SELECT id, order_id, user_id, type, amount, description, created_at
		FROM accounting_records
		WHERE order_id = $1
	`
	rows, err := r.db.Conn().Query(query, orderID)
	if err != nil {
		utils.LogError(err)
		return nil, fmt.Errorf("failed to get records by order: %v", err)
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

// GetRecordsByUser retrieves accounting records for a user
func (r *PostgresRepository) GetRecordsByUser(userID int64) ([]models.AccountingRecord, error) {
	query := `
		SELECT id, order_id, user_id, type, amount, description, created_at
		FROM accounting_records
		WHERE user_id = $1
	`
	rows, err := r.db.Conn().Query(query, userID)
	if err != nil {
		utils.LogError(err)
		return nil, fmt.Errorf("failed to get records by user: %v", err)
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

// GetRecordsByType retrieves accounting records by type
func (r *PostgresRepository) GetRecordsByType(recordType string) ([]models.AccountingRecord, error) {
	query := `
		SELECT id, order_id, user_id, type, amount, description, created_at
		FROM accounting_records
		WHERE type = $1
	`
	rows, err := r.db.Conn().Query(query, recordType)
	if err != nil {
		utils.LogError(err)
		return nil, fmt.Errorf("failed to get records by type: %v", err)
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

// UpdateRecord updates an existing accounting record
func (r *PostgresRepository) UpdateRecord(record *models.AccountingRecord) error {
	query := `
		UPDATE accounting_records
		SET order_id = $1, user_id = $2, type = $3, amount = $4, description = $5, created_at = $6
		WHERE id = $7
	`
	var orderID sql.NullInt64
	if record.OrderID != 0 {
		orderID.Valid = true
		orderID.Int64 = int64(record.OrderID)
	}
	_, err := r.db.Conn().Exec(
		query,
		orderID, record.UserID, record.Type, record.Amount, record.Description,
		record.CreatedAt, record.ID,
	)
	if err != nil {
		utils.LogError(err)
		return fmt.Errorf("failed to update accounting record: %v", err)
	}
	return nil
}