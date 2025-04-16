package order

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

// CreateOrder creates a new order
func (r *PostgresRepository) CreateOrder(order *models.Order) error {
	query := `
		INSERT INTO orders (
			user_id, category, subcategory, photos, video, date, time, phone, address, 
			description, status, reason, cost, payment_method, payment_confirmed, 
			created_at, updated_at, confirmed
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
		RETURNING id
	`
	var date, timeVal sql.NullTime
	var video, reason, paymentMethod sql.NullString
	if !order.Date.IsZero() {
		date.Valid = true
		date.Time = order.Date
	}
	if !order.Time.IsZero() {
		timeVal.Valid = true
		timeVal.Time = order.Time
	}
	if order.Video != "" {
		video.Valid = true
		video.String = order.Video
	}
	if order.Reason != "" {
		reason.Valid = true
		reason.String = order.Reason
	}
	if order.PaymentMethod != "" {
		paymentMethod.Valid = true
		paymentMethod.String = order.PaymentMethod
	}
	err := r.db.Conn().QueryRow(
		query,
		order.UserID, order.Category, order.Subcategory, order.Photos, video,
		date, timeVal, order.Phone, order.Address, order.Description,
		order.Status, reason, order.Cost, paymentMethod, order.PaymentConfirmed,
		order.CreatedAt, order.UpdatedAt, order.Confirmed,
	).Scan(&order.ID)
	if err != nil {
		utils.LogError(err)
		return fmt.Errorf("failed to create order: %v", err)
	}
	return nil
}

// GetOrder retrieves an order by ID
func (r *PostgresRepository) GetOrder(id int) (*models.Order, error) {
	query := `
		SELECT id, user_id, category, subcategory, photos, video, date, time, phone, 
		       address, description, status, reason, cost, payment_method, payment_confirmed, 
		       created_at, updated_at, confirmed
		FROM orders
		WHERE id = $1
	`
	order := &models.Order{}
	var date, timeVal sql.NullTime
	var video, reason, paymentMethod sql.NullString
	err := r.db.Conn().QueryRow(query, id).Scan(
		&order.ID, &order.UserID, &order.Category, &order.Subcategory, &order.Photos,
		&video, &date, &timeVal, &order.Phone, &order.Address,
		&order.Description, &order.Status, &reason, &order.Cost, &paymentMethod,
		&order.PaymentConfirmed, &order.CreatedAt, &order.UpdatedAt, &order.Confirmed,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("order not found")
	}
	if err != nil {
		utils.LogError(err)
		return nil, fmt.Errorf("failed to get order: %v", err)
	}
	if video.Valid {
		order.Video = video.String
	}
	if date.Valid {
		order.Date = date.Time
	}
	if timeVal.Valid {
		order.Time = timeVal.Time
	}
	if reason.Valid {
		order.Reason = reason.String
	}
	if paymentMethod.Valid {
		order.PaymentMethod = paymentMethod.String
	}

	// Fetch executors
	executors, err := r.getExecutors(id)
	if err != nil {
		utils.LogError(err)
	}
	order.Executors = executors

	return order, nil
}

// GetOrdersByStatus retrieves orders by status
func (r *PostgresRepository) GetOrdersByStatus(status string) ([]models.Order, error) {
	query := `
		SELECT id, user_id, category, subcategory, photos, video, date, time, phone, 
		       address, description, status, reason, cost, payment_method, payment_confirmed, 
		       created_at, updated_at, confirmed
		FROM orders
		WHERE status = $1
	`
	rows, err := r.db.Conn().Query(query, status)
	if err != nil {
		utils.LogError(err)
		return nil, fmt.Errorf("failed to get orders by status: %v", err)
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		var date, timeVal sql.NullTime
		var video, reason, paymentMethod sql.NullString
		if err := rows.Scan(
			&order.ID, &order.UserID, &order.Category, &order.Subcategory, &order.Photos,
			&video, &date, &timeVal, &order.Phone, &order.Address,
			&order.Description, &order.Status, &reason, &order.Cost, &paymentMethod,
			&order.PaymentConfirmed, &order.CreatedAt, &order.UpdatedAt, &order.Confirmed,
		); err != nil {
			utils.LogError(err)
			continue
		}
		if video.Valid {
			order.Video = video.String
		}
		if date.Valid {
			order.Date = date.Time
		}
		if timeVal.Valid {
			order.Time = timeVal.Time
		}
		if reason.Valid {
			order.Reason = reason.String
		}
		if paymentMethod.Valid {
			order.PaymentMethod = paymentMethod.String
		}
		// Fetch executors
		executors, err := r.getExecutors(order.ID)
		if err != nil {
			utils.LogError(err)
		}
		order.Executors = executors
		orders = append(orders, order)
	}
	return orders, nil
}

// GetOrdersByStatusAndCategory retrieves orders by status and category
func (r *PostgresRepository) GetOrdersByStatusAndCategory(status, category string) ([]models.Order, error) {
	query := `
		SELECT id, user_id, category, subcategory, photos, video, date, time, phone, 
		       address, description, status, reason, cost, payment_method, payment_confirmed, 
		       created_at, updated_at, confirmed
		FROM orders
		WHERE status = $1 AND category = $2
	`
	rows, err := r.db.Conn().Query(query, status, category)
	if err != nil {
		utils.LogError(err)
		return nil, fmt.Errorf("failed to get orders by status and category: %v", err)
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		var date, timeVal sql.NullTime
		var video, reason, paymentMethod sql.NullString
		if err := rows.Scan(
			&order.ID, &order.UserID, &order.Category, &order.Subcategory, &order.Photos,
			&video, &date, &timeVal, &order.Phone, &order.Address,
			&order.Description, &order.Status, &reason, &order.Cost, &paymentMethod,
			&order.PaymentConfirmed, &order.CreatedAt, &order.UpdatedAt, &order.Confirmed,
		); err != nil {
			utils.LogError(err)
			continue
		}
		if video.Valid {
			order.Video = video.String
		}
		if date.Valid {
			order.Date = date.Time
		}
		if timeVal.Valid {
			order.Time = timeVal.Time
		}
		if reason.Valid {
			order.Reason = reason.String
		}
		if paymentMethod.Valid {
			order.PaymentMethod = paymentMethod.String
		}
		// Fetch executors
		executors, err := r.getExecutors(order.ID)
		if err != nil {
			utils.LogError(err)
		}
		order.Executors = executors
		orders = append(orders, order)
	}
	return orders, nil
}

// GetExecutorOrders retrieves orders assigned to an executor
func (r *PostgresRepository) GetExecutorOrders(userID int64) ([]models.Order, error) {
	query := `
		SELECT o.id, o.user_id, o.category, o.subcategory, o.photos, o.video, o.date, o.time, 
		       o.phone, o.address, o.description, o.status, o.reason, o.cost, o.payment_method, 
		       o.payment_confirmed, o.created_at, o.updated_at, o.confirmed
		FROM orders o
		JOIN executors e ON o.id = e.order_id
		WHERE e.user_id = $1
	`
	rows, err := r.db.Conn().Query(query, userID)
	if err != nil {
		utils.LogError(err)
		return nil, fmt.Errorf("failed to get executor orders: %v", err)
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		var date, timeVal sql.NullTime
		var video, reason, paymentMethod sql.NullString
		if err := rows.Scan(
			&order.ID, &order.UserID, &order.Category, &order.Subcategory, &order.Photos,
			&video, &date, &timeVal, &order.Phone, &order.Address,
			&order.Description, &order.Status, &reason, &order.Cost, &paymentMethod,
			&order.PaymentConfirmed, &order.CreatedAt, &order.UpdatedAt, &order.Confirmed,
		); err != nil {
			utils.LogError(err)
			continue
		}
		if video.Valid {
			order.Video = video.String
		}
		if date.Valid {
			order.Date = date.Time
		}
		if timeVal.Valid {
			order.Time = timeVal.Time
		}
		if reason.Valid {
			order.Reason = reason.String
		}
		if paymentMethod.Valid {
			order.PaymentMethod = paymentMethod.String
		}
		// Fetch executors
		executors, err := r.getExecutors(order.ID)
		if err != nil {
			utils.LogError(err)
		}
		order.Executors = executors
		orders = append(orders, order)
	}
	return orders, nil
}

// GetOrderClientID retrieves the client ID for an order
func (r *PostgresRepository) GetOrderClientID(orderID int) (int64, error) {
	query := `SELECT user_id FROM orders WHERE id = $1`
	var userID int64
	err := r.db.Conn().QueryRow(query, orderID).Scan(&userID)
	if err == sql.ErrNoRows {
		return 0, fmt.Errorf("order not found")
	}
	if err != nil {
		utils.LogError(err)
		return 0, fmt.Errorf("failed to get order client ID: %v", err)
	}
	return userID, nil
}

// UpdateOrder updates an existing order
func (r *PostgresRepository) UpdateOrder(order *models.Order) error {
	query := `
		UPDATE orders
		SET user_id = $1, category = $2, subcategory = $3, photos = $4, video = $5, 
		    date = $6, time = $7, phone = $8, address = $9, description = $10, 
		    status = $11, reason = $12, cost = $13, payment_method = $14, 
		    payment_confirmed = $15, created_at = $16, updated_at = $17, confirmed = $18
		WHERE id = $19
	`
	var date, timeVal sql.NullTime
	var video, reason, paymentMethod sql.NullString
	if !order.Date.IsZero() {
		date.Valid = true
		date.Time = order.Date
	}
	if !order.Time.IsZero() {
		timeVal.Valid = true
		timeVal.Time = order.Time
	}
	if order.Video != "" {
		video.Valid = true
		video.String = order.Video
	}
	if order.Reason != "" {
		reason.Valid = true
		reason.String = order.Reason
	}
	if order.PaymentMethod != "" {
		paymentMethod.Valid = true
		paymentMethod.String = order.PaymentMethod
	}
	_, err := r.db.Conn().Exec(
		query,
		order.UserID, order.Category, order.Subcategory, order.Photos, video,
		date, timeVal, order.Phone, order.Address, order.Description,
		order.Status, reason, order.Cost, paymentMethod, order.PaymentConfirmed,
		order.CreatedAt, order.UpdatedAt, order.Confirmed, order.ID,
	)
	if err != nil {
		utils.LogError(err)
		return fmt.Errorf("failed to update order: %v", err)
	}
	return nil
}

// ConfirmOrder confirms an order
func (r *PostgresRepository) ConfirmOrder(orderID int, userID int64) error {
	// Update executor confirmation
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

	// Check if all executors confirmed
	query = `
		SELECT COUNT(*) AS total, SUM(CASE WHEN confirmed THEN 1 ELSE 0 END) AS confirmed
		FROM executors
		WHERE order_id = $1
	`
	var total, confirmed int
	err = r.db.Conn().QueryRow(query, orderID).Scan(&total, &confirmed)
	if err != nil {
		utils.LogError(err)
		return fmt.Errorf("failed to check executors: %v", err)
	}

	if total > 0 && total == confirmed {
		query = `
			UPDATE orders
			SET status = 'completed', confirmed = TRUE
			WHERE id = $1
		`
		_, err = r.db.Conn().Exec(query, orderID)
		if err != nil {
			utils.LogError(err)
			return fmt.Errorf("failed to confirm order: %v", err)
		}
	}
	return nil
}

// getExecutors retrieves executors for an order
func (r *PostgresRepository) getExecutors(orderID int) ([]models.Executor, error) {
	query := `
		SELECT id, order_id, user_id, role, confirmed, notified, created_at
		FROM executors
		WHERE order_id = $1
	`
	rows, err := r.db.Conn().Query(query, orderID)
	if err != nil {
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