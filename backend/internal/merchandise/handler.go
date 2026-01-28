package merchandise

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetAll(c *gin.Context) {
	category := c.Query("category")

	var items []MerchandiseResponse
	var err error

	if category != "" && category != "all" {
		items, err = h.service.GetByCategory(category)
	} else {
		items, err = h.service.GetAll()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch merchandise"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": items})
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	merch := router.Group("/merchandise")
	{
		merch.GET("", h.GetAll)
	}
}
