package accommodation

import (
	"fmt"

	"github.com/wayne-stays/backend/internal/models"
	"github.com/wayne-stays/backend/internal/storage"
)

// Service handles business logic for accommodations
type Service struct {
	repo    *Repository
	storage storage.Storage
}

// NewService creates a new accommodation service
func NewService(repo *Repository, storage storage.Storage) *Service {
	return &Service{
		repo:    repo,
		storage: storage,
	}
}

// AccommodationResponse is the API response format
type AccommodationResponse struct {
	ID          string           `json:"id"`
	Slug        string           `json:"slug"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Location    string           `json:"location"`
	PriceStart  float64          `json:"priceStart"` // in dollars
	Capacity    int              `json:"capacity"`
	Category    string           `json:"category"`
	Amenities   []string         `json:"amenities"`
	Images      ImageResponseSet `json:"images"`
}

// ImageResponseSet holds image URLs
type ImageResponseSet struct {
	Hero    string   `json:"hero"`
	Gallery []string `json:"gallery"`
}

// GetAll retrieves all accommodations with full image URLs
func (s *Service) GetAll() ([]AccommodationResponse, error) {
	accommodations, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	return s.toResponseList(accommodations), nil
}

// GetByID retrieves an accommodation by ID with full image URLs
func (s *Service) GetByID(id string) (*AccommodationResponse, error) {
	accommodation, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	response := s.toResponse(accommodation)
	return &response, nil
}

// GetBySlug retrieves an accommodation by slug with full image URLs
func (s *Service) GetBySlug(slug string) (*AccommodationResponse, error) {
	accommodation, err := s.repo.FindBySlug(slug)
	if err != nil {
		return nil, err
	}

	response := s.toResponse(accommodation)
	return &response, nil
}

// GetByCategory retrieves accommodations by category
func (s *Service) GetByCategory(category string) ([]AccommodationResponse, error) {
	accommodations, err := s.repo.FindByCategory(category)
	if err != nil {
		return nil, err
	}

	return s.toResponseList(accommodations), nil
}

// toResponse converts a model to a response with full image URLs
func (s *Service) toResponse(acc *models.Accommodation) AccommodationResponse {
	return AccommodationResponse{
		ID:          acc.ID,
		Slug:        acc.Slug,
		Title:       acc.Title,
		Description: acc.Description,
		Location:    acc.Location,
		PriceStart:  float64(acc.PriceStart) / 100, // Convert cents to dollars
		Capacity:    acc.Capacity,
		Category:    acc.Category,
		Amenities:   []string(acc.Amenities),
		Images:      s.buildImageURLs(acc.Images),
	}
}

// toResponseList converts a list of models to responses
func (s *Service) toResponseList(accommodations []models.Accommodation) []AccommodationResponse {
	responses := make([]AccommodationResponse, len(accommodations))
	for i, acc := range accommodations {
		responses[i] = s.toResponse(&acc)
	}
	return responses
}

// buildImageURLs converts image paths to full URLs
func (s *Service) buildImageURLs(images models.ImageSet) ImageResponseSet {
	gallery := make([]string, len(images.Gallery))
	for i, img := range images.Gallery {
		gallery[i] = s.storage.GetURL(img)
	}

	return ImageResponseSet{
		Hero:    s.storage.GetURL(images.Hero),
		Gallery: gallery,
	}
}

// ValidateCategory checks if a category is valid
func (s *Service) ValidateCategory(category string) error {
	validCategories := map[string]bool{
		"urban-luxury": true,
		"island":       true,
		"historic":     true,
		"countryside":  true,
	}

	if !validCategories[category] {
		return fmt.Errorf("invalid category: %s", category)
	}
	return nil
}
