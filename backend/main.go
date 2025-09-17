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
	// Load environment variables from .env file if available
	godotenv.Load() // Silently ignore errors - env file is optional

	// Initialize database
	db, err := InitDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize services
	userService := NewUserService(db)
	bankingService := NewBankingService(db)
	cardService := NewCardService(db)
	webhookService := NewWebhookService()

	// Initialize handlers
	authHandler := NewAuthHandler(userService, webhookService)
	bankingHandler := NewBankingHandler(bankingService, userService, webhookService, cardService)
	adminHandler := NewAdminHandler(bankingService, userService, webhookService)

	// Ensure PokéBank has fixed balance on startup
	bankingService.EnsurePokeBankBalance()

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

	// Serve static files (CSS, JS, assets)
	r.Static("/css", "./css")
	r.Static("/js", "./js")
	r.Static("/assets", "./assets")

	// Serve main HTML files
	r.StaticFile("/", "./index.html")
	r.StaticFile("/debug_test.html", "./debug_test.html")

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
			protected.GET("/card", bankingHandler.GetCardHandler)
			protected.POST("/card/refresh", bankingHandler.RefreshCardHandler)
		}

		// Admin routes (require admin authentication)
		admin := api.Group("/admin")
		admin.Use(adminAuthMiddleware())
		{
			admin.POST("/adjust-balance", adminHandler.AdjustBalanceHandler)
			admin.POST("/merchant-transaction", adminHandler.CreateMerchantTransactionHandler)
			admin.POST("/bank-transfer", adminHandler.BankTransferHandler)
			admin.GET("/users", adminHandler.GetAllUsersHandler)
			admin.GET("/user/:account", adminHandler.GetUserByAccountHandler)
		}
	}

	// Start server
	port := getEnv("PORT", "8080")
	log.Printf("Server starting on port %s", port)
	log.Fatal(r.Run(":" + port))
}
