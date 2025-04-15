package chat

import (
	"context"
	"database/sql"

	"github.com/skyzeper/telegram-bot/internal/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) SaveMessage(ctx context.Context, msg *models.Message) error {
	query := `
		INSERT INTO chat_messages (user_id, operator_id, message, is_from_user)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`
	return r.db.QueryRowContext(ctx, query,
		msg.UserID, msg.OperatorID, msg.Message, msg.IsFromUser,
	).Scan(&msg.ID, &msg.CreatedAt)
}
