package accommodation

import (
	"github.com/wayne-stays/backend/internal/models"
	"gorm.io/gorm"
)

// Repository handles database operations for accommodations
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new accommodation repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// listOrder keeps list results stable: listing order, then id as a tiebreak
const listOrder = "created_at, id"

// FindAll retrieves all active accommodations
func (r *Repository) FindAll() ([]models.Accommodation, error) {
	var accommodations []models.Accommodation
	err := r.db.Where("status = ?", "active").Order(listOrder).Find(&accommodations).Error
	return accommodations, err
}

// FindByID retrieves an accommodation by ID
func (r *Repository) FindByID(id string) (*models.Accommodation, error) {
	var accommodation models.Accommodation
	err := r.db.Where("id = ? AND status = ?", id, "active").First(&accommodation).Error
	if err != nil {
		return nil, err
	}
	return &accommodation, nil
}

// FindBySlug retrieves an accommodation by slug
func (r *Repository) FindBySlug(slug string) (*models.Accommodation, error) {
	var accommodation models.Accommodation
	err := r.db.Where("slug = ? AND status = ?", slug, "active").First(&accommodation).Error
	if err != nil {
		return nil, err
	}
	return &accommodation, nil
}

// FindByCategory retrieves accommodations by category
func (r *Repository) FindByCategory(category string) ([]models.Accommodation, error) {
	var accommodations []models.Accommodation
	err := r.db.Where("category = ? AND status = ?", category, "active").Order(listOrder).Find(&accommodations).Error
	return accommodations, err
}

// Create creates a new accommodation
func (r *Repository) Create(accommodation *models.Accommodation) error {
	return r.db.Create(accommodation).Error
}

// Update updates an existing accommodation
func (r *Repository) Update(accommodation *models.Accommodation) error {
	return r.db.Save(accommodation).Error
}

// Delete soft deletes an accommodation
func (r *Repository) Delete(id string) error {
	return r.db.Delete(&models.Accommodation{}, "id = ?", id).Error
}
