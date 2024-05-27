package service

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

func InitConsumer(adr string, topic string, ctx context.Context) error {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{adr},
		Topic:    topic,
		GroupID:  "temporal-group", // 确保每个工作流实例有一个唯一的组ID
		MinBytes: 10e3,             // 10KB
		MaxBytes: 10e6,             // 10MB
		MaxWait:  2 * time.Millisecond,
	})

	defer func(r *kafka.Reader) {
		err := r.Close()
		if err != nil {

		}
	}(r)

	count := 0
	for count < 2 {
		select {
		case <-ctx.Done():
			return nil // 上下文被取消，结束消费
		default:
			msg, err := r.ReadMessage(ctx)
			if err != nil {
				return err
			}
			log.Printf("Consumed message in sub-workflow: %s\n", string(msg.Value))
			count++
		}
	}

	return nil

}
