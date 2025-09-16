package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vishalanarase/order-system/internal/models"
	"github.com/vishalanarase/order-system/internal/services"
	"go.uber.org/zap"
)

type OrderHandler struct {
	log          *zap.Logger
	orderService *services.OrderService
}

func NewOrderHandler(log *zap.Logger, orderService *services.OrderService) *OrderHandler {
	return &OrderHandler{
		log:          log,
		orderService: orderService,
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	h.log.Info("OrderHandler: CreateOrder called")

	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		h.log.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Generate order ID and set timestamps
	order.ID = uuid.New().String()
	order.Status = models.OrderCreated
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	// Validate required fields
	if order.UserID == "" || order.Email == "" || len(order.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	// Calculate total amount
	var totalAmount float64
	for _, item := range order.Items {
		totalAmount += item.Price * float64(item.Quantity)
	}
	order.TotalAmount = totalAmount

	if err := h.orderService.CreateOrder(order); err != nil {
		h.log.Error("Failed to create order", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Order created successfully",
		"order_id": order.ID,
		"status":   order.Status,
	})
}
