package main

import (
	"go.temporal.io/sdk/worker"
	"log"
	"os"
	"solution/service"
)

func main() {
	if os.Getenv("ROLE") == "register" {
		register()
		select {}
	} else {
		workerInit()
	}
}

func register() {
	c, err := service.CreatClient()
	if err != nil {
		log.Fatalln("无法创建 Temporal 客户端", err)
	}
	service.InitProducerWorkFlow(err, c)
	service.InitConsumerWorkFlow(err, c)
}

func workerInit() {

	c, err := service.CreatClient()
	if err != nil {
		log.Fatalln("无法创建 Temporal 客户端", err)
	}
	w := service.WorkflowRegister(c)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}

}
