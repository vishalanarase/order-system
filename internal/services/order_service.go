package services

import (
	"time"

	"github.com/vishalanarase/order-system/internal/kafka"
	"github.com/vishalanarase/order-system/internal/models"
	"go.uber.org/zap"
)

type OrderService struct {
	log      *zap.Logger
	producer *kafka.Producer
}

func NewOrderService(log *zap.Logger, producer *kafka.Producer) *OrderService {
	return &OrderService{
		log:      log,
		producer: producer,
	}
}

func (s *OrderService) CreateOrder(order models.Order) error {
	s.log.Info("Service: CreateOrder called")
	orderEvent := models.OrderEvent{
		Type:      "ORDER_CREATED",
		Order:     order,
		Timestamp: time.Now(),
		Metadata: models.Metadata{
			CorrelationID: order.ID,
			Source:        "order-service",
		},
	}

	if err := s.producer.ProduceMessage(order.ID, orderEvent); err != nil {
		return err
	}

	s.log.Info("Order created and sent to Kafka",
		zap.String("order_id", order.ID),
		zap.String("user_id", order.UserID),
		zap.Float64("total_amount", order.TotalAmount),
	)

	return nil
}
