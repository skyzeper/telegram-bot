package chat

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

// CreateMessage saves a chat message
func (r *PostgresRepository) CreateMessage(message *models.Message) error {
	query := `
		INSERT INTO messages (user_id, operator_id, message, is_from_user, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	var operatorID sql.NullInt64
	if message.OperatorID != 0 {
		operatorID.Valid = true
		operatorID.Int64 = message.OperatorID
	}
	err := r.db.Conn().QueryRow(
		query,
		message.UserID, operatorID, message.Message, message.IsFromUser, message.CreatedAt,
	).Scan(&message.ID)
	if err != nil {
		utils.LogError(err)
		return fmt.Errorf("failed to create message: %v", err)
	}
	return nil
}

// GetMessagesByUser retrieves all messages for a user
func (r *PostgresRepository) GetMessagesByUser(userID int64) ([]models.Message, error) {
	query := `
		SELECT id, user_id, operator_id, message, is_from_user, created_at
		FROM messages
		WHERE user_id = $1
	`
	rows, err := r.db.Conn().Query(query, userID)
	if err != nil {
		utils.LogError(err)
		return nil, fmt.Errorf("failed to get messages by user: %v", err)
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		var operatorID sql.NullInt64
		if err := rows.Scan(
			&msg.ID, &msg.UserID, &operatorID, &msg.Message, &msg.IsFromUser, &msg.CreatedAt,
		); err != nil {
			utils.LogError(err)
			continue
		}
		if operatorID.Valid {
			msg.OperatorID = operatorID.Int64
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

// GetActiveOperator retrieves an active operator ID
func (r *PostgresRepository) GetActiveOperator() (int64, error) {
	query := `
		SELECT chat_id
		FROM users
		WHERE role IN ('operator', 'main_operator') AND is_blocked = FALSE
		ORDER BY RANDOM()
		LIMIT 1
	`
	var operatorID int64
	err := r.db.Conn().QueryRow(query).Scan(&operatorID)
	if err == sql.ErrNoRows {
		return 0, fmt.Errorf("no active operators found")
	}
	if err != nil {
		utils.LogError(err)
		return 0, fmt.Errorf("failed to get active operator: %v", err)
	}
	return operatorID, nil
}