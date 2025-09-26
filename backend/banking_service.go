package main

import (
	"database/sql"
	"fmt"
	"time"
)

// BankingService handles banking-related database operations
type BankingService struct {
	db *sql.DB
}

// NewBankingService creates a new BankingService
func NewBankingService(db *sql.DB) *BankingService {
	return &BankingService{db: db}
}

// GetUserBalance gets the current balance for a user
func (s *BankingService) GetUserBalance(userID int) (float64, error) {
	query := `SELECT balance FROM users WHERE id = ?`
	var balance float64
	err := s.db.QueryRow(query, userID).Scan(&balance)
	return balance, err
}

// Transfer transfers money between users
func (s *BankingService) Transfer(fromUserID int, toAccountNumber string, amount float64, description string) (*Transaction, error) {
	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Get sender details
	var fromBalance float64
	err = tx.QueryRow(`SELECT balance FROM users WHERE id = ?`, fromUserID).Scan(&fromBalance)
	if err != nil {
		return nil, err
	}

	// Check sufficient balance
	if fromBalance < amount {
		return nil, fmt.Errorf("insufficient balance")
	}

	// Get recipient details
	var toUserID int
	err = tx.QueryRow(`SELECT id FROM users WHERE account_number = ?`, toAccountNumber).Scan(&toUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("recipient account not found")
		}
		return nil, err
	}

	// Check not transferring to self
	if fromUserID == toUserID {
		return nil, fmt.Errorf("cannot transfer to yourself")
	}

	// Update sender balance
	_, err = tx.Exec(`UPDATE users SET balance = balance - ? WHERE id = ?`, amount, fromUserID)
	if err != nil {
		return nil, err
	}

	// Update recipient balance
	_, err = tx.Exec(`UPDATE users SET balance = balance + ? WHERE id = ?`, amount, toUserID)
	if err != nil {
		return nil, err
	}

	// Create transaction record
	var transactionID int
	err = tx.QueryRow(`
		INSERT INTO transactions (from_user_id, to_user_id, amount, transaction_type, description, status, created_at)
		VALUES (?, ?, ?, 'transfer', ?, 'completed', ?)
		RETURNING id
	`, fromUserID, toUserID, amount, description, time.Now()).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	// Reset PokéBank balance if involved in transaction
	err = s.resetPokeBankBalanceInTx(tx, fromUserID, toUserID)
	if err != nil {
		return nil, err
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	// Get the created transaction with user details
	return s.GetTransactionByID(transactionID)
}

// CreatePaymentRequest creates a new payment request
func (s *BankingService) CreatePaymentRequest(fromUserID int, toAccountNumber string, amount float64, reason, message string) (*PaymentRequest, error) {
	// Get recipient user ID
	var toUserID int
	err := s.db.QueryRow(`SELECT id FROM users WHERE account_number = ?`, toAccountNumber).Scan(&toUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	// Check not requesting from self
	if fromUserID == toUserID {
		return nil, fmt.Errorf("cannot request money from yourself")
	}

	// Create payment request
	query := `
		INSERT INTO payment_requests (from_user_id, to_user_id, amount, reason, message, status)
		VALUES (?, ?, ?, ?, ?, 'pending')
		RETURNING id, from_user_id, to_user_id, amount, reason, message, status, created_at
	`

	request := &PaymentRequest{}
	err = s.db.QueryRow(query, fromUserID, toUserID, amount, reason, message).Scan(
		&request.ID, &request.FromUserID, &request.ToUserID, &request.Amount,
		&request.Reason, &request.Message, &request.Status, &request.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	// Add usernames for display
	request.FromUsername, _ = s.getUsernameByID(request.FromUserID)
	request.ToUsername, _ = s.getUsernameByID(request.ToUserID)

	return request, nil
}

// GetUserTransactions gets transactions for a user
func (s *BankingService) GetUserTransactions(userID int, limit int) ([]Transaction, error) {
	query := `
		SELECT t.id, t.from_user_id, t.to_user_id, 
		       CASE 
		           WHEN t.from_user_id = ? THEN -t.amount 
		           ELSE t.amount 
		       END as amount,
		       t.transaction_type, t.description, t.status, t.created_at,
		       u1.username as from_username, u2.username as to_username
		FROM transactions t
		LEFT JOIN users u1 ON t.from_user_id = u1.id
		LEFT JOIN users u2 ON t.to_user_id = u2.id
		WHERE t.from_user_id = ? OR t.to_user_id = ?
		ORDER BY t.created_at DESC, t.id DESC
		LIMIT ?
	`

	rows, err := s.db.Query(query, userID, userID, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var t Transaction
		err := rows.Scan(
			&t.ID, &t.FromUserID, &t.ToUserID, &t.Amount, &t.TransactionType,
			&t.Description, &t.Status, &t.CreatedAt, &t.FromUsername, &t.ToUsername,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}

// CreateAdminTransaction creates an administrative transaction for balance adjustment
func (s *BankingService) CreateAdminTransaction(userID int, amount float64, description, merchantName string) (*Transaction, error) {
	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Get PokéBank user ID
	var pokeBankID int
	err = tx.QueryRow(`SELECT id FROM users WHERE username = 'PokéBank'`).Scan(&pokeBankID)
	if err != nil {
		return nil, err
	}

	// Get current balance to check if withdrawal is possible
	if amount < 0 {
		var currentBalance float64
		err = tx.QueryRow(`SELECT balance FROM users WHERE id = ?`, userID).Scan(&currentBalance)
		if err != nil {
			return nil, err
		}

		if currentBalance+amount < 0 {
			return nil, fmt.Errorf("insufficient balance for withdrawal")
		}
	}

	// Update user balance
	_, err = tx.Exec(`UPDATE users SET balance = balance + ? WHERE id = ?`, amount, userID)
	if err != nil {
		return nil, err
	}

	// Create transaction record
	var fromUserID, toUserID int
	if amount > 0 {
		// Credit: money coming from PokéBank to user
		fromUserID = pokeBankID
		toUserID = userID
	} else {
		// Debit: money going from user to PokéBank
		fromUserID = userID
		toUserID = pokeBankID
	}

	transactionType := "admin_adjustment"
	status := "completed"

	query := `
		INSERT INTO transactions (from_user_id, to_user_id, amount, transaction_type, description, status, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	result, err := tx.Exec(query, fromUserID, toUserID, amount, transactionType, description, status, time.Now())
	if err != nil {
		return nil, err
	}

	transactionID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Reset PokéBank balance since it's involved in the transaction
	err = s.resetPokeBankBalanceInTx(tx, fromUserID, toUserID)
	if err != nil {
		return nil, err
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// Get the created transaction with user details
	return s.GetTransactionByID(int(transactionID))
}

// CreateMerchantTransaction creates a transaction to/from a virtual merchant
func (s *BankingService) CreateMerchantTransaction(userID int, amount float64, description, merchantName string) (*Transaction, error) {
	return s.CreateAdminTransaction(userID, amount, description, merchantName)
}

// EnsurePokeBankBalance ensures PokéBank always has the fixed balance
func (s *BankingService) EnsurePokeBankBalance() error {
	const pokeBankBalance = 999999999.99
	
	// Find PokéBank user
	var pokeBankID int
	err := s.db.QueryRow(`SELECT id FROM users WHERE username = 'PokéBank'`).Scan(&pokeBankID)
	if err != nil {
		return err // PokéBank doesn't exist
	}

	// Update balance to fixed amount
	_, err = s.db.Exec(`UPDATE users SET balance = ? WHERE id = ?`, pokeBankBalance, pokeBankID)
	return err
}

// GetUserPaymentRequests gets payment requests for a user
func (s *BankingService) GetUserPaymentRequests(userID int) ([]PaymentRequest, []PaymentRequest, error) {
	// Get incoming requests (where user is the recipient)
	incomingQuery := `
		SELECT pr.id, pr.from_user_id, pr.to_user_id, pr.amount, pr.reason, 
		       pr.message, pr.status, pr.created_at, u.username as from_username
		FROM payment_requests pr
		JOIN users u ON pr.from_user_id = u.id
		WHERE pr.to_user_id = ?
		ORDER BY pr.created_at DESC
	`

	rows, err := s.db.Query(incomingQuery, userID)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var incoming []PaymentRequest
	for rows.Next() {
		var pr PaymentRequest
		err := rows.Scan(
			&pr.ID, &pr.FromUserID, &pr.ToUserID, &pr.Amount, &pr.Reason,
			&pr.Message, &pr.Status, &pr.CreatedAt, &pr.FromUsername,
		)
		if err != nil {
			return nil, nil, err
		}
		incoming = append(incoming, pr)
	}

	// Get outgoing requests (where user is the sender)
	outgoingQuery := `
		SELECT pr.id, pr.from_user_id, pr.to_user_id, pr.amount, pr.reason, 
		       pr.message, pr.status, pr.created_at, u.username as to_username
		FROM payment_requests pr
		JOIN users u ON pr.to_user_id = u.id
		WHERE pr.from_user_id = ?
		ORDER BY pr.created_at DESC
	`

	rows, err = s.db.Query(outgoingQuery, userID)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var outgoing []PaymentRequest
	for rows.Next() {
		var pr PaymentRequest
		err := rows.Scan(
			&pr.ID, &pr.FromUserID, &pr.ToUserID, &pr.Amount, &pr.Reason,
			&pr.Message, &pr.Status, &pr.CreatedAt, &pr.ToUsername,
		)
		if err != nil {
			return nil, nil, err
		}
		outgoing = append(outgoing, pr)
	}

	return incoming, outgoing, nil
}

// ApprovePaymentRequest approves a payment request and processes the transfer
func (s *BankingService) ApprovePaymentRequest(requestID, userID int) error {
	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get payment request details
	var pr PaymentRequest
	err = tx.QueryRow(`
		SELECT id, from_user_id, to_user_id, amount, reason, status
		FROM payment_requests 
		WHERE id = ? AND to_user_id = ? AND status = 'pending'
	`, requestID, userID).Scan(&pr.ID, &pr.FromUserID, &pr.ToUserID, &pr.Amount, &pr.Reason, &pr.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("payment request not found or already processed")
		}
		return err
	}

	// Check user has sufficient balance
	var balance float64
	err = tx.QueryRow(`SELECT balance FROM users WHERE id = ?`, userID).Scan(&balance)
	if err != nil {
		return err
	}

	if balance < pr.Amount {
		return fmt.Errorf("insufficient balance")
	}

	// Process transfer (userID pays to fromUserID)
	_, err = tx.Exec(`UPDATE users SET balance = balance - ? WHERE id = ?`, pr.Amount, userID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`UPDATE users SET balance = balance + ? WHERE id = ?`, pr.Amount, pr.FromUserID)
	if err != nil {
		return err
	}

	// Create transaction record
	_, err = tx.Exec(`
		INSERT INTO transactions (from_user_id, to_user_id, amount, transaction_type, description, status, created_at)
		VALUES (?, ?, ?, 'transfer', ?, 'completed', ?)
	`, userID, pr.FromUserID, pr.Amount, "Payment for: "+pr.Reason, time.Now())
	if err != nil {
		return err
	}

	// Update payment request status
	_, err = tx.Exec(`UPDATE payment_requests SET status = 'approved' WHERE id = ?`, requestID)
	if err != nil {
		return err
	}

	// Reset PokéBank balance if involved in transaction
	err = s.resetPokeBankBalanceInTx(tx, userID, pr.FromUserID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// RejectPaymentRequest rejects a payment request
func (s *BankingService) RejectPaymentRequest(requestID, userID int) error {
	query := `
		UPDATE payment_requests 
		SET status = 'rejected' 
		WHERE id = ? AND to_user_id = ? AND status = 'pending'
	`
	result, err := s.db.Exec(query, requestID, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("payment request not found or already processed")
	}

	return nil
}

// CancelPaymentRequest cancels a payment request (for the requester)
func (s *BankingService) CancelPaymentRequest(requestID, userID int) error {
	query := `
		UPDATE payment_requests 
		SET status = 'cancelled' 
		WHERE id = ? AND from_user_id = ? AND status = 'pending'
	`
	result, err := s.db.Exec(query, requestID, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("payment request not found or already processed")
	}

	return nil
}

// GetTransactionByID gets a transaction by ID
func (s *BankingService) GetTransactionByID(id int) (*Transaction, error) {
	query := `
		SELECT t.id, t.from_user_id, t.to_user_id, t.amount, t.transaction_type, 
		       t.description, t.status, t.created_at,
		       u1.username as from_username, u2.username as to_username
		FROM transactions t
		LEFT JOIN users u1 ON t.from_user_id = u1.id
		LEFT JOIN users u2 ON t.to_user_id = u2.id
		WHERE t.id = ?
	`

	transaction := &Transaction{}
	err := s.db.QueryRow(query, id).Scan(
		&transaction.ID, &transaction.FromUserID, &transaction.ToUserID, &transaction.Amount,
		&transaction.TransactionType, &transaction.Description, &transaction.Status, &transaction.CreatedAt,
		&transaction.FromUsername, &transaction.ToUsername,
	)

	if err != nil {
		return nil, err
	}

	return transaction, nil
}

// getUsernameByID helper function to get username by user ID
func (s *BankingService) getUsernameByID(userID int) (string, error) {
	var username string
	err := s.db.QueryRow(`SELECT username FROM users WHERE id = ?`, userID).Scan(&username)
	return username, err
}

// resetPokeBankBalanceInTx resets PokéBank account balance to 999999999.99 if involved in transaction
func (s *BankingService) resetPokeBankBalanceInTx(tx *sql.Tx, fromUserID, toUserID int) error {
	const pokeBankBalance = 999999999.99
	
	// Check if sender is PokéBank and reset balance
	var senderUsername string
	err := tx.QueryRow(`SELECT username FROM users WHERE id = ?`, fromUserID).Scan(&senderUsername)
	if err == nil && senderUsername == "PokéBank" {
		_, err = tx.Exec(`UPDATE users SET balance = ? WHERE id = ?`, pokeBankBalance, fromUserID)
		if err != nil {
			return err
		}
	}
	
	// Check if recipient is PokéBank and reset balance
	var recipientUsername string
	err = tx.QueryRow(`SELECT username FROM users WHERE id = ?`, toUserID).Scan(&recipientUsername)
	if err == nil && recipientUsername == "PokéBank" {
		_, err = tx.Exec(`UPDATE users SET balance = ? WHERE id = ?`, pokeBankBalance, toUserID)
		if err != nil {
			return err
		}
	}
	
	return nil
}
