package accommodation

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Handler handles HTTP requests for accommodations
type Handler struct {
	service *Service
}

// NewHandler creates a new accommodation handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// GetAll handles GET /api/v1/accommodations
func (h *Handler) GetAll(c *gin.Context) {
	// Check for category filter
	category := c.Query("category")

	var accommodations []AccommodationResponse
	var err error

	if category != "" {
		// Validate category
		if err := h.service.ValidateCategory(category); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		accommodations, err = h.service.GetByCategory(category)
	} else {
		accommodations, err = h.service.GetAll()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve accommodations",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  accommodations,
		"count": len(accommodations),
	})
}

// GetByID handles GET /api/v1/accommodations/:id
func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")

	// A malformed UUID can never match, and would otherwise be a Postgres error
	if _, err := uuid.Parse(id); err != nil {
		respondNotFound(c)
		return
	}

	accommodation, err := h.service.GetByID(id)
	if err != nil {
		respondLookupError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": accommodation,
	})
}

// GetBySlug handles GET /api/v1/accommodations/slug/:slug
func (h *Handler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")

	accommodation, err := h.service.GetBySlug(slug)
	if err != nil {
		respondLookupError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": accommodation,
	})
}

// respondLookupError maps a single-record lookup error to a 404 or a 500
func respondLookupError(c *gin.Context, err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		respondNotFound(c)
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": "Failed to retrieve accommodation",
	})
}

func respondNotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"error": "Accommodation not found",
	})
}

// RegisterRoutes registers all accommodation routes
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/accommodations", h.GetAll)
	router.GET("/accommodations/:id", h.GetByID)
	router.GET("/accommodations/slug/:slug", h.GetBySlug)
}
