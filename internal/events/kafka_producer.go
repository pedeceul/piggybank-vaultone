package events

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	w *kafka.Writer
}

func NewKafkaProducer(brokers []string, clientID string) *KafkaProducer {
	w := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
		Async:        false,
	}
	_ = clientID // reserved for future headers/metrics
	return &KafkaProducer{w: w}
}

func (p *KafkaProducer) ProduceJSON(ctx context.Context, topic string, key string, payload []byte) error {
	msg := kafka.Message{Key: []byte(key), Value: payload, Time: time.Now()}
	return p.w.WriteMessages(ctx, kafka.Message{Topic: topic, Key: msg.Key, Value: msg.Value, Time: msg.Time})
}

func (p *KafkaProducer) Close() error { return p.w.Close() }
