package api

import (
	"github.com/gin-gonic/gin"
	"github.com/vishalanarase/order-system/internal/services"
	"go.uber.org/zap"
)

func RegisterRoutes(log *zap.Logger, router *gin.Engine, orderService *services.OrderService) {
	healthHandler := NewHealthHandler(log)

	orderHandler := NewOrderHandler(log, orderService)

	// Health check endpoint
	router.GET("/health", healthHandler.HealthCheck)

	// Order endpoints
	orders := router.Group("/orders")
	{
		orders.POST("", orderHandler.CreateOrder)
	}
}
