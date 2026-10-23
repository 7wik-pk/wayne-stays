package database

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"

	"github.com/pressly/goose/v3"
	"github.com/wayne-stays/backend/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

//go:embed migrations/*.sql
var embeddedMigrations embed.FS

// Connect establishes a database connection
func Connect(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connected successfully")
	return db, nil
}

// Migrate applies any pending SQL migrations from the migrations/ directory.
// Each migration runs once; goose records applied versions in goose_db_version.
func Migrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	migrations, err := fs.Sub(embeddedMigrations, "migrations")
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	provider, err := goose.NewProvider(goose.DialectPostgres, sqlDB, migrations)
	if err != nil {
		return fmt.Errorf("failed to create migration provider: %w", err)
	}

	results, err := provider.Up(context.Background())
	if err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}
	for _, r := range results {
		log.Printf("   applied %s (%s)", r.Source.Path, r.Duration)
	}

	log.Println("Migrations completed successfully")
	return nil
}
