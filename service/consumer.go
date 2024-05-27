package service

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

func InitConsumer(adr string, topic string, ctx context.Context) {
	readerConfig := kafka.ReaderConfig{
		Brokers:        []string{adr},
		Topic:          topic,
		CommitInterval: 1 * time.Second,
		GroupID:        "demo-consumer",
		StartOffset:    kafka.LastOffset,
	}
	reader := kafka.NewReader(readerConfig)
	m, err := reader.FetchMessage(ctx)
	if err != nil {
		return
	}
	log.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	if err := reader.CommitMessages(ctx, m); err != nil {
		log.Fatal("failed to commit messages:", err)
	}

}
