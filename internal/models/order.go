package models

import "time"

type OrderStatus string

const (
	OrderCreated    OrderStatus = "CREATED"
	OrderProcessing OrderStatus = "PROCESSING"
	OrderPaid       OrderStatus = "PAID"
	OrderShipped    OrderStatus = "SHIPPED"
	OrderCompleted  OrderStatus = "COMPLETED"
	OrderCancelled  OrderStatus = "CANCELLED"
	PaymentFailed   OrderStatus = "PAYMENT_FAILED"
)

type Order struct {
	ID          string      `json:"id"`
	UserID      string      `json:"user_id"`
	Email       string      `json:"email"`
	Items       []OrderItem `json:"items"`
	TotalAmount float64     `json:"total_amount"`
	Status      OrderStatus `json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type OrderItem struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
}
