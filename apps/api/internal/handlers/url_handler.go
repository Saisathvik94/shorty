package handlers

import (
	"net/http"

	"github.com/Saisathvik94/shorty/apps/api/internal/services"
	"github.com/gin-gonic/gin"
)

type CreateURLRequest struct {
	OriginalURL string `json:"url" binding:"required"`
}

type URLHandler struct {
	service *services.URLService
}

func NewURLHandler(service *services.URLService) *URLHandler {
	return &URLHandler{
		service: service,
	}
}

func (h *URLHandler) CreateURL(c *gin.Context) {
	var req CreateURLRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shortCode, err := h.service.CreateURL(c.Request.Context(), req.OriginalURL)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"short_code": shortCode,
	})
}

func (h *URLHandler) Redirect(c *gin.Context) {
	shortCode := c.Param("shortCode")
	originalUrl, err := h.service.GetOriginalURL(c.Request.Context(), shortCode)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, originalUrl)

}
