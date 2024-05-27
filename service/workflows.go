package service

import (
	"context"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
	"google.golang.org/protobuf/types/known/durationpb"
	"log"
	"os"
	"time"
)

type CronResult struct {
	RunTime time.Time
}

func CreatClient() (client.Client, error) {
	c, err := client.Dial(client.Options{
		HostPort: "10.100.183.60:7233",
	})
	namespaceClient, err := client.NewNamespaceClient(client.Options{
		HostPort: "10.100.183.60:7233",
	})
	err = namespaceClient.Register(context.Background(), &workflowservice.RegisterNamespaceRequest{
		Namespace:                        "default",
		WorkflowExecutionRetentionPeriod: durationpb.New(time.Hour * 24 * 30),
	})
	if err.Error() == "Namespace already exists." {
		return c, nil
	} else if err != nil {
		return c, err
	}
	return c, err
}

func InitProducerWorkFlow(err error, c client.Client) {
	workflowID := "producer-parent-workflow"
	workflowOptions := client.StartWorkflowOptions{
		ID:                    workflowID,
		TaskQueue:             "workflow",
		CronSchedule:          "* * * * *",
		WorkflowIDReusePolicy: enums.WORKFLOW_ID_REUSE_POLICY_TERMINATE_IF_RUNNING,
	}

	workflowRun, err := c.ExecuteWorkflow(context.Background(), workflowOptions, ProducerWorkflowFn)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}
	log.Println("Started workflow",
		"WorkflowID", workflowRun.GetID(), "RunID", workflowRun.GetRunID())
}

func ProducerChildWorkflowFn(ctx workflow.Context) (string, error) {

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx1 := workflow.WithActivityOptions(ctx, ao)
	var result string
	// Start from 0 for first cron job
	lastRunTime := time.Time{}
	// Check to see if there was a previous cron job
	if workflow.HasLastCompletionResult(ctx) {
		var lastResult CronResult
		if err := workflow.GetLastCompletionResult(ctx, &lastResult); err == nil {
			lastRunTime = lastResult.RunTime
		}
	}
	thisRunTime := workflow.Now(ctx)

	err := workflow.ExecuteActivity(ctx1, ProducerChildActiveFn, lastRunTime, thisRunTime).Get(ctx, &result)
	if err != nil {
		//
	}
	return result, nil
}

func ProducerChildActiveFn(ctx context.Context, lastRunTime, thisRunTime time.Time) error {
	InitProducer(GenConf(), "demo", ctx)
	return nil
}

func ProducerWorkflowFn(ctx workflow.Context) (string, error) {
	workflow.GetLogger(ctx).Info("Hello, Temporal!")
	childWorkflowOptions := workflow.ChildWorkflowOptions{
		WorkflowID:         "producer-child-workflow",
		WorkflowRunTimeout: 2 * time.Minute,
	}

	ctx = workflow.WithChildOptions(ctx, childWorkflowOptions)
	var result string
	for i := 0; i < 200; i++ {
		err := workflow.ExecuteChildWorkflow(ctx, ProducerChildWorkflowFn).Get(ctx, &result)
		if err != nil {
			//
		}
	}
	return result, nil

}

func ConsumerChildWorkflowFn(ctx workflow.Context) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx1 := workflow.WithActivityOptions(ctx, ao)
	// Start from 0 for first cron job
	lastRunTime := time.Time{}
	// Check to see if there was a previous cron job
	if workflow.HasLastCompletionResult(ctx) {
		var lastResult CronResult
		if err := workflow.GetLastCompletionResult(ctx, &lastResult); err == nil {
			lastRunTime = lastResult.RunTime
		}
	}
	thisRunTime := workflow.Now(ctx)
	err := workflow.ExecuteActivity(ctx1, ConsumerChildActiveFn, lastRunTime, thisRunTime).Get(ctx, nil)
	if err != nil {
		// Cron job failed
		// Next cron will still be scheduled by the Server
		workflow.GetLogger(ctx).Error("Cron job failed.", "Error", err)
		return nil
	}
	return nil
}

func ConsumerChildActiveFn(ctx context.Context, lastRunTime, thisRunTime time.Time) error {
	InitConsumer(GenConf(), "demo", ctx)
	return nil
}

func ConsumerWorkflowFn(ctx workflow.Context) error {
	childWorkflowOptions := workflow.ChildWorkflowOptions{
		WorkflowID:         "consumer-child-workflow",
		WorkflowRunTimeout: 2 * time.Minute,
	}
	ctx = workflow.WithChildOptions(ctx, childWorkflowOptions)
	var result string
	for i := 0; i < 200; i++ {
		err := workflow.ExecuteChildWorkflow(ctx, ConsumerChildWorkflowFn).Get(ctx, &result)
		if err != nil {
			//
		}
	}
	return nil
}

func InitConsumerWorkFlow(err error, c client.Client) {
	workflowID := "consumer-parent-workflow"
	workflowOptions := client.StartWorkflowOptions{
		ID:                    workflowID,
		TaskQueue:             "workflow",
		CronSchedule:          "*/2 * * * *",
		WorkflowIDReusePolicy: enums.WORKFLOW_ID_REUSE_POLICY_TERMINATE_IF_RUNNING,
	}

	workflowRun, err := c.ExecuteWorkflow(context.Background(), workflowOptions, ConsumerWorkflowFn)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}
	log.Println("Started workflow",
		"WorkflowID", workflowRun.GetID(), "RunID", workflowRun.GetRunID())
}

func GenConf() string {
	env, b := os.LookupEnv("KAFKA_HOME")
	if b {
		//
		return "10.100.216.49"
	}
	return env
}

func WorkflowRegister(c client.Client) worker.Worker {
	// 创建工作流任务工作者
	w := worker.New(c, "workflow", worker.Options{})

	// 注册工作流和子工作流
	w.RegisterWorkflow(ConsumerWorkflowFn)
	w.RegisterWorkflow(ConsumerChildWorkflowFn)
	w.RegisterWorkflow(ProducerWorkflowFn)
	w.RegisterWorkflow(ProducerChildWorkflowFn)
	w.RegisterActivity(ConsumerChildActiveFn)
	w.RegisterActivity(ProducerChildActiveFn)
	return w
}
