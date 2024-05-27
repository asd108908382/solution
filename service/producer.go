package service

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

func InitProducer(adr string, topic string, ctx context.Context) {
	writer := kafka.Writer{
		Addr:                   kafka.TCP(adr),
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
	log.Print("success producer")

}
