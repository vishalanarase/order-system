package kafka

import (
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/vishalanarase/order-system/pkg/logger"

	"go.uber.org/zap"
)

type Producer struct {
	producer *kafka.Producer
	topic    string
}

func NewProducer(brokers []string, topic string) (*Producer, error) {
	config := &kafka.ConfigMap{
		"bootstrap.servers": brokers[0],
		"client.id":         "order-service-producer",
		"acks":              "all",
		"retries":           3,
		"retry.backoff.ms":  1000,
	}

	p, err := kafka.NewProducer(config)
	if err != nil {
		return nil, err
	}

	// Start delivery report goroutine
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					logger.Log.Error("Delivery failed",
						zap.Error(ev.TopicPartition.Error),
						zap.String("topic", *ev.TopicPartition.Topic),
					)
				} else {
					logger.Log.Debug("Message delivered",
						zap.String("topic", *ev.TopicPartition.Topic),
						zap.Int32("partition", ev.TopicPartition.Partition),
						zap.Int64("offset", int64(ev.TopicPartition.Offset)),
					)
				}
			}
		}
	}()

	return &Producer{
		producer: p,
		topic:    topic,
	}, nil
}

func (p *Producer) ProduceMessage(key string, value interface{}) error {
	valueBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &p.topic,
			Partition: kafka.PartitionAny,
		},
		Key:   []byte(key),
		Value: valueBytes,
	}

	return p.producer.Produce(message, nil)
}

func (p *Producer) Close() {
	p.producer.Flush(15 * 1000)
	p.producer.Close()
}
