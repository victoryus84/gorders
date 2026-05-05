package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(addr string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(addr),
			Balancer: &kafka.LeastBytes{},
			MaxAttempts: 3,
			RequiredAcks: kafka.RequireOne,
		},
	}
}

func (p *Producer) Publish(ctx context.Context, topic string, key string, payload []byte) error {
	return p.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic, // Îl dăm la fiecare mesaj în parte
		Key:   []byte(key),
		Value: payload,
	})
}

func (p *Producer) Close() error {
	return p.writer.Close()
}