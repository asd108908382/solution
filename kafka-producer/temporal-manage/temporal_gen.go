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

func CreatWorkflowAndRegister(c client.Client) worker.Worker {
	// create worker
	w := worker.New(c, "my-group", worker.Options{})

	// register workflow active
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
	// 创建解析器
	decoder := yaml.NewDecoder(file)
	// 配置对象
	var conf config.KafkaConfig
	// 解析 YAML 数据
	err = decoder.Decode(&conf)
	background := context.Background()
	kafka.InitProducer(conf, topic, background)
	return nil
}

func WorkflowFn(ctx workflow.Context) error {
	workflow.GetLogger(ctx).Info("Hello, Temporal!")
	childWorkflowOptions := workflow.ChildWorkflowOptions{
		WorkflowID:         "child-workflow",
		WorkflowRunTimeout: time.Minute,
	}

	for i := 0; i < 200; i++ {
		workflow.ExecuteChildWorkflow(ctx, childWorkflowOptions, ChildActiveFn)
	}
	return nil

}
