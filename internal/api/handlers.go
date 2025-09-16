package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RegisterRoutes(log *zap.Logger, router *gin.Engine) {
	healthHandler := NewHealthHandler(log)

	// Health check endpoint
	router.GET("/health", healthHandler.HealthCheck)
}
