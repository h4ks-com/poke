package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// BankingHandler handles banking-related HTTP requests
type BankingHandler struct {
	service        *BankingService
	userService    *UserService
	webhookService *WebhookService
	cardService    *CardService
}

// NewBankingHandler creates a new BankingHandler
func NewBankingHandler(service *BankingService, userService *UserService, webhookService *WebhookService, cardService *CardService) *BankingHandler {
	return &BankingHandler{
		service:        service,
		userService:    userService,
		webhookService: webhookService,
		cardService:    cardService,
	}
}

// GetBalanceHandler handles GET /api/balance
func (h *BankingHandler) GetBalanceHandler(c *gin.Context) {
	userID := c.GetInt("userID")

	balance, err := h.service.GetUserBalance(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get balance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": balance})
}

// TransferRequest represents the request body for transfers
type TransferRequest struct {
	To          string  `json:"to" binding:"required"`          // Can be username or account number
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Description string  `json:"description"`
}

// TransferHandler handles POST /api/transfer
func (h *BankingHandler) TransferHandler(c *gin.Context) {
	userID := c.GetInt("userID")

	var req TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Resolve the target user by username or account number
	targetUser, err := h.userService.GetUserByUsernameOrAccountNumber(req.To)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipient not found"})
		return
	}

	transaction, err := h.service.Transfer(userID, targetUser.AccountNumber, req.Amount, req.Description)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Send webhook notification
	go h.webhookService.SendTransferWebhook(transaction)

	c.JSON(http.StatusOK, gin.H{
		"message":     "Transfer completed successfully",
		"transaction": transaction,
	})
}

// PaymentRequestRequest represents the request body for payment requests
type PaymentRequestRequest struct {
	To      string  `json:"to" binding:"required"`      // Can be username or account number
	Amount  float64 `json:"amount" binding:"required,gt=0"`
	Reason  string  `json:"reason" binding:"required"`
	Message string  `json:"message"`
}

// CreatePaymentRequestHandler handles POST /api/payment-requests
func (h *BankingHandler) CreatePaymentRequestHandler(c *gin.Context) {
	userID := c.GetInt("userID")

	var req PaymentRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Resolve the target user by username or account number
	targetUser, err := h.userService.GetUserByUsernameOrAccountNumber(req.To)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipient not found"})
		return
	}

	paymentRequest, err := h.service.CreatePaymentRequest(userID, targetUser.AccountNumber, req.Amount, req.Reason, req.Message)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Send webhook notification
	go h.webhookService.SendPaymentRequestWebhook(paymentRequest)

	c.JSON(http.StatusCreated, gin.H{
		"message":        "Payment request created successfully",
		"paymentRequest": paymentRequest,
	})
}

// GetTransactionsHandler handles GET /api/transactions
func (h *BankingHandler) GetTransactionsHandler(c *gin.Context) {
	userID := c.GetInt("userID")

	// Parse limit parameter (default 50)
	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100 // Max limit
	}

	transactions, err := h.service.GetUserTransactions(userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get transactions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transactions": transactions})
}

// GetPaymentRequestsHandler handles GET /api/payment-requests
func (h *BankingHandler) GetPaymentRequestsHandler(c *gin.Context) {
	userID := c.GetInt("userID")

	incoming, outgoing, err := h.service.GetUserPaymentRequests(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get payment requests"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"incoming": incoming,
		"outgoing": outgoing,
	})
}

// PaymentRequestActionRequest represents the request body for payment request actions
type PaymentRequestActionRequest struct {
	Action string `json:"action" binding:"required,oneof=approve reject cancel"`
}

// HandlePaymentRequestHandler handles PUT /api/payment-requests/:id
func (h *BankingHandler) HandlePaymentRequestHandler(c *gin.Context) {
	userID := c.GetInt("userID")

	requestIDStr := c.Param("id")
	requestID, err := strconv.Atoi(requestIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	var req PaymentRequestActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	switch req.Action {
	case "approve":
		err = h.service.ApprovePaymentRequest(requestID, userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Send webhook notification for approval
		go h.webhookService.SendPaymentRequestApprovalWebhook(requestID, userID)

		c.JSON(http.StatusOK, gin.H{"message": "Payment request approved and transfer completed"})

	case "reject":
		err = h.service.RejectPaymentRequest(requestID, userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Send webhook notification for rejection
		go h.webhookService.SendPaymentRequestRejectionWebhook(requestID, userID)

		c.JSON(http.StatusOK, gin.H{"message": "Payment request rejected"})

	case "cancel":
		err = h.service.CancelPaymentRequest(requestID, userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Send webhook notification for cancellation
		go h.webhookService.SendPaymentRequestRejectionWebhook(requestID, userID)

		c.JSON(http.StatusOK, gin.H{"message": "Payment request cancelled"})

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action"})
	}
}

// GetAccountInfoHandler handles GET /api/account
func (h *BankingHandler) GetAccountInfoHandler(c *gin.Context) {
	userID := c.GetInt("userID")

	user, err := h.getUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get account info"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":            user.ID,
		"username":      user.Username,
		"email":         user.Email,
		"accountNumber": user.AccountNumber,
		"balance":       user.Balance,
		"createdAt":     user.CreatedAt,
	})
}

// getUserByID helper function to get user by ID
func (h *BankingHandler) getUserByID(userID int) (*User, error) {
	user := &User{}
	query := `
		SELECT id, username, email, account_number, balance, created_at
		FROM users WHERE id = $1
	`
	err := h.service.db.QueryRow(query, userID).Scan(
		&user.ID, &user.Username, &user.Email, &user.AccountNumber, &user.Balance, &user.CreatedAt,
	)
	return user, err
}

// GetCardHandler handles GET /api/card
func (h *BankingHandler) GetCardHandler(c *gin.Context) {
	userID := c.GetInt("userID")

	// Get user info to access account number
	user, err := h.getUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}

	// Get or create user card
	card, err := h.cardService.GetUserCard(userID, user.AccountNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get card"})
		return
	}

	// Calculate time until next refresh
	var timeUntilRefresh map[string]interface{}
	if duration := h.cardService.GetTimeUntilNextRefresh(card); duration != nil {
		hours := int(duration.Hours())
		minutes := int(duration.Minutes()) % 60
		timeUntilRefresh = map[string]interface{}{
			"hours":   hours,
			"minutes": minutes,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"card":             card,
		"canRefresh":       h.cardService.GetTimeUntilNextRefresh(card) == nil,
		"timeUntilRefresh": timeUntilRefresh,
	})
}

// RefreshCardHandler handles POST /api/card/refresh
func (h *BankingHandler) RefreshCardHandler(c *gin.Context) {
	userID := c.GetInt("userID")

	// Get user info to access account number
	user, err := h.getUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}

	// Refresh the card
	newCard, err := h.cardService.RefreshCard(userID, user.AccountNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Send webhook notification for card refresh
	go h.webhookService.SendCardRefreshNotification(user.Username, newCard.CardNumber)

	c.JSON(http.StatusOK, gin.H{
		"message": "Card refreshed successfully",
		"card":    newCard,
		"canRefresh": false,
	})
}
