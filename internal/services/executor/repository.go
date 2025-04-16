package executor

import (
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

// AssignExecutor assigns an executor to an order
func (r *PostgresRepository) AssignExecutor(orderID int, userID int64, role string) error {
	query := `
		INSERT INTO executors (order_id, user_id, role, confirmed, notified, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	var id int
	err := r.db.Conn().QueryRow(
		query,
		orderID, userID, role, false, false, time.Now(),
	).Scan(&id)
	if err != nil {
		utils.LogError(err)
		return fmt.Errorf("failed to assign executor: %v", err)
	}
	return nil
}

// RemoveExecutor removes an executor from an order
func (r *PostgresRepository) RemoveExecutor(orderID int, userID int64) error {
	query := `DELETE FROM executors WHERE order_id = $1 AND user_id = $2`
	_, err := r.db.Conn().Exec(query, orderID, userID)
	if err != nil {
		utils.LogError(err)
		return fmt.Errorf("failed to remove executor: %v", err)
	}
	return nil
}

// GetExecutors retrieves all executors for an order
func (r *PostgresRepository) GetExecutors(orderID int) ([]models.Executor, error) {
	query := `
		SELECT id, order_id, user_id, role, confirmed, notified, created_at
		FROM executors
		WHERE order_id = $1
	`
	rows, err := r.db.Conn().Query(query, orderID)
	if err != nil {
		utils.LogError(err)
		return nil, fmt.Errorf("failed to get executors: %v", err)
	}
	defer rows.Close()

	var executors []models.Executor
	for rows.Next() {
		var exec models.Executor
		if err := rows.Scan(
			&exec.ID, &exec.OrderID, &exec.UserID, &exec.Role, &exec.Confirmed,
			&exec.Notified, &exec.CreatedAt,
		); err != nil {
			utils.LogError(err)
			continue
		}
		executors = append(executors, exec)
	}
	return executors, nil
}

// ConfirmExecutor confirms an executor's completion
func (r *PostgresRepository) ConfirmExecutor(orderID int, userID int64) error {
	query := `
		UPDATE executors
		SET confirmed = TRUE
		WHERE order_id = $1 AND user_id = $2
	`
	_, err := r.db.Conn().Exec(query, orderID, userID)
	if err != nil {
		utils.LogError(err)
		return fmt.Errorf("failed to confirm executor: %v", err)
	}
	return nil
}