package service

import (
	"context"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
	"log"
	"os"
	"time"
)

func CreatClient() (client.Client, error) {
	c, err := client.NewClient(client.Options{
		HostPort: client.DefaultHostPort,
	})
	return c, err
}

func CreatProducerWorkflowAndRegister(c client.Client) worker.Worker {
	// create worker
	w := worker.New(c, "my-group", worker.Options{})

	// register workflow active
	w.RegisterWorkflow(ProducerWorkflowFn)
	w.RegisterActivity(ProducerChildActiveFn)
	return w
}

func ProducerChildActiveFn() error {
	background := context.Background()
	InitProducer(GenConf(), "demo", background)
	return nil
}

func ProducerWorkflowFn(ctx workflow.Context) error {
	workflow.GetLogger(ctx).Info("Hello, Temporal!")
	childWorkflowOptions := workflow.ChildWorkflowOptions{
		WorkflowID:         "child-workflow",
		WorkflowRunTimeout: time.Minute,
	}

	for i := 0; i < 200; i++ {
		workflow.ExecuteChildWorkflow(ctx, childWorkflowOptions, ProducerChildActiveFn)
	}
	return nil

}

func CreateConsumerWorkflowAndRegister(c client.Client) worker.Worker {
	// 创建工作流任务工作者
	w := worker.New(c, "my-group", worker.Options{})

	// 注册工作流和子工作流
	w.RegisterWorkflow(ConsumerWorkflowFn)
	w.RegisterActivity(ConsumerChildActiveFn)
	return w
}

func ConsumerChildActiveFn() error {
	background := context.Background()
	InitConsumer(GenConf(), "deom", background)
	return nil
}

func ConsumerWorkflowFn(ctx workflow.Context) error {
	childWorkflowOptions := workflow.ChildWorkflowOptions{
		WorkflowID:         "child-workflow",
		WorkflowRunTimeout: time.Minute,
	}
	for i := 0; i < 200; i++ {
		log.Println("execute word {}", i)
		workflow.ExecuteActivity(ctx, childWorkflowOptions, ConsumerChildActiveFn)
	}
	return nil
}

func GenConf() string {
	env, b := os.LookupEnv("KAFKA_HOME")
	if b {
		//
	}
	return env
}
