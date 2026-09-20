package handlers

import (
	"errors"
	"net/http"

	"github.com/Saisathvik94/shorty/apps/api/internal/services"
	"github.com/gin-gonic/gin"
)

type CreateURLRequest struct {
	OriginalURL string `json:"url" binding:"required"`
	ExpiresAt   string `json:"expires_at" binding:"required"`
}

type UpdateExpirationRequest struct {
	ExpiresAt string `json:"expires_at" binding:"required"`
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

	shortCode, err := h.service.CreateURL(c.Request.Context(), req.OriginalURL, req.ExpiresAt)

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

func (h *URLHandler) DeactivateURL(c *gin.Context) {
	shortCode := c.Param("shortCode")
	err := h.service.DeactivateURL(c.Request.Context(), shortCode)

	if err != nil {
		if errors.Is(err, services.ErrorURLNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"short_code": shortCode,
		"message":    "URL deactivated successfully",
	})
}
func (h *URLHandler) DeleteURL(c *gin.Context) {
	shortCode := c.Param("shortCode")
	err := h.service.DeleteURL(c.Request.Context(), shortCode)

	if err != nil {
		if errors.Is(err, services.ErrorURLNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"short_code": shortCode,
		"message":    "URL deleted successfully",
	})
}
func (h *URLHandler) UpdateExpiration(c *gin.Context) {
	shortCode := c.Param("shortCode")

	var req UpdateExpirationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	err := h.service.UpdateExpiration(c.Request.Context(), req.ExpiresAt, shortCode)

	if err != nil {
		if errors.Is(err, services.ErrorURLNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"short_code": shortCode,
		"message":    "Updated successfully",
	})
}
