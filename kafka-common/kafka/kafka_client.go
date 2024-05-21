package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"solution/kafka-common/config"
	"time"
)

func InitConsumer(config config.KafkaConfig, topic string, ctx context.Context) {
	readerConfig := kafka.ReaderConfig{
		Brokers:        []string{config.Address},
		Topic:          topic,
		CommitInterval: 1 * time.Second,
		GroupID:        "demo-consumer",
		StartOffset:    kafka.LastOffset,
	}
	for true {
		reader := kafka.NewReader(readerConfig)
		defer reader.Close()
		message, err := reader.ReadMessage(ctx)
		if err != nil {
			//错误处理
		}
		log.Fatalln("Message at offset %d: %s\n", message.Offset, message.Value)
	}
}

func InitProducer(kafkaConfig config.KafkaConfig, topic string, ctx context.Context) {
	writer := kafka.Writer{
		Addr:                   kafka.TCP(kafkaConfig.Address),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		WriteTimeout:           1 * time.Second,
		RequiredAcks:           kafka.RequireNone,
		AllowAutoTopicCreation: true,
	}
	defer func(writer *kafka.Writer) {
		err := writer.Close()
		if err != nil {

		}
	}(&writer)

	err := writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(time.Now().String()),
		Value: []byte(time.Now().String() + "hello world"),
	})
	if err != nil {
		//retry or other?
	}

}
