package merchandise

import (
	"github.com/wayne-stays/backend/internal/models"
	"github.com/wayne-stays/backend/internal/storage"
)

type Service interface {
	GetAll() ([]MerchandiseResponse, error)
	GetByCategory(category string) ([]MerchandiseResponse, error)
}

type service struct {
	repo    Repository
	storage storage.Storage
}

func NewService(repo Repository, storage storage.Storage) Service {
	return &service{
		repo:    repo,
		storage: storage,
	}
}

type MerchandiseResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price"` // in dollars
	Category    string  `json:"category"`
	ImageURL    string  `json:"imageUrl"`
	IsNew       bool    `json:"isNew"`
}

func (s *service) GetAll() ([]MerchandiseResponse, error) {
	items, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	return s.toResponseList(items), nil
}

func (s *service) GetByCategory(category string) ([]MerchandiseResponse, error) {
	items, err := s.repo.GetByCategory(category)
	if err != nil {
		return nil, err
	}
	return s.toResponseList(items), nil
}

func (s *service) toResponse(m *models.Merchandise) MerchandiseResponse {
	return MerchandiseResponse{
		ID:          m.ID.String(),
		Name:        m.Name,
		Description: m.Description,
		Price:       float64(m.Price) / 100,
		Category:    m.Category,
		ImageURL:    s.storage.GetURL(m.ImageURL),
		IsNew:       m.IsNew,
	}
}

func (s *service) toResponseList(items []models.Merchandise) []MerchandiseResponse {
	responses := make([]MerchandiseResponse, len(items))
	for i, item := range items {
		responses[i] = s.toResponse(&item)
	}
	return responses
}
