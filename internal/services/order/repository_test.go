package order_test

import (
	"database/sql"
	"testing"
	"time"
	"github.com/skyzeper/telegram-bot/internal/db"
	"github.com/skyzeper/telegram-bot/internal/models"
	"github.com/skyzeper/telegram-bot/internal/services/order"
	"github.com/stretchr/testify/assert"
	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) (*db.DB, func()) {
	sqlDB, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}

	_, err = sqlDB.Exec(`
		CREATE TABLE orders (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			category TEXT,
			subcategory TEXT,
			photos TEXT,
			video TEXT,
			date DATETIME,
			time DATETIME,
			phone TEXT,
			address TEXT,
			description TEXT,
			status TEXT,
			reason TEXT,
			cost REAL,
			payment_method TEXT,
			payment_confirmed BOOLEAN,
			created_at DATETIME,
			updated_at DATETIME,
			confirmed BOOLEAN
		);
		CREATE TABLE executors (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			order_id INTEGER,
			user_id INTEGER,
			role TEXT,
			confirmed BOOLEAN,
			notified BOOLEAN,
			created_at DATETIME
		);
	`)
	if err != nil {
		t.Fatalf("failed to create tables: %v", err)
	}

	dbConn := &db.DB{SQLDB: sqlDB}
	return dbConn, func() { sqlDB.Close() }
}

func TestPostgresRepository_CreateOrder(t *testing.T) {
	dbConn, cleanup := setupTestDB(t)
	defer cleanup()

	repo := order.NewPostgresRepository(dbConn)

	t.Run("ValidOrder", func(t *testing.T) {
		order := &models.Order{
			UserID:      123,
			Category:    "вывоз мусора",
			Subcategory: "строительный мусор",
			Phone:       "+1234567890",
			Address:     "ул. Тестовая, 1",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		err := repo.CreateOrder(order)
		assert.NoError(t, err)
		assert.NotZero(t, order.ID)

		// Verify in DB
		var count int
		err = dbConn.Conn().QueryRow("SELECT COUNT(*) FROM orders WHERE id = ?", order.ID).Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
	})
}

func TestPostgresRepository_GetOrder(t *testing.T) {
	dbConn, cleanup := setupTestDB(t)
	defer cleanup()

	repo := order.NewPostgresRepository(dbConn)

	t.Run("ExistingOrder", func(t *testing.T) {
		order := &models.Order{
			UserID:      123,
			Category:    "вывоз мусора",
			Subcategory: "строительный мусор",
			Phone:       "+1234567890",
			Address:     "ул. Тестовая, 1",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		err := repo.CreateOrder(order)
		assert.NoError(t, err)

		result, err := repo.GetOrder(order.ID)
		assert.NoError(t, err)
		assert.Equal(t, order.UserID, result.UserID)
		assert.Equal(t, order.Category, result.Category)
	})

	t.Run("NonExistingOrder", func(t *testing.T) {
		result, err := repo.GetOrder(999)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "order not found", err.Error())
	})
}

func TestPostgresRepository_GetOrdersByStatus(t *testing.T) {
	dbConn, cleanup := setupTestDB(t)
	defer cleanup()

	repo := order.NewPostgresRepository(dbConn)

	t.Run("ExistingOrders", func(t *testing.T) {
		order1 := &models.Order{
			UserID:      123,
			Category:    "вывоз мусора",
			Subcategory: "строительный мусор",
			Phone:       "+1234567890",
			Address:     "ул. Тестовая, 1",
			Status:      "new",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		order2 := &models.Order{
			UserID:      124,
			Category:    "демонтаж",
			Subcategory: "стены",
			Phone:       "+1234567891",
			Address:     "ул. Тестовая, 2",
			Status:      "new",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		err := repo.CreateOrder(order1)
		assert.NoError(t, err)
		err = repo.CreateOrder(order2)
		assert.NoError(t, err)

		orders, err := repo.GetOrdersByStatus("new")
		assert.NoError(t, err)
		assert.Len(t, orders, 2)
		assert.Equal(t, order1.UserID, orders[0].UserID)
		assert.Equal(t, order2.UserID, orders[1].UserID)
	})

	t.Run("NoOrders", func(t *testing.T) {
		orders, err := repo.GetOrdersByStatus("completed")
		assert.NoError(t, err)
		assert.Empty(t, orders)
	})
}