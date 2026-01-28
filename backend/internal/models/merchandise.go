package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Merchandise struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `json:"description,omitempty"`
	Price       int64          `gorm:"not null" json:"price"` // stored in cents
	Category    string         `gorm:"index" json:"category"` // apparel, home, collectibles, accessories, lifestyle
	ImageURL    string         `json:"image_url"`
	IsNew       bool           `gorm:"default:false" json:"is_new"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
