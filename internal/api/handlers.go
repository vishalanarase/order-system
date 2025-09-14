package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
}

func NewOrderHandler() *OrderHandler {
	return &OrderHandler{}
}

func (h *OrderHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   "order-service",
	})
}

func RegisterRoutes(router *gin.Engine) {

	handler := NewOrderHandler()

	// Health check endpoint
	router.GET("/health", handler.HealthCheck)
}
