package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	userService    *UserService
	webhookService *WebhookService
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(userService *UserService, webhookService *WebhookService) *AuthHandler {
	return &AuthHandler{
		userService:    userService,
		webhookService: webhookService,
	}
}

// LoginRequest represents a login request
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
	User    *User  `json:"user,omitempty"`
}

// RegisterRequest represents a registration request
type RegisterRequest struct {
	Username        string `json:"username" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirmPassword" binding:"required"`
}

// ChangePasswordRequest represents a password change request
type ChangePasswordRequest struct {
	CurrentPassword    string `json:"currentPassword" binding:"required"`
	NewPassword        string `json:"newPassword" binding:"required"`
	ConfirmNewPassword string `json:"confirmNewPassword" binding:"required"`
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	// Get user by username
	user, err := h.userService.GetUserByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, LoginResponse{
			Success: false,
			Message: "Invalid username or password",
		})
		return
	}

	// Verify password
	var passwordValid bool
	
	// Special case for PokéBank account - use admin key as password
	if user.Username == "PokéBank" {
		adminKey := getEnv("ADMIN_KEY", "")
		if adminKey == "" {
			c.JSON(http.StatusServiceUnavailable, LoginResponse{
				Success: false,
				Message: "PokéBank login not available",
			})
			return
		}
		
		passwordValid = (req.Password == adminKey)
	} else {
		// Regular user authentication
		passwordValid = h.userService.VerifyPassword(user, req.Password)
	}
	
	if !passwordValid {
		c.JSON(http.StatusUnauthorized, LoginResponse{
			Success: false,
			Message: "Invalid username or password",
		})
		return
	}

	// Generate session token
	token, err := generateSessionToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate session token"})
		return
	}

	// Create session (expires in 24 hours)
	expiresAt := time.Now().Add(24 * time.Hour)
	if err := h.userService.CreateSession(user.ID, token, expiresAt); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	// Clear password hash from response
	user.PasswordHash = ""

	// Send webhook notification
	go h.webhookService.SendUserAuthWebhook(user.ID, user.Username, user.Email, "login")

	c.JSON(http.StatusOK, LoginResponse{
		Success: true,
		Message: "Login successful",
		Token:   token,
		User:    user,
	})
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	// Validate input
	if err := h.validateRegistration(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if username already exists
	if _, err := h.userService.GetUserByUsername(req.Username); err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	// Create user
	user, err := h.userService.CreateUser(req.Username, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "details": err.Error()})
		return
	}

	// Send webhook notification
	go h.webhookService.SendUserAuthWebhook(user.ID, user.Username, user.Email, "register")

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Account created successfully",
		"user":    user,
	})
}

// ChangePassword handles password changes
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	// Get current user from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get user details
	user, err := h.userService.GetUserByID(userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Verify current password
	var currentPasswordValid bool
	
	// Special case for PokéBank account - use admin key for current password verification
	if user.Username == "PokéBank" {
		adminKey := getEnv("ADMIN_KEY", "")
		if adminKey == "" {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "PokéBank password change not available"})
			return
		}
		currentPasswordValid = (req.CurrentPassword == adminKey)
	} else {
		// Regular user authentication
		currentPasswordValid = h.userService.VerifyPassword(user, req.CurrentPassword)
	}
	
	if !currentPasswordValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Current password is incorrect"})
		return
	}

	// Prevent password changes for PokéBank account
	if user.Username == "PokéBank" {
		c.JSON(http.StatusForbidden, gin.H{"error": "PokéBank password cannot be changed. It is tied to the admin key."})
		return
	}

	// Validate new password
	if err := h.validatePassword(req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check password confirmation
	if req.NewPassword != req.ConfirmNewPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "New passwords do not match"})
		return
	}

	// Check if new password is different from current
	if req.CurrentPassword == req.NewPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "New password must be different from current password"})
		return
	}

	// Update password
	if err := h.userService.UpdatePassword(user.ID, req.NewPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	// Invalidate all user sessions (force re-login)
	h.userService.DeleteAllUserSessions(user.ID)

	// Send webhook notification
	go h.webhookService.SendUserAuthWebhook(user.ID, user.Username, user.Email, "password_change")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Password changed successfully. Please log in again with your new password.",
		"username": user.Username,
	})
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// Get token from header
	token := extractTokenFromHeader(c.GetHeader("Authorization"))
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No token provided"})
		return
	}

	// Delete session
	h.userService.DeleteSession(token)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logged out successfully",
	})
}

// validateRegistration validates registration input
func (h *AuthHandler) validateRegistration(req *RegisterRequest) error {
	// Validate username
	if len(req.Username) < 3 || len(req.Username) > 20 {
		return fmt.Errorf("username must be between 3 and 20 characters")
	}

	// Restrict reserved usernames
	if req.Username == "PokéBank" || req.Username == "PokeBank" || req.Username == "POKEBANK" {
		return fmt.Errorf("this username is reserved and cannot be used")
	}

	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !usernameRegex.MatchString(req.Username) {
		return fmt.Errorf("username can only contain letters, numbers, and underscores")
	}

	// Validate email
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return fmt.Errorf("invalid email address format")
	}

	// Validate password
	if err := h.validatePassword(req.Password); err != nil {
		return err
	}

	// Check password confirmation
	if req.Password != req.ConfirmPassword {
		return fmt.Errorf("passwords do not match")
	}

	return nil
}

// validatePassword validates password strength
func (h *AuthHandler) validatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`\d`).MatchString(password)

	if !hasLetter || !hasNumber {
		return fmt.Errorf("password must contain at least one letter and one number")
	}

	return nil
}

// generateSessionToken generates a random session token
func generateSessionToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// extractTokenFromHeader extracts the token from Authorization header
func extractTokenFromHeader(authHeader string) string {
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}
