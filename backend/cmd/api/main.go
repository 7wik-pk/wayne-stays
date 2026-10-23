package main

import (
	"log"

	"github.com/wayne-stays/backend/internal/config"
	"github.com/wayne-stays/backend/internal/database"
	"github.com/wayne-stays/backend/internal/server"
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
	log.Printf("Storage initialized (%s)", cfg.Storage.Type)

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
	if err := database.Seed(db); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	router := server.NewRouter(cfg, db, store)

	// Start server
	addr := ":" + cfg.Server.Port
	log.Printf("Server starting on http://localhost%s", addr)
	log.Printf("API available at http://localhost%s/api/v1", addr)

	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
