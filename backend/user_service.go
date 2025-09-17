package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// UserService handles user-related database operations
type UserService struct {
	db *sql.DB
}

// NewUserService creates a new UserService
func NewUserService(db *sql.DB) *UserService {
	return &UserService{db: db}
}

// CreateUser creates a new user account
func (s *UserService) CreateUser(username, email, password string) (*User, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Generate unique account number
	accountNumber := s.generateAccountNumber()

	// Create user with initial balance
	query := `
		INSERT INTO users (username, email, password_hash, account_number, balance)
		VALUES (?, ?, ?, ?, ?)
	`

	result, err := s.db.Exec(query, username, email, string(hashedPassword), accountNumber, 1000.00)
	if err != nil {
		return nil, err
	}

	// Get the last insert ID
	userID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Fetch the created user
	user := &User{}
	selectQuery := `SELECT id, username, email, account_number, balance, created_at FROM users WHERE id = ?`
	err = s.db.QueryRow(selectQuery, userID).Scan(
		&user.ID, &user.Username, &user.Email, &user.AccountNumber, &user.Balance, &user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByUsername retrieves a user by username
func (s *UserService) GetUserByUsername(username string) (*User, error) {
	query := `
		SELECT id, username, email, password_hash, account_number, balance, created_at
		FROM users WHERE username = ?
	`

	user := &User{}
	err := s.db.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.AccountNumber, &user.Balance, &user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(id int) (*User, error) {
	query := `
		SELECT id, username, email, password_hash, account_number, balance, created_at
		FROM users WHERE id = ?
	`

	user := &User{}
	err := s.db.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.AccountNumber, &user.Balance, &user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByAccountNumber retrieves a user by account number
func (s *UserService) GetUserByAccountNumber(accountNumber string) (*User, error) {
	query := `
		SELECT id, username, email, password_hash, account_number, balance, created_at
		FROM users WHERE account_number = ?
	`

	user := &User{}
	err := s.db.QueryRow(query, accountNumber).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.AccountNumber, &user.Balance, &user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdatePassword updates a user's password
func (s *UserService) UpdatePassword(userID int, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `UPDATE users SET password_hash = ? WHERE id = ?`
	_, err = s.db.Exec(query, string(hashedPassword), userID)
	return err
}

// VerifyPassword verifies a user's password
func (s *UserService) VerifyPassword(user *User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	return err == nil
}

// CreateSession creates a new user session
func (s *UserService) CreateSession(userID int, token string, expiresAt time.Time) error {
	query := `INSERT INTO user_sessions (user_id, session_token, expires_at) VALUES (?, ?, ?)`
	_, err := s.db.Exec(query, userID, token, expiresAt)
	return err
}

// GetSessionByToken retrieves a session by token
func (s *UserService) GetSessionByToken(token string) (*UserSession, error) {
	query := `
		SELECT id, user_id, session_token, expires_at, created_at
		FROM user_sessions 
		WHERE session_token = ? AND expires_at > CURRENT_TIMESTAMP
	`

	session := &UserSession{}
	err := s.db.QueryRow(query, token).Scan(
		&session.ID, &session.UserID, &session.SessionToken, &session.ExpiresAt, &session.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return session, nil
}

// DeleteSession deletes a session
func (s *UserService) DeleteSession(token string) error {
	query := `DELETE FROM user_sessions WHERE session_token = ?`
	_, err := s.db.Exec(query, token)
	return err
}

// DeleteAllUserSessions deletes all sessions for a user
func (s *UserService) DeleteAllUserSessions(userID int) error {
	query := `DELETE FROM user_sessions WHERE user_id = ?`
	_, err := s.db.Exec(query, userID)
	return err
}

// generateAccountNumber generates a unique account number
func (s *UserService) generateAccountNumber() string {
	// Generate a random 10-digit account number
	rand.Seed(time.Now().UnixNano())
	accountNumber := fmt.Sprintf("%010d", rand.Intn(10000000000))
	
	// Check if it already exists, regenerate if needed
	for {
		var exists bool
		query := `SELECT EXISTS(SELECT 1 FROM users WHERE account_number = ?)`
		s.db.QueryRow(query, accountNumber).Scan(&exists)
		
		if !exists {
			break
		}
		
		accountNumber = fmt.Sprintf("%010d", rand.Intn(10000000000))
	}
	
	return accountNumber
}
