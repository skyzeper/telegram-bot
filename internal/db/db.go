package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

// DB represents the database connection
type DB struct {
	conn *sql.DB
}

// Config holds database connection parameters
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// NewDB creates a new database connection
func NewDB(cfg Config) (*DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName,
	)

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	db := &DB{conn: conn}
	if err := db.initSchema(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to initialize schema: %v", err)
	}

	return db, nil
}

// initSchema creates database tables
func (db *DB) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		chat_id BIGINT UNIQUE NOT NULL,
		role VARCHAR(50) NOT NULL,
		first_name VARCHAR(100),
		last_name VARCHAR(100),
		nickname VARCHAR(100),
		phone VARCHAR(20),
		is_blocked BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL
	);

	CREATE TABLE IF NOT EXISTS orders (
		id SERIAL PRIMARY KEY,
		user_id BIGINT NOT NULL,
		category VARCHAR(100) NOT NULL,
		subcategory VARCHAR(100) NOT NULL,
		photos TEXT[],
		video TEXT,
		date TIMESTAMP,
		time TIME,
		phone VARCHAR(20) NOT NULL,
		address TEXT NOT NULL,
		description TEXT,
		status VARCHAR(50) NOT NULL,
		reason TEXT,
		cost FLOAT DEFAULT 0,
		payment_method VARCHAR(50),
		payment_confirmed BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL,
		confirmed BOOLEAN DEFAULT FALSE,
		FOREIGN KEY (user_id) REFERENCES users(chat_id)
	);

	CREATE TABLE IF NOT EXISTS executors (
		id SERIAL PRIMARY KEY,
		order_id INTEGER NOT NULL,
		user_id BIGINT NOT NULL,
		role VARCHAR(50) NOT NULL,
		confirmed BOOLEAN DEFAULT FALSE,
		notified BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP NOT NULL,
		FOREIGN KEY (order_id) REFERENCES orders(id),
		FOREIGN KEY (user_id) REFERENCES users(chat_id)
	);

	CREATE TABLE IF NOT EXISTS payments (
		id SERIAL PRIMARY KEY,
		order_id INTEGER NOT NULL,
		user_id BIGINT NOT NULL,
		amount FLOAT NOT NULL,
		method VARCHAR(50) NOT NULL,
		driver_id BIGINT,
		confirmed BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP NOT NULL,
		FOREIGN KEY (order_id) REFERENCES orders(id),
		FOREIGN KEY (user_id) REFERENCES users(chat_id),
		FOREIGN KEY (driver_id) REFERENCES users(chat_id)
	);

	CREATE TABLE IF NOT EXISTS reviews (
		id SERIAL PRIMARY KEY,
		order_id INTEGER NOT NULL,
		user_id BIGINT NOT NULL,
		rating INTEGER NOT NULL,
		comment TEXT,
		created_at TIMESTAMP NOT NULL,
		FOREIGN KEY (order_id) REFERENCES orders(id),
		FOREIGN KEY (user_id) REFERENCES users(chat_id),
		UNIQUE (order_id)
	);

	CREATE TABLE IF NOT EXISTS referrals (
		id SERIAL PRIMARY KEY,
		inviter_id BIGINT NOT NULL,
		invitee_id BIGINT NOT NULL,
		order_id INTEGER,
		payout_requested BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP NOT NULL,
		FOREIGN KEY (inviter_id) REFERENCES users(chat_id),
		FOREIGN KEY (invitee_id) REFERENCES users(chat_id),
		FOREIGN KEY (order_id) REFERENCES orders(id),
		UNIQUE (invitee_id)
	);

	CREATE TABLE IF NOT EXISTS messages (
		id SERIAL PRIMARY KEY,
		user_id BIGINT NOT NULL,
		operator_id BIGINT,
		message TEXT NOT NULL,
		is_from_user BOOLEAN NOT NULL,
		created_at TIMESTAMP NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(chat_id),
		FOREIGN KEY (operator_id) REFERENCES users(chat_id)
	);

	CREATE TABLE IF NOT EXISTS notifications (
		id SERIAL PRIMARY KEY,
		user_id BIGINT NOT NULL,
		type VARCHAR(100) NOT NULL,
		message TEXT NOT NULL,
		sent_at TIMESTAMP,
		created_at TIMESTAMP NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(chat_id)
	);

	CREATE TABLE IF NOT EXISTS accounting_records (
		id SERIAL PRIMARY KEY,
		order_id INTEGER,
		user_id BIGINT NOT NULL,
		type VARCHAR(50) NOT NULL,
		amount FLOAT NOT NULL,
		description TEXT,
		created_at TIMESTAMP NOT NULL,
		FOREIGN KEY (order_id) REFERENCES orders(id),
		FOREIGN KEY (user_id) REFERENCES users(chat_id)
	);
	`

	_, err := db.conn.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to execute schema: %v", err)
	}
	return nil
}

// Conn returns the database connection
func (db *DB) Conn() *sql.DB {
	return db.conn
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}
