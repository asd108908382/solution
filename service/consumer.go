package service

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"sync"
	"time"
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
					MaxBytes: 10e6, // 10MB
					MaxWait:  5 * time.Minute,
				}),
			}
			return singleInstance
		} else {
			return singleInstance
		}
	} else {
		return singleInstance
	}

	return singleInstance
}

func ConsumerMessage(ctx context.Context) error {
	//r := kafka.NewReader(kafka.ReaderConfig{
	//	Brokers:  []string{GenConf()},
	//	Topic:    "demo",
	//	GroupID:  "test",
	//	MaxBytes: 10e6, // 10MB
	//	MaxWait:  2 * time.Millisecond,
	//})
	//
	//defer func(r *kafka.Reader) {
	//	err := r.Close()
	//	if err != nil {
	//
	//	}
	//}(r)
	instance := GetInstance()
	msg, err := instance.Reader.ReadMessage(ctx)
	if err != nil {
		return err
	}
	log.Printf("Consumed message in sub-workflow: %s\n", string(msg.Value))
	if err := instance.Reader.CommitMessages(ctx, msg); err != nil {
		return err
	}
	return nil
}
