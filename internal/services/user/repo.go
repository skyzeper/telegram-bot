package user

import (
	"bot/internal/models"
	"context"
	"database/sql"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (chat_id, role, first_name, last_name, nickname, phone)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`
	return r.db.QueryRowContext(ctx, query,
		user.ChatID, user.Role, user.FirstName, user.LastName, user.Nickname, user.Phone,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *Repository) GetUserByChatID(ctx context.Context, chatID int64) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, chat_id, role, first_name, last_name, nickname, phone, is_blocked, created_at, updated_at
		FROM users WHERE chat_id = $1`
	err := r.db.QueryRowContext(ctx, query, chatID).Scan(
		&user.ID, &user.ChatID, &user.Role, &user.FirstName, &user.LastName,
		&user.Nickname, &user.Phone, &user.IsBlocked, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	return user, err
}

func (r *Repository) UpdateUser(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET role = $1, first_name = $2, last_name = $3, nickname = $4, phone = $5, is_blocked = $6, updated_at = CURRENT_TIMESTAMP
		WHERE chat_id = $7`
	_, err := r.db.ExecContext(ctx, query,
		user.Role, user.FirstName, user.LastName, user.Nickname, user.Phone, user.IsBlocked, user.ChatID,
	)
	return err
}

func (r *Repository) LogBlock(ctx context.Context, chatID int64, reason string) error {
	query := `INSERT INTO chat_messages (user_id, message, is_from_user) VALUES ((SELECT id FROM users WHERE chat_id = $1), $2, FALSE)`
	_, err := r.db.ExecContext(ctx, query, chatID, fmt.Sprintf("Вы заблокированы. Причина: %s", reason))
	return err
}
