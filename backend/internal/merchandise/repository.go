package merchandise

import (
	"strings"

	"github.com/wayne-stays/backend/internal/models"
	"gorm.io/gorm"
)

type Repository interface {
	GetAll() ([]models.Merchandise, error)
	GetByCategory(category string) ([]models.Merchandise, error)
}

// listOrder keeps list results stable: listing order, then id as a tiebreak
const listOrder = "created_at, id"

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetAll() ([]models.Merchandise, error) {
	var items []models.Merchandise
	err := r.db.Order(listOrder).Find(&items).Error
	return items, err
}

func (r *repository) GetByCategory(category string) ([]models.Merchandise, error) {
	var items []models.Merchandise
	err := r.db.Where("category = ?", strings.ToLower(category)).Order(listOrder).Find(&items).Error
	return items, err
}
