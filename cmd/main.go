package main

import (
	"log"
	"solution/service"
)

func main() {
	c, err := service.CreatClient()
	defer c.Close()
	if err != nil {
		log.Fatalln("无法创建 Temporal 客户端", err)
	}
	w := service.CreateConsumerWorkflowAndRegister(c)
	w2 := service.CreatProducerWorkflowAndRegister(c)
	err = w.Start()
	defer w.Stop()
	if err != nil {
		log.Fatalln("can not init Temporal worker", err)
	}
	err = w2.Start()
	defer w2.Stop()
	if err != nil {
		log.Fatalln("can not init Temporal worker", err)
	}
	select {}

}
