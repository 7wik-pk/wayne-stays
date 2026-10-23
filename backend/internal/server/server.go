package server

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/wayne-stays/backend/internal/accommodation"
	"github.com/wayne-stays/backend/internal/config"
	"github.com/wayne-stays/backend/internal/merchandise"
	"github.com/wayne-stays/backend/internal/storage"
	"gorm.io/gorm"
)

// NewRouter wires up middleware, static files and all API routes.
// It is separate from main so tests can exercise the real routes over HTTP.
func NewRouter(cfg *config.Config, db *gorm.DB, store storage.Storage) *gin.Engine {
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

		// Initialize merchandise components
		merchandiseRepo := merchandise.NewRepository(db)
		merchandiseService := merchandise.NewService(merchandiseRepo, store)
		merchandiseHandler := merchandise.NewHandler(merchandiseService)

		// Register merchandise routes
		merchandiseHandler.RegisterRoutes(v1)
	}

	return router
}
