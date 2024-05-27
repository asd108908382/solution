package main

import (
	"log"
	"os"
	"solution/service"
)

func main() {
	service.GetInstance()
	if os.Getenv("ROLE") == "register" {
		register()
		workerInit()
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
	w := service.ProducerWorkflowRegister(c)
	w1 := service.ConsumerWorkflowRegister(c)
	err = w.Start()
	if err != nil {
		log.Fatalln("无法启动生产者工作流", err)
	}
	defer w.Stop()
	err = w1.Start()
	if err != nil {
		log.Fatalln("无法启动消费者工作流", err)
	}
	defer w1.Stop()
	select {}

}
