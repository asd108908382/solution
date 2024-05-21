package main

import (
	"go.temporal.io/sdk/worker"
	"log"
	"solution/kafka-consumer/temporal-manage"
)

func main() {
	c, err := temporal_manage.CreatClient()
	defer c.Close()
	if err != nil {
		log.Fatalln("can not create Temporal client", err)
	}
	w := temporal_manage.CreatWorkflow(c)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("can not init Temporal worker", err)
	}
}
