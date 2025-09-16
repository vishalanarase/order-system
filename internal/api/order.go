package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
}

func NewOrderHandler() *OrderHandler {
	return &OrderHandler{}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {

	c.JSON(http.StatusCreated, gin.H{
		"message": "Order created successfully",
	})
}
