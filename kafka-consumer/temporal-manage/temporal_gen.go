package temporal_manage

import (
	"context"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"solution/kafka-common/config"
	"solution/kafka-common/kafka"
	"time"
)

var topic = "demo"

func CreatClient() (client.Client, error) {
	c, err := client.NewClient(client.Options{
		HostPort: client.DefaultHostPort,
	})
	return c, err
}

func CreatWorkflow(c client.Client) worker.Worker {
	// 创建工作流任务工作者
	w := worker.New(c, "my-group", worker.Options{})

	// 注册工作流和子工作流
	w.RegisterWorkflow(WorkflowFn)
	w.RegisterActivity(ChildActiveFn)
	return w
}

func ChildActiveFn() error {
	file, err := os.Open("config.yaml")
	if err != nil {
		log.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()
	decoder := yaml.NewDecoder(file)
	var conf config.KafkaConfig
	err = decoder.Decode(&conf)
	background := context.Background()
	kafka.InitConsumer(conf, topic, background)
	return nil
}

func WorkflowFn(ctx workflow.Context) error {
	childWorkflowOptions := workflow.ChildWorkflowOptions{
		WorkflowID:         "child-workflow",
		WorkflowRunTimeout: time.Minute,
	}
	for i := 0; i < 200; i++ {
		log.Println("execute word {}", i)
		workflow.ExecuteActivity(ctx, childWorkflowOptions, ChildActiveFn)
	}
	return nil

}
