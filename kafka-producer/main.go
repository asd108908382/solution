package main

import (
	"go.temporal.io/sdk/worker"
	"log"
	"solution/kafka-producer/temporal-manage"
)

var topic = "demo"

func main() {
	// 创建 Temporal 客户端连接
	c, err := temporal_manage.CreatClient()
	defer c.Close()
	if err != nil {
		log.Fatalln("无法创建 Temporal 客户端", err)
	}
	//创建工作流
	w := temporal_manage.CreatWorkflowAndRegister(c)

	err = w.Run(worker.InterruptCh())
	defer w.Stop()
	if err != nil {
		log.Fatalln("无法启动 Temporal worker", err)
	}
}
