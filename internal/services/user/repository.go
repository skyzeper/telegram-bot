package user

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

// CreateUser creates a new user
func (r *PostgresRepository) CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (chat_id, role, first_name, last_name, nickname, phone, is_blocked, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`
	err := r.db.Conn().QueryRow(
		query,
		user.ChatID, user.Role, user.FirstName, user.LastName, user.Nickname,
		user.Phone, user.IsBlocked, user.CreatedAt, user.UpdatedAt,
	).Scan(&user.ID)
	if err != nil {
		utils.LogError(err)
		return fmt.Errorf("failed to create user: %v", err)
	}
	return nil
}

// GetUser retrieves a user by ChatID
func (r *PostgresRepository) GetUser(chatID int64) (*models.User, error) {
	query := `
		SELECT id, chat_id, role, first_name, last_name, nickname, phone, is_blocked, created_at, updated_at
		FROM users
		WHERE chat_id = $1
	`
	user := &models.User{}
	err := r.db.Conn().QueryRow(query, chatID).Scan(
		&user.ID, &user.ChatID, &user.Role, &user.FirstName, &user.LastName,
		&user.Nickname, &user.Phone, &user.IsBlocked, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		utils.LogError(err)
		return nil, fmt.Errorf("failed to get user: %v", err)
	}
	return user, nil
}

// GetUserByRole retrieves a user by ChatID and role
func (r *PostgresRepository) GetUserByRole(chatID int64, role string) (*models.User, error) {
	query := `
		SELECT id, chat_id, role, first_name, last_name, nickname, phone, is_blocked, created_at, updated_at
		FROM users
		WHERE chat_id = $1 AND role = $2
	`
	user := &models.User{}
	err := r.db.Conn().QueryRow(query, chatID, role).Scan(
		&user.ID, &user.ChatID, &user.Role, &user.FirstName, &user.LastName,
		&user.Nickname, &user.Phone, &user.IsBlocked, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found with role %s", role)
	}
	if err != nil {
		utils.LogError(err)
		return nil, fmt.Errorf("failed to get user by role: %v", err)
	}
	return user, nil
}

// ListUsersByRole retrieves users by role
func (r *PostgresRepository) ListUsersByRole(role string) ([]models.User, error) {
	query := `
		SELECT id, chat_id, role, first_name, last_name, nickname, phone, is_blocked, created_at, updated_at
		FROM users
		WHERE role = $1 AND is_blocked = FALSE
	`
	rows, err := r.db.Conn().Query(query, role)
	if err != nil {
		utils.LogError(err)
		return nil, fmt.Errorf("failed to list users by role: %v", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.ID, &user.ChatID, &user.Role, &user.FirstName, &user.LastName,
			&user.Nickname, &user.Phone, &user.IsBlocked, &user.CreatedAt, &user.UpdatedAt,
		); err != nil {
			utils.LogError(err)
			continue
		}
		users = append(users, user)
	}
	return users, nil
}

// UpdateUser updates an existing user
func (r *PostgresRepository) UpdateUser(user *models.User) error {
	query := `
		UPDATE users
		SET role = $1, first_name = $2, last_name = $3, nickname = $4, phone = $5, 
		    is_blocked = $6, created_at = $7, updated_at = $8
		WHERE chat_id = $9
	`
	_, err := r.db.Conn().Exec(
		query,
		user.Role, user.FirstName, user.LastName, user.Nickname, user.Phone,
		user.IsBlocked, user.CreatedAt, user.UpdatedAt, user.ChatID,
	)
	if err != nil {
		utils.LogError(err)
		return fmt.Errorf("failed to update user: %v", err)
	}
	return nil
}

// DeleteUser deletes a user
func (r *PostgresRepository) DeleteUser(chatID int64) error {
	query := `DELETE FROM users WHERE chat_id = $1`
	_, err := r.db.Conn().Exec(query, chatID)
	if err != nil {
		utils.LogError(err)
		return fmt.Errorf("failed to delete user: %v", err)
	}
	return nil
}