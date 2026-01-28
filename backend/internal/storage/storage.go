package storage

import (
	"context"
	"io"
)

// Storage is the interface for file storage operations
type Storage interface {
	// GetURL returns the public URL for a given file key
	GetURL(key string) string

	// Upload uploads a file and returns its key
	Upload(ctx context.Context, key string, file io.Reader) (string, error)

	// Delete removes a file
	Delete(ctx context.Context, key string) error

	// Exists checks if a file exists
	Exists(ctx context.Context, key string) (bool, error)
}
