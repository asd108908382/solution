package service

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"sync"
)

var (
	readerInstance *kafka.Reader
	readerOnce     sync.Once
)

func getKafkaReader() *kafka.Reader {
	readerOnce.Do(func() {
		var err error
		config := kafka.ReaderConfig{
			Brokers:     []string{GenConf()},
			Topic:       "demo",
			GroupID:     "test",
			MinBytes:    10e3, // 10KB
			MaxBytes:    10e6, // 10MB
			StartOffset: kafka.LastOffset,
		}
		readerInstance = kafka.NewReader(config)
		if err != nil {
			log.Fatalf("Error creating kafka reader: %v", err)
		}
	})
	return readerInstance
}

func consumerMessage(ctx context.Context) error {
	instance := getKafkaReader()
	msg, err := instance.ReadMessage(ctx)
	if err != nil {
		return err
	}
	log.Printf("Consumed message in sub-workflow: %s\n", string(msg.Value))
	if err := instance.CommitMessages(ctx, msg); err != nil {
		return err
	}
	return nil
}
