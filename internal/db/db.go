package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"bot/internal/config"
	_ "github.com/lib/pq"
)

func Connect(cfg *config.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := initTables(db); err != nil {
		return nil, err
	}

	go reconnectLoop(db, connStr)

	return db, nil
}

func initTables(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			chat_id BIGINT UNIQUE,
			role TEXT NOT NULL CHECK (role IN ('user','operator','main_operator','driver','loader','owner')),
			first_name VARCHAR(100),
			last_name VARCHAR(100),
			nickname VARCHAR(100) UNIQUE,
			phone VARCHAR(20),
			is_blocked BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id),
			category TEXT NOT NULL CHECK (category IN ('waste_removal','demolition','construction_materials','other')),
			subcategory VARCHAR(100),
			photos TEXT[],
			video TEXT,
			date DATE,
			time TIME,
			phone VARCHAR(20),
			address TEXT,
			description TEXT,
			status TEXT NOT NULL CHECK (status IN ('new','awaiting_payment','in_progress','completed','canceled')),
			reason TEXT,
			cost DECIMAL,
			payment_method TEXT CHECK (payment_method IN ('cash','transfer','driver_card')),
			payment_confirmed BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS executors (
			id SERIAL PRIMARY KEY,
			order_id INTEGER REFERENCES orders(id),
			user_id INTEGER REFERENCES users(id),
			role TEXT NOT NULL CHECK (role IN ('driver','loader')),
			confirmed BOOLEAN DEFAULT FALSE,
			notified BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS referrals (
			id SERIAL PRIMARY KEY,
			inviter_id INTEGER REFERENCES users(id),
			invitee_id INTEGER REFERENCES users(id),
			order_id INTEGER REFERENCES orders(id),
			payout_requested BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS chat_messages (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id),
			operator_id INTEGER REFERENCES users(id),
			message TEXT,
			is_from_user BOOLEAN,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS reviews (
			id SERIAL PRIMARY KEY,
			order_id INTEGER REFERENCES orders(id),
			user_id INTEGER REFERENCES users(id),
			rating INTEGER CHECK (rating BETWEEN 1 AND 5),
			comment TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS payments (
			id SERIAL PRIMARY KEY,
			order_id INTEGER REFERENCES orders(id),
			user_id INTEGER REFERENCES users(id),
			amount DECIMAL,
			method TEXT NOT NULL CHECK (method IN ('cash','transfer','driver_card')),
			driver_id INTEGER REFERENCES users(id),
			confirmed BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS accounting_records (
			id SERIAL PRIMARY KEY,
			order_id INTEGER REFERENCES orders(id),
			user_id INTEGER REFERENCES users(id),
			type TEXT NOT NULL CHECK (type IN ('fuel','other','salary','referral')),
			amount DECIMAL,
			description TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE INDEX IF NOT EXISTS idx_users_chat_id ON users(chat_id);`,
		`CREATE INDEX IF NOT EXISTS idx_orders_user_id_status ON orders(user_id, status);`,
		`CREATE INDEX IF NOT EXISTS idx_executors_order_id ON executors(order_id);`,
		`CREATE INDEX IF NOT EXISTS idx_referrals_inviter_id ON referrals(inviter_id);`,
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			return fmt.Errorf("failed to execute query: %v", err)
		}
	}

	return nil
}

func reconnectLoop(db *sql.DB, connStr string) {
	for {
		time.Sleep(60 * time.Second)
		if err := db.Ping(); err != nil {
			log.Printf("Database ping failed: %v. Reconnecting...", err)
			newDB, err := sql.Open("postgres", connStr)
			if err == nil && newDB.Ping() == nil {
				db.Close()
				*db = *newDB
				log.Println("Reconnected to database")
			}
		}
	}
}
