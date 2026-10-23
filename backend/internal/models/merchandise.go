package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Merchandise represents a Gotham-themed item sold in the shop
type Merchandise struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `json:"description,omitempty"`
	Price       int64          `gorm:"not null" json:"price"` // stored in cents
	Category    string         `gorm:"index" json:"category"` // apparel, home, collectibles, accessories, lifestyle
	ImageURL    string         `json:"imageUrl"`
	IsNew       bool           `gorm:"default:false" json:"isNew"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for Merchandise
func (Merchandise) TableName() string {
	return "merchandise"
}
