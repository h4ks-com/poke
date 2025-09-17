package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize database
	db, err := InitDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize services
	userService := NewUserService(db)
	bankingService := NewBankingService(db)
	webhookService := NewWebhookService()

	// Initialize handlers
	authHandler := NewAuthHandler(userService, webhookService)
	bankingHandler := NewBankingHandler(bankingService, webhookService)

	// Initialize Gin router
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API routes
	api := r.Group("/api")
	{
		// Authentication routes
		api.POST("/register", authHandler.Register)
		api.POST("/login", authHandler.Login)
		api.POST("/change-password", authMiddleware(), authHandler.ChangePassword)

		// Protected banking routes
		protected := api.Group("/")
		protected.Use(authMiddleware())
		{
			protected.GET("/account", bankingHandler.GetAccountInfoHandler)
			protected.GET("/balance", bankingHandler.GetBalanceHandler)
			protected.POST("/transfer", bankingHandler.TransferHandler)
			protected.GET("/transactions", bankingHandler.GetTransactionsHandler)
			protected.POST("/payment-requests", bankingHandler.CreatePaymentRequestHandler)
			protected.GET("/payment-requests", bankingHandler.GetPaymentRequestsHandler)
			protected.PUT("/payment-requests/:id", bankingHandler.HandlePaymentRequestHandler)
		}
	}

	// Start server
	port := getEnv("PORT", "8080")
	log.Printf("Server starting on port %s", port)
	log.Fatal(r.Run(":" + port))
}
