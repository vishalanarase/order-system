package api

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	healthHandler := NewHealthHandler()

	// Health check endpoint
	router.GET("/health", healthHandler.HealthCheck)
}
