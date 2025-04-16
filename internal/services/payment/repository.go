package payment

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

// CreatePayment creates a new payment
func (r *PostgresRepository) CreatePayment(payment *models.Payment) error {
	query := `
		INSERT INTO payments (order_id, user_id, amount, method, driver_id, confirmed, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	var driverID sql.NullInt64
	if payment.DriverID != 0 {
		driverID.Valid = true
		driverID.Int64 = payment.DriverID
	}
	err := r.db.Conn().QueryRow(
		query,
		payment.OrderID, payment.UserID, payment.Amount, payment.Method,
		driverID, payment.Confirmed, payment.CreatedAt,
	).Scan(&payment.ID)
	if err != nil {
		utils.LogError(err)
		return fmt.Errorf("failed to create payment: %v", err)
	}
	return nil
}

// GetPendingPayments retrieves pending payments for an order
func (r *PostgresRepository) GetPendingPayments(orderID int64) ([]models.Payment, error) {
	query := `
		SELECT id, order_id, user_id, amount, method, driver_id, confirmed, created_at
		FROM payments
		WHERE order_id = $1 AND confirmed = FALSE
	`
	rows, err := r.db.Conn().Query(query, orderID)
	if err != nil {
		utils.LogError(err)
		return nil, fmt.Errorf("failed to get pending payments: %v", err)
	}
	defer rows.Close()

	var payments []models.Payment
	for rows.Next() {
		var payment models.Payment
		var driverID sql.NullInt64
		if err := rows.Scan(
			&payment.ID, &payment.OrderID, &payment.UserID, &payment.Amount, &payment.Method,
			&driverID, &payment.Confirmed, &payment.CreatedAt,
		); err != nil {
			utils.LogError(err)
			continue
		}
		if driverID.Valid {
			payment.DriverID = driverID.Int64
		}
		payments = append(payments, payment)
	}
	return payments, nil
}

// ConfirmPayment confirms a payment
func (r *PostgresRepository) ConfirmPayment(orderID int64, driverID int64) error {
	query := `
		UPDATE payments
		SET confirmed = TRUE
		WHERE order_id = $1 AND driver_id = $2
	`
	_, err := r.db.Conn().Exec(query, orderID, driverID)
	if err != nil {
		utils.LogError(err)
		return fmt.Errorf("failed to confirm payment: %v", err)
	}
	return nil
}

// GetPayment retrieves a specific payment
func (r *PostgresRepository) GetPayment(orderID int64, driverID int64) (*models.Payment, error) {
	query := `
		SELECT id, order_id, user_id, amount, method, driver_id, confirmed, created_at
		FROM payments
		WHERE order_id = $1 AND driver_id = $2
	`
	payment := &models.Payment{}
	var driverIDVal sql.NullInt64
	err := r.db.Conn().QueryRow(query, orderID, driverID).Scan(
		&payment.ID, &payment.OrderID, &payment.UserID, &payment.Amount, &payment.Method,
		&driverIDVal, &payment.Confirmed, &payment.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("payment not found")
	}
	if err != nil {
		utils.LogError(err)
		return nil, fmt.Errorf("failed to get payment: %v", err)
	}
	if driverIDVal.Valid {
		payment.DriverID = driverIDVal.Int64
	}
	return payment, nil
}