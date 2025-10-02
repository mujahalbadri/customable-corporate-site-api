package main

import (
	"customable-corporate-site-api/internal/config"
	"customable-corporate-site-api/internal/database"
	"customable-corporate-site-api/internal/handlers"
	"customable-corporate-site-api/internal/middleware"
	"customable-corporate-site-api/internal/repositories/postgres"
	"customable-corporate-site-api/internal/services"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configurations
	config := config.Load()

	// Set up database connection
	db, err := database.ConnectDB(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate database schemas
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to auto-migrate database: %v", err)
	}

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, config.JWT.Secret, config.JWT.ExpiresIn)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)

	// Set up Gin router
	router := setupRouter(authHandler, config.JWT.Secret)

	// Start the server
	log.Printf("Starting server on port %s...", config.Server.Port)
	log.Printf("API Health Check Endpoint: http://localhost:%s/api/v1/health", config.Server.Port)
	log.Fatal(router.Run(":" + config.Server.Port))
}

func setupRouter(authHandler *handlers.AuthHandler, jwtSecret string) *gin.Engine {
	// Create a Gin router
	router := gin.Default()

	// Global middleware
	router.Use(middleware.CORS())
	router.Use(middleware.Logger())

	// API v1 routes
	api := router.Group("/api/v1")

	// Public auth routes
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
	}

	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.JWTAuth(jwtSecret))
	{
		protected.GET("/auth/profile", authHandler.GetProfile)
		protected.PUT("/auth/profile", authHandler.UpdateProfile)
	}

	// Health check endpoint
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"message": "API is healthy",
			"version": "1.0.0",
		})
	})

	return router
}
