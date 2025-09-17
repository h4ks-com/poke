package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// WebhookService handles sending webhook notifications
type WebhookService struct {
	webhookURL string
	client     *http.Client
}

// NewWebhookService creates a new WebhookService
func NewWebhookService() *WebhookService {
	return &WebhookService{
		webhookURL: os.Getenv("WEBHOOK_URL"),
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// WebhookPayload represents the structure of webhook notifications
type WebhookPayload struct {
	Event     string      `json:"event"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data"`
}

// TransferWebhookData represents transfer webhook data
type TransferWebhookData struct {
	TransactionID   int     `json:"transactionId"`
	FromUserID      int     `json:"fromUserId"`
	FromUsername    string  `json:"fromUsername"`
	ToUserID        int     `json:"toUserId"`
	ToUsername      string  `json:"toUsername"`
	Amount          float64 `json:"amount"`
	Description     string  `json:"description"`
	TransactionType string  `json:"transactionType"`
	Status          string  `json:"status"`
}

// PaymentRequestWebhookData represents payment request webhook data
type PaymentRequestWebhookData struct {
	RequestID    int     `json:"requestId"`
	FromUserID   int     `json:"fromUserId"`
	FromUsername string  `json:"fromUsername"`
	ToUserID     int     `json:"toUserId"`
	ToUsername   string  `json:"toUsername"`
	Amount       float64 `json:"amount"`
	Reason       string  `json:"reason"`
	Message      string  `json:"message"`
	Status       string  `json:"status"`
}

// UserAuthWebhookData represents user authentication webhook data
type UserAuthWebhookData struct {
	UserID   int    `json:"userId"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Action   string `json:"action"` // "login", "register", "password_change"
}

// PaymentRequestActionWebhookData represents payment request action webhook data
type PaymentRequestActionWebhookData struct {
	RequestID    int    `json:"requestId"`
	UserID       int    `json:"userId"`
	Username     string `json:"username"`
	Action       string `json:"action"` // "approve", "reject"
	ActionUserID int    `json:"actionUserId"`
}

// CardRefreshWebhookData represents card refresh webhook data
type CardRefreshWebhookData struct {
	UserID     int    `json:"userId"`
	Username   string `json:"username"`
	CardNumber string `json:"cardNumber"`
	Action     string `json:"action"` // "refresh"
}

// SendTransferWebhook sends a webhook notification for transfers
func (w *WebhookService) SendTransferWebhook(transaction *Transaction) {
	if w.webhookURL == "" {
		return // No webhook URL configured
	}

	data := TransferWebhookData{
		TransactionID:   transaction.ID,
		FromUserID:      transaction.FromUserID,
		FromUsername:    transaction.FromUsername,
		ToUserID:        transaction.ToUserID,
		ToUsername:      transaction.ToUsername,
		Amount:          transaction.Amount,
		Description:     transaction.Description,
		TransactionType: transaction.TransactionType,
		Status:          transaction.Status,
	}

	payload := WebhookPayload{
		Event:     "transfer_completed",
		Timestamp: time.Now(),
		Data:      data,
	}

	w.sendWebhook(payload)
}

// SendPaymentRequestWebhook sends a webhook notification for payment requests
func (w *WebhookService) SendPaymentRequestWebhook(paymentRequest *PaymentRequest) {
	if w.webhookURL == "" {
		return
	}

	data := PaymentRequestWebhookData{
		RequestID:    paymentRequest.ID,
		FromUserID:   paymentRequest.FromUserID,
		FromUsername: paymentRequest.FromUsername,
		ToUserID:     paymentRequest.ToUserID,
		ToUsername:   paymentRequest.ToUsername,
		Amount:       paymentRequest.Amount,
		Reason:       paymentRequest.Reason,
		Message:      paymentRequest.Message,
		Status:       paymentRequest.Status,
	}

	payload := WebhookPayload{
		Event:     "payment_request_created",
		Timestamp: time.Now(),
		Data:      data,
	}

	w.sendWebhook(payload)
}

// SendUserAuthWebhook sends a webhook notification for user authentication events
func (w *WebhookService) SendUserAuthWebhook(userID int, username, email, action string) {
	if w.webhookURL == "" {
		return
	}

	data := UserAuthWebhookData{
		UserID:   userID,
		Username: username,
		Email:    email,
		Action:   action,
	}

	payload := WebhookPayload{
		Event:     "user_auth",
		Timestamp: time.Now(),
		Data:      data,
	}

	w.sendWebhook(payload)
}

// SendPaymentRequestApprovalWebhook sends a webhook notification when a payment request is approved
func (w *WebhookService) SendPaymentRequestApprovalWebhook(requestID, actionUserID int) {
	if w.webhookURL == "" {
		return
	}

	data := PaymentRequestActionWebhookData{
		RequestID:    requestID,
		Action:       "approve",
		ActionUserID: actionUserID,
	}

	payload := WebhookPayload{
		Event:     "payment_request_approved",
		Timestamp: time.Now(),
		Data:      data,
	}

	w.sendWebhook(payload)
}

// SendPaymentRequestRejectionWebhook sends a webhook notification when a payment request is rejected
func (w *WebhookService) SendPaymentRequestRejectionWebhook(requestID, actionUserID int) {
	if w.webhookURL == "" {
		return
	}

	data := PaymentRequestActionWebhookData{
		RequestID:    requestID,
		Action:       "reject",
		ActionUserID: actionUserID,
	}

	payload := WebhookPayload{
		Event:     "payment_request_rejected",
		Timestamp: time.Now(),
		Data:      data,
	}

	w.sendWebhook(payload)
}

// SendCardRefreshNotification sends a webhook notification for card refresh
func (w *WebhookService) SendCardRefreshNotification(username, cardNumber string) {
	if w.webhookURL == "" {
		return // No webhook URL configured
	}

	data := CardRefreshWebhookData{
		Username:   username,
		CardNumber: cardNumber,
		Action:     "refresh",
	}

	payload := WebhookPayload{
		Event:     "card_refreshed",
		Timestamp: time.Now(),
		Data:      data,
	}

	w.sendWebhook(payload)
}

// sendWebhook sends a webhook payload to the configured URL
func (w *WebhookService) sendWebhook(payload WebhookPayload) {
	// Skip webhook if URL is not configured or is a placeholder
	if w.webhookURL == "" || w.webhookURL == "https://your-webhook-endpoint.com/webhook" {
		// Silently skip webhook - no error logging for development
		return
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error marshaling webhook payload: %v\n", err)
		return
	}

	req, err := http.NewRequest("POST", w.webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error creating webhook request: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Viridian-Bank-Webhook/1.0")

	resp, err := w.client.Do(req)
	if err != nil {
		fmt.Printf("Error sending webhook: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Printf("Webhook returned non-success status: %d\n", resp.StatusCode)
	}
}

// AdminTransactionWebhookData represents admin transaction webhook data
type AdminTransactionWebhookData struct {
	TransactionID int     `json:"transactionId"`
	UserID        int     `json:"userId"`
	Username      string  `json:"username"`
	Amount        float64 `json:"amount"`
	Description   string  `json:"description"`
	MerchantName  string  `json:"merchantName"`
	ActionType    string  `json:"actionType"` // "credit" or "debit"
}

// MerchantTransactionWebhookData represents merchant transaction webhook data
type MerchantTransactionWebhookData struct {
	TransactionID int     `json:"transactionId"`
	UserID        int     `json:"userId"`
	Username      string  `json:"username"`
	Amount        float64 `json:"amount"`
	Description   string  `json:"description"`
	MerchantName  string  `json:"merchantName"`
}

// SendAdminTransactionWebhook sends a webhook notification for admin transactions
func (w *WebhookService) SendAdminTransactionWebhook(transactionID, userID int, username string, amount float64, description, merchantName string) {
	if w.webhookURL == "" {
		return // No webhook URL configured
	}

	actionType := "credit"
	if amount < 0 {
		actionType = "debit"
	}

	data := AdminTransactionWebhookData{
		TransactionID: transactionID,
		UserID:        userID,
		Username:      username,
		Amount:        amount,
		Description:   description,
		MerchantName:  merchantName,
		ActionType:    actionType,
	}

	payload := WebhookPayload{
		Event:     "admin_transaction",
		Timestamp: time.Now(),
		Data:      data,
	}

	w.sendWebhook(payload)
}

// SendMerchantTransactionWebhook sends a webhook notification for merchant transactions
func (w *WebhookService) SendMerchantTransactionWebhook(transactionID, userID int, username string, amount float64, description, merchantName string) {
	if w.webhookURL == "" {
		return // No webhook URL configured
	}

	data := MerchantTransactionWebhookData{
		TransactionID: transactionID,
		UserID:        userID,
		Username:      username,
		Amount:        amount,
		Description:   description,
		MerchantName:  merchantName,
	}

	payload := WebhookPayload{
		Event:     "merchant_transaction",
		Timestamp: time.Now(),
		Data:      data,
	}

	w.sendWebhook(payload)
}
