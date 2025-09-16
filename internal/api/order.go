package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type OrderHandler struct {
	log *zap.Logger
}

func NewOrderHandler(log *zap.Logger) *OrderHandler {
	return &OrderHandler{
		log: log,
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	h.log.Info("OrderHandler: CreateOrder called")
	c.JSON(http.StatusCreated, gin.H{
		"message": "Order created successfully",
	})
}
