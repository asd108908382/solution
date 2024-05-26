package main

import (
	"go.temporal.io/sdk/worker"
	"log"
	"solution/service"
)

func main() {
	c, err := service.CreatClient()
	if err != nil {
		log.Fatalln("无法创建 Temporal 客户端", err)
	}
	service.InitProducerWorkFlow(err, c)
	service.InitConsumerWorkFlow(err, c)

	w := service.WorkflowRegister(c)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
