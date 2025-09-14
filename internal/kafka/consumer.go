package kafka

import (
	"context"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/vishalanarase/order-system/pkg/logger"
	"go.uber.org/zap"
)

type MessageHandler func(key string, value []byte, headers []kafka.Header) error

type Consumer struct {
	consumer *kafka.Consumer
	topics   []string
	handler  MessageHandler
	groupID  string
}

func NewConsumer(brokers []string, groupID string, topics []string, handler MessageHandler) (*Consumer, error) {
	config := &kafka.ConfigMap{
		"bootstrap.servers":         brokers[0],
		"group.id":                  groupID,
		"auto.offset.reset":         "earliest",
		"enable.auto.commit":        false,
		"max.poll.interval.ms":      300000,
		"session.timeout.ms":        10000,
		"heartbeat.interval.ms":     3000,
		"max.partition.fetch.bytes": 1048576,
	}

	c, err := kafka.NewConsumer(config)
	if err != nil {
		return nil, err
	}

	err = c.SubscribeTopics(topics, nil)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		consumer: c,
		topics:   topics,
		handler:  handler,
		groupID:  groupID,
	}, nil
}

func (c *Consumer) Consume(ctx context.Context) {
	logger.Log.Info("Starting consumer",
		zap.String("group_id", c.groupID),
		zap.Strings("topics", c.topics),
	)

	for {
		select {
		case <-ctx.Done():
			logger.Log.Info("Stopping consumer", zap.String("group_id", c.groupID))
			return
		default:
			msg, err := c.consumer.ReadMessage(100 * time.Millisecond)
			if err != nil {
				if err.(kafka.Error).Code() != kafka.ErrTimedOut {
					logger.Log.Error("Consumer error", zap.Error(err))
				}
				continue
			}

			logger.Log.Debug("Received message",
				zap.String("topic", *msg.TopicPartition.Topic),
				zap.String("key", string(msg.Key)),
			)

			if err := c.handler(string(msg.Key), msg.Value, msg.Headers); err != nil {
				logger.Log.Error("Error processing message",
					zap.Error(err),
					zap.String("key", string(msg.Key)),
				)
				// TODO: Implementing dead letter queue here
			} else {
				if _, err := c.consumer.CommitMessage(msg); err != nil {
					logger.Log.Error("Failed to commit message", zap.Error(err))
				}
			}
		}
	}
}

func (c *Consumer) Close() {
	c.consumer.Close()
}
