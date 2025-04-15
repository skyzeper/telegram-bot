package security

import (
	"context"
	"database/sql"

	"github.com/skyzeper/telegram-bot/internal/services/user"
)

func CheckAccess(ctx context.Context, chatID int64, db *sql.DB) bool {
	userService := user.NewService(db, nil)
	user, err := userService.GetUser(ctx, chatID)
	if err != nil {
		return false
	}
	return !user.IsBlocked
}
