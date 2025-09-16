package services

import (
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
	return nil
}
