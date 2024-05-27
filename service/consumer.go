package service

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
	"log"
	"sync"
)

var lock = &sync.Mutex{}

type Consumer struct {

	// kafka.Reader
	kafka.Reader
}

var singleInstance *Consumer

func GetInstance() *Consumer {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleInstance == nil {
			fmt.Println("Creating single instance now.")

			singleInstance = &Consumer{
				Reader: *kafka.NewReader(kafka.ReaderConfig{
					Brokers:  []string{GenConf()},
					Topic:    "demo",
					GroupID:  "test",
					MinBytes: 10e3, // 10KB
					MaxBytes: 10e6, // 10MB
					Dialer: &kafka.Dialer{
						SASLMechanism: plain.Mechanism{
							Username: "user1",
							Password: "eNlA6pwOgR",
						},
					},
				}),
			}
		} else {
			return singleInstance
		}
	} else {
		return singleInstance
	}

	return singleInstance
}

func InitConsumer(ctx context.Context) error {
	r := GetInstance()
	count := 0
	for count < 2 {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			return err
		}
		log.Printf("Consumed message in sub-workflow: %s\n", string(msg.Value))
		err = r.CommitMessages(ctx, msg)
		if err != nil {
			return err
		}
		count++

	}
	return nil
}
