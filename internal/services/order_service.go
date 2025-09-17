package services

import (
	"context"
	"encoding/json"
	"time"

	gokafka "github.com/confluentinc/confluent-kafka-go/kafka"
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

func (s *OrderService) StartOrderProcessor(ctx context.Context, brokers []string) {
	handler := func(key string, value []byte, headers []gokafka.Header) error {
		var event models.OrderEvent
		if err := json.Unmarshal(value, &event); err != nil {
			return err
		}

		switch event.Type {
		case "ORDER_CREATED":
			s.log.Info("Processing order", zap.String("order_id", event.Order.ID))

			// Simulate order processing
			time.Sleep(500 * time.Millisecond)

			// Update order status
			event.Order.Status = models.OrderProcessing
			event.Order.UpdatedAt = time.Now()

			// Create processing event
			processingEvent := models.OrderEvent{
				Type:      "ORDER_PROCESSING",
				Order:     event.Order,
				Timestamp: time.Now(),
				Metadata:  event.Metadata,
			}

			if err := s.producer.ProduceMessage(event.Order.ID, processingEvent); err != nil {
				return err
			}

			s.log.Info("Order processed", zap.String("order_id", event.Order.ID))
		}

		return nil
	}

	consumer, err := kafka.NewConsumer(brokers, "order-processor", []string{"orders"}, handler)
	if err != nil {
		s.log.Error("Failed to create order processor consumer", zap.Error(err))
		return
	}
	defer consumer.Close()

	go consumer.Consume(ctx)
}
