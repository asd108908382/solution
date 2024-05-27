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
	message, err := reader.ReadMessage(ctx)
	if err != nil {
		//错误处理
	}
	err = reader.CommitMessages(ctx, message)
	err = reader.Close()
	log.Println("Message at offset %d: %s\n", message.Offset, message.Value)
}
