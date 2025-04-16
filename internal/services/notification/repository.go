package notification

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

// CreateNotification saves a notification
func (r *PostgresRepository) CreateNotification(notification *models.Notification) error {
	query := `
		INSERT INTO notifications (user_id, type, message, sent_at, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	var sentAt sql.NullTime
	if !notification.SentAt.IsZero() {
		sentAt.Valid = true
		sentAt.Time = notification.SentAt
	}
	err := r.db.Conn().QueryRow(
		query,
		notification.UserID, notification.Type, notification.Message,
		sentAt, notification.CreatedAt,
	).Scan(&notification.ID)
	if err != nil {
		utils.LogError(err)
		return fmt.Errorf("failed to create notification: %v", err)
	}
	return nil
}

// GetPendingNotifications retrieves unsent notifications
func (r *PostgresRepository) GetPendingNotifications() ([]models.Notification, error) {
	query := `
		SELECT id, user_id, type, message, sent_at, created_at
		FROM notifications
		WHERE sent_at IS NULL
	`
	rows, err := r.db.Conn().Query(query)
	if err != nil {
		utils.LogError(err)
		return nil, fmt.Errorf("failed to get pending notifications: %v", err)
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var notification models.Notification
		var sentAt sql.NullTime
		if err := rows.Scan(
			&notification.ID, &notification.UserID, &notification.Type,
			&notification.Message, &sentAt, &notification.CreatedAt,
		); err != nil {
			utils.LogError(err)
			continue
		}
		if sentAt.Valid {
			notification.SentAt = sentAt.Time
		}
		notifications = append(notifications, notification)
	}
	return notifications, nil
}

// MarkNotificationSent marks a notification as sent
func (r *PostgresRepository) MarkNotificationSent(notificationID int) error {
	query := `
		UPDATE notifications
		SET sent_at = $1
		WHERE id = $2
	`
	_, err := r.db.Conn().Exec(query, time.Now(), notificationID)
	if err != nil {
		utils.LogError(err)
		return fmt.Errorf("failed to mark notification sent: %v", err)
	}
	return nil
}