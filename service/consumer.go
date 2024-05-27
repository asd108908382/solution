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
		MaxWait:        1 * time.Second,
	}
	// 创建消费者
	reader := kafka.NewReader(readerConfig)

	// 读取一条消息
	message, err := reader.ReadMessage(context.Background())
	if err != nil {
		log.Fatalf("Failed to read message: %v", err)
	}

	// 打印消息内容
	log.Printf("Consumed message at offset %d: %s = %s\n", message.Offset, string(message.Key), string(message.Value))

	// 提交消费位移
	if err := reader.CommitMessages(context.Background(), message); err != nil {
		log.Fatalf("Failed to commit offset: %v", err)
	}

	// 关闭消费者
	if err := reader.Close(); err != nil {
		log.Fatalf("Failed to close reader: %v", err)
	}

	log.Println("Consumer closed")
}
