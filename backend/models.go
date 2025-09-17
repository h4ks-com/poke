package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// User represents a bank user
type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	AccountNumber string   `json:"account_number"`
	Balance      float64   `json:"balance"`
	CreatedAt    time.Time `json:"created_at"`
}

// Transaction represents a banking transaction
type Transaction struct {
	ID            int       `json:"id"`
	FromUserID    int       `json:"from_user_id"`
	ToUserID      int       `json:"to_user_id"`
	Amount        float64   `json:"amount"`
	TransactionType string  `json:"transaction_type"` // "transfer", "deposit", "withdrawal"
	Description   string    `json:"description"`
	Status        string    `json:"status"` // "pending", "completed", "failed"
	CreatedAt     time.Time `json:"created_at"`
	
	// Additional fields for display
	FromUsername  string    `json:"from_username,omitempty"`
	ToUsername    string    `json:"to_username,omitempty"`
}

// PaymentRequest represents a money request between users
type PaymentRequest struct {
	ID          int       `json:"id"`
	FromUserID  int       `json:"from_user_id"`
	ToUserID    int       `json:"to_user_id"`
	Amount      float64   `json:"amount"`
	Reason      string    `json:"reason"`
	Message     string    `json:"message"`
	Status      string    `json:"status"` // "pending", "approved", "rejected"
	CreatedAt   time.Time `json:"created_at"`
	
	// Additional fields for display
	FromUsername string    `json:"from_username,omitempty"`
	ToUsername   string    `json:"to_username,omitempty"`
}

// UserSession represents an active user session
type UserSession struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	SessionToken string    `json:"session_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

// Card represents a user's virtual bank card
type Card struct {
	ID                int       `json:"id"`
	UserID            int       `json:"user_id"`
	CardNumber        string    `json:"card_number"`
	ExpiryDate        string    `json:"expiry_date"`
	RefreshSeed       int       `json:"refresh_seed"`
	LastRefreshDate   *time.Time `json:"last_refresh_date"`
	IsActive          bool      `json:"is_active"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// InitDB initializes the database connection and creates tables
func InitDB() (*sql.DB, error) {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./viridian_bank.db"
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Create tables if they don't exist
	if err := createTables(db); err != nil {
		return nil, err
	}

	return db, nil
}

// createTables creates the necessary database tables
func createTables(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			account_number TEXT UNIQUE NOT NULL,
			balance REAL DEFAULT 0.00,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		
		`CREATE TABLE IF NOT EXISTS transactions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			from_user_id INTEGER REFERENCES users(id),
			to_user_id INTEGER REFERENCES users(id),
			amount REAL NOT NULL,
			transaction_type TEXT NOT NULL,
			description TEXT,
			status TEXT DEFAULT 'completed',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		
		`CREATE TABLE IF NOT EXISTS payment_requests (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			from_user_id INTEGER REFERENCES users(id),
			to_user_id INTEGER REFERENCES users(id),
			amount REAL NOT NULL,
			reason TEXT NOT NULL,
			message TEXT,
			status TEXT DEFAULT 'pending',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		
		`CREATE TABLE IF NOT EXISTS user_sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER REFERENCES users(id),
			session_token TEXT UNIQUE NOT NULL,
			expires_at DATETIME NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		
		`CREATE TABLE IF NOT EXISTS cards (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER REFERENCES users(id),
			card_number TEXT NOT NULL,
			expiry_date TEXT NOT NULL,
			refresh_seed INTEGER DEFAULT 0,
			last_refresh_date DATETIME,
			is_active BOOLEAN DEFAULT TRUE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		
		// Create indexes for better performance
		`CREATE INDEX IF NOT EXISTS idx_users_username ON users(username)`,
		`CREATE INDEX IF NOT EXISTS idx_users_account_number ON users(account_number)`,
		`CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)`,
		`CREATE INDEX IF NOT EXISTS idx_transactions_from_user ON transactions(from_user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_transactions_to_user ON transactions(to_user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_payment_requests_from_user ON payment_requests(from_user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_payment_requests_to_user ON payment_requests(to_user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_sessions_token ON user_sessions(session_token)`,
		`CREATE INDEX IF NOT EXISTS idx_sessions_user ON user_sessions(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_cards_user ON cards(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_cards_active ON cards(is_active)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query: %v", err)
		}
	}

	return nil
}
