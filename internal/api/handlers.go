package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RegisterRoutes(log *zap.Logger, router *gin.Engine) {
	healthHandler := NewHealthHandler(log)

	orderHandler := NewOrderHandler(log)

	// Health check endpoint
	router.GET("/health", healthHandler.HealthCheck)

	// Order endpoints
	orders := router.Group("/orders")
	{
		orders.POST("", orderHandler.CreateOrder)
	}
}
