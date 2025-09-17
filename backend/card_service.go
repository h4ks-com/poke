package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

// CardService handles card-related operations
type CardService struct {
	db *sql.DB
}

// NewCardService creates a new card service
func NewCardService(db *sql.DB) *CardService {
	return &CardService{db: db}
}

// GetUserCard retrieves the active card for a user, creating one if it doesn't exist
func (cs *CardService) GetUserCard(userID int, accountNumber string) (*Card, error) {
	// First try to get existing active card
	card, err := cs.getActiveCard(userID)
	if err == nil {
		return card, nil
	}
	
	// If no active card exists, create a new one
	return cs.createCard(userID, accountNumber)
}

// getActiveCard retrieves the user's active card
func (cs *CardService) getActiveCard(userID int) (*Card, error) {
	query := `
		SELECT id, user_id, card_number, expiry_date, refresh_seed, 
		       last_refresh_date, is_active, created_at, updated_at
		FROM cards 
		WHERE user_id = ? AND is_active = TRUE
		ORDER BY created_at DESC
		LIMIT 1
	`
	
	var card Card
	var lastRefreshDate sql.NullTime
	
	err := cs.db.QueryRow(query, userID).Scan(
		&card.ID, &card.UserID, &card.CardNumber, &card.ExpiryDate,
		&card.RefreshSeed, &lastRefreshDate, &card.IsActive,
		&card.CreatedAt, &card.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	if lastRefreshDate.Valid {
		card.LastRefreshDate = &lastRefreshDate.Time
	}
	
	return &card, nil
}

// createCard creates a new card for a user
func (cs *CardService) createCard(userID int, accountNumber string) (*Card, error) {
	cardNumber := cs.generateCardNumber(accountNumber, 0)
	expiryDate := cs.generateExpiryDate()
	
	query := `
		INSERT INTO cards (user_id, card_number, expiry_date, refresh_seed, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`
	
	result, err := cs.db.Exec(query, userID, cardNumber, expiryDate, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to create card: %v", err)
	}
	
	cardID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get card ID: %v", err)
	}
	
	return &Card{
		ID:              int(cardID),
		UserID:          userID,
		CardNumber:      cardNumber,
		ExpiryDate:      expiryDate,
		RefreshSeed:     0,
		LastRefreshDate: nil,
		IsActive:        true,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}, nil
}

// RefreshCard generates a new card for the user if allowed
func (cs *CardService) RefreshCard(userID int, accountNumber string) (*Card, error) {
	// Get current active card
	currentCard, err := cs.getActiveCard(userID)
	if err != nil {
		return nil, fmt.Errorf("no active card found: %v", err)
	}
	
	// Check if refresh is allowed (once per day)
	if !cs.canRefreshCard(currentCard) {
		return nil, fmt.Errorf("card can only be refreshed once per day")
	}
	
	// Start transaction
	tx, err := cs.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()
	
	// Deactivate current card
	_, err = tx.Exec("UPDATE cards SET is_active = FALSE, updated_at = CURRENT_TIMESTAMP WHERE id = ?", currentCard.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to deactivate current card: %v", err)
	}
	
	// Create new card with incremented refresh seed
	newRefreshSeed := currentCard.RefreshSeed + 1
	newCardNumber := cs.generateCardNumber(accountNumber, newRefreshSeed)
	newExpiryDate := cs.generateExpiryDate()
	now := time.Now()
	
	query := `
		INSERT INTO cards (user_id, card_number, expiry_date, refresh_seed, last_refresh_date, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`
	
	result, err := tx.Exec(query, userID, newCardNumber, newExpiryDate, newRefreshSeed, now)
	if err != nil {
		return nil, fmt.Errorf("failed to create new card: %v", err)
	}
	
	newCardID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get new card ID: %v", err)
	}
	
	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}
	
	return &Card{
		ID:              int(newCardID),
		UserID:          userID,
		CardNumber:      newCardNumber,
		ExpiryDate:      newExpiryDate,
		RefreshSeed:     newRefreshSeed,
		LastRefreshDate: &now,
		IsActive:        true,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}, nil
}

// canRefreshCard checks if a card can be refreshed (once per day limit)
func (cs *CardService) canRefreshCard(card *Card) bool {
	if card.LastRefreshDate == nil {
		return true
	}
	
	// Check if 24 hours have passed since last refresh
	now := time.Now()
	timeSinceRefresh := now.Sub(*card.LastRefreshDate)
	return timeSinceRefresh >= 24*time.Hour
}

// generateCardNumber generates a card number using the same algorithm as frontend
func (cs *CardService) generateCardNumber(accountNumber string, refreshSeed int) string {
	// Extract numeric part from account number
	accountSeed := cs.extractNumeric(accountNumber)
	if accountSeed == 0 {
		accountSeed = 1234
	}
	
	// Add refresh seed to generate new numbers when card is refreshed
	combinedSeed := accountSeed + refreshSeed
	
	// Viridian City Bank BIN (Bank Identification Number) - using 4532 (Visa format)
	cardNumber := "4532"
	
	// Generate middle digits based on combined seed
	middle8Digits := cs.generateMiddleDigits(combinedSeed)
	cardNumber += middle8Digits
	
	// Add check digit using Luhn algorithm
	checkDigit := cs.calculateLuhnCheckDigit(cardNumber)
	cardNumber += strconv.Itoa(checkDigit)
	
	return cardNumber
}

// extractNumeric extracts numeric characters from a string and converts to int
func (cs *CardService) extractNumeric(s string) int {
	var numeric string
	for _, char := range s {
		if char >= '0' && char <= '9' {
			numeric += string(char)
		}
	}
	
	if numeric == "" {
		return 0
	}
	
	result, err := strconv.Atoi(numeric)
	if err != nil {
		return 0
	}
	
	return result
}

// generateMiddleDigits generates 8 middle digits using a pseudo-random generator
func (cs *CardService) generateMiddleDigits(seed int) string {
	random := seed
	result := ""
	
	for i := 0; i < 8; i++ {
		random = (random*9301 + 49297) % 233280
		digit := (random * 10) / 233280
		result += strconv.Itoa(digit)
	}
	
	// Ensure it's exactly 8 digits
	for len(result) < 8 {
		result = "0" + result
	}
	
	return result[:8]
}

// calculateLuhnCheckDigit calculates the Luhn check digit for a card number
func (cs *CardService) calculateLuhnCheckDigit(cardNumber string) int {
	sum := 0
	isEven := true
	
	// Process digits from right to left
	for i := len(cardNumber) - 1; i >= 0; i-- {
		digit := int(cardNumber[i] - '0')
		
		if isEven {
			digit *= 2
			if digit > 9 {
				digit = digit - 9
			}
		}
		
		sum += digit
		isEven = !isEven
	}
	
	return (10 - (sum % 10)) % 10
}

// generateExpiryDate generates an expiry date 3 years from now
func (cs *CardService) generateExpiryDate() string {
	now := time.Now()
	expiryYear := now.Year() + 3
	expiryMonth := now.Month()
	
	return fmt.Sprintf("%02d/%02d", int(expiryMonth), expiryYear%100)
}

// GetTimeUntilNextRefresh returns the time remaining until next refresh is allowed
func (cs *CardService) GetTimeUntilNextRefresh(card *Card) *time.Duration {
	if card.LastRefreshDate == nil {
		return nil
	}
	
	now := time.Now()
	nextRefreshTime := card.LastRefreshDate.Add(24 * time.Hour)
	
	if now.After(nextRefreshTime) || now.Equal(nextRefreshTime) {
		return nil
	}
	
	duration := nextRefreshTime.Sub(now)
	return &duration
}
