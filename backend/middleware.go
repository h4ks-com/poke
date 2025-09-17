package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// authMiddleware validates the session token and sets user context
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		token := extractTokenFromHeader(authHeader)

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No valid authorization token"})
			c.Abort()
			return
		}

		// Initialize user service (in a real app, this would be injected)
		db, err := InitDB()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed"})
			c.Abort()
			return
		}
		defer db.Close()

		userService := NewUserService(db)

		// Validate session
		session, err := userService.GetSessionByToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired session"})
			c.Abort()
			return
		}

		// Get user details
		user, err := userService.GetUserByID(session.UserID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		// Set user context
		c.Set("userID", user.ID)
		c.Set("user", user)
		c.Set("sessionToken", token)

		c.Next()
	}
}
