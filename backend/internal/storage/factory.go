package storage

import (
	"fmt"

	"github.com/wayne-stays/backend/internal/config"
)

// New creates a new storage instance based on configuration
func New(cfg *config.StorageConfig) (Storage, error) {
	switch cfg.Type {
	case "local":
		return NewLocalStorage(cfg.LocalPath, cfg.LocalBaseURL)
	case "s3":
		return nil, fmt.Errorf("S3 storage not yet implemented")
	default:
		return nil, fmt.Errorf("unknown storage type: %s", cfg.Type)
	}
}
