package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminHandler handles administrative operations
type AdminHandler struct {
	bankingService *BankingService
	userService    *UserService
	webhookService *WebhookService
}

// NewAdminHandler creates a new AdminHandler
func NewAdminHandler(bankingService *BankingService, userService *UserService, webhookService *WebhookService) *AdminHandler {
	return &AdminHandler{
		bankingService: bankingService,
		userService:    userService,
		webhookService: webhookService,
	}
}

// AdjustBalanceRequest represents a balance adjustment request
type AdjustBalanceRequest struct {
	UserID          int     `json:"user_id" binding:"required"`
	Amount          float64 `json:"amount" binding:"required"`
	Description     string  `json:"description" binding:"required"`
	MerchantName    string  `json:"merchant_name" binding:"required"`
	TransactionType string  `json:"transaction_type"` // "credit" or "debit" - optional, inferred from amount
}

// AdjustBalanceHandler handles administrative balance adjustments
func (h *AdminHandler) AdjustBalanceHandler(c *gin.Context) {
	var req AdjustBalanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	// Validate user exists
	user, err := h.userService.GetUserByID(req.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Determine transaction type if not provided
	if req.TransactionType == "" {
		if req.Amount > 0 {
			req.TransactionType = "credit"
		} else {
			req.TransactionType = "debit"
		}
	}

	// Create the administrative transaction
	transaction, err := h.bankingService.CreateAdminTransaction(req.UserID, req.Amount, req.Description, req.MerchantName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to adjust balance", "details": err.Error()})
		return
	}

	// Get updated balance
	newBalance, err := h.bankingService.GetUserBalance(req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get updated balance"})
		return
	}

	// Send webhook notification
	go h.webhookService.SendAdminTransactionWebhook(transaction.ID, req.UserID, user.Username, req.Amount, req.Description, req.MerchantName)

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"message":     "Balance adjusted successfully",
		"transaction": transaction,
		"new_balance": newBalance,
	})
}

// GetAllUsersHandler returns all users with their balances (admin only)
func (h *AdminHandler) GetAllUsersHandler(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"users":   users,
	})
}

// GetUserByAccountHandler gets user info by account number (admin only)
func (h *AdminHandler) GetUserByAccountHandler(c *gin.Context) {
	accountNumber := c.Param("account")
	if accountNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account number is required"})
		return
	}

	user, err := h.userService.GetUserByAccountNumber(accountNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user":    user,
	})
}

// CreateMerchantTransactionRequest represents a merchant transaction request
type CreateMerchantTransactionRequest struct {
	UserID      int     `json:"user_id" binding:"required"`
	Amount      float64 `json:"amount" binding:"required"`
	Description string  `json:"description" binding:"required"`
	MerchantName string `json:"merchant_name" binding:"required"`
}

// BankTransferRequest represents a transfer from PokéBank to a user
type BankTransferRequest struct {
	To          string  `json:"to" binding:"required"`        // Can be username or account number
	Amount      float64 `json:"amount" binding:"required"`
	Description string  `json:"description" binding:"required"`
}

// BankTransferHandler handles transfers from PokéBank to users
func (h *AdminHandler) BankTransferHandler(c *gin.Context) {
	var req BankTransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	// Validate amount is positive
	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be positive"})
		return
	}

	// Get PokéBank user
	pokeBankUser, err := h.userService.GetUserByUsername("PokéBank")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "PokéBank account not found"})
		return
	}

	// Validate recipient exists
	recipient, err := h.userService.GetUserByUsernameOrAccountNumber(req.To)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipient not found"})
		return
	}

	// Check PokéBank can't transfer to itself
	if pokeBankUser.AccountNumber == recipient.AccountNumber {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot transfer to PokéBank account"})
		return
	}

	// Perform the transfer using the existing transfer method
	transaction, err := h.bankingService.Transfer(pokeBankUser.ID, recipient.AccountNumber, req.Amount, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transfer failed", "details": err.Error()})
		return
	}

	// Send webhook notification
	go h.webhookService.SendTransferWebhook(transaction)

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"message":     "Transfer from PokéBank completed successfully",
		"transaction": transaction,
	})
}

// CreateMerchantTransactionHandler creates a transaction to/from a virtual merchant
func (h *AdminHandler) CreateMerchantTransactionHandler(c *gin.Context) {
	var req CreateMerchantTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	// Validate user exists
	user, err := h.userService.GetUserByID(req.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Create merchant transaction
	transaction, err := h.bankingService.CreateMerchantTransaction(req.UserID, req.Amount, req.Description, req.MerchantName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create merchant transaction", "details": err.Error()})
		return
	}

	// Get updated balance
	newBalance, err := h.bankingService.GetUserBalance(req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get updated balance"})
		return
	}

	// Send webhook notification
	go h.webhookService.SendMerchantTransactionWebhook(transaction.ID, req.UserID, user.Username, req.Amount, req.Description, req.MerchantName)

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"message":     "Merchant transaction created successfully",
		"transaction": transaction,
		"new_balance": newBalance,
	})
}
