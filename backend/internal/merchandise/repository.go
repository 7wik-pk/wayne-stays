package merchandise

import (
	"github.com/wayne-stays/backend/internal/models"
	"gorm.io/gorm"
)

type Repository interface {
	GetAll() ([]models.Merchandise, error)
	GetByCategory(category string) ([]models.Merchandise, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetAll() ([]models.Merchandise, error) {
	var items []models.Merchandise
	err := r.db.Find(&items).Error
	return items, err
}

func (r *repository) GetByCategory(category string) ([]models.Merchandise, error) {
	var items []models.Merchandise
	err := r.db.Where("category ILIKE ?", category).Find(&items).Error
	return items, err
}
