package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/wayne-stays/backend/internal/accommodation"
	"github.com/wayne-stays/backend/internal/config"
	"github.com/wayne-stays/backend/internal/database"
	"github.com/wayne-stays/backend/internal/storage"
)

func main() {
	// Load configuration
	cfg := config.MustLoad()

	// Initialize storage
	store, err := storage.New(&cfg.Storage)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
	log.Printf("✅ Storage initialized (%s)", cfg.Storage.Type)

	// Connect to database
	db, err := database.Connect(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Seed database
	if err := database.Seed(db, store); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	// Initialize Gin
	gin.SetMode(cfg.Server.Mode)
	router := gin.Default()

	// CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.Server.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Serve static files (images) in development
	if cfg.Storage.Type == "local" {
		router.Static("/static", cfg.Storage.LocalPath)
		log.Printf("📁 Serving static files from: %s", cfg.Storage.LocalPath)
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "wayne-stays-backend",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Initialize accommodation components
		accommodationRepo := accommodation.NewRepository(db)
		accommodationService := accommodation.NewService(accommodationRepo, store)
		accommodationHandler := accommodation.NewHandler(accommodationService)

		// Register accommodation routes
		accommodationHandler.RegisterRoutes(v1)
	}

	// Start server
	addr := ":" + cfg.Server.Port
	log.Printf("🚀 Server starting on http://localhost%s", addr)
	log.Printf("📍 API available at http://localhost%s/api/v1", addr)

	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
