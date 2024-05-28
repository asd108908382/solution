package service

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
)

func consumerMessage(ctx context.Context) error {

	config := kafka.ReaderConfig{
		Brokers:     []string{GenConf()},
		Topic:       "demo",
		GroupID:     "test",
		StartOffset: kafka.LastOffset,
	}
	r := kafka.NewReader(config)
	defer func(r *kafka.Reader) {
		err := r.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(r)
	msg, err := r.ReadMessage(ctx)
	if err != nil {
		return err
	}
	log.Printf("Consumed message in sub-workflow: %s\n", string(msg.Value))
	if err := r.CommitMessages(ctx, msg); err != nil {
		return err
	}
	return nil
}
