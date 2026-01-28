package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Accommodation represents a luxury property available for booking
type Accommodation struct {
	ID          string `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Slug        string `gorm:"uniqueIndex;not null" json:"slug"`
	Title       string `gorm:"not null" json:"title"`
	Description string `gorm:"type:text" json:"description"`
	Location    string `json:"location"`
	PriceStart  int    `gorm:"not null" json:"priceStart"` // in cents
	Capacity    int    `gorm:"default:2" json:"capacity"`
	Category    string `json:"category"`                       // urban-luxury, island, historic, countryside
	Status      string `gorm:"default:'active'" json:"status"` // active, maintenance, inactive

	// JSON fields for flexible data
	Amenities StringArray `gorm:"type:jsonb" json:"amenities"`
	Images    ImageSet    `gorm:"type:jsonb" json:"images"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// StringArray is a custom type for []string that implements sql.Scanner and driver.Valuer
type StringArray []string

// Value implements driver.Valuer interface
func (s StringArray) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	return json.Marshal(s)
}

// Scan implements sql.Scanner interface
func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan StringArray: %v", value)
	}
	return json.Unmarshal(bytes, s)
}

// ImageSet holds references to accommodation images
type ImageSet struct {
	Hero    string   `json:"hero"`    // Main hero image
	Gallery []string `json:"gallery"` // Additional gallery images
}

// Value implements driver.Valuer interface
func (i ImageSet) Value() (driver.Value, error) {
	return json.Marshal(i)
}

// Scan implements sql.Scanner interface
func (i *ImageSet) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan ImageSet: %v", value)
	}
	return json.Unmarshal(bytes, i)
}

// TableName specifies the table name for Accommodation
func (Accommodation) TableName() string {
	return "accommodations"
}
