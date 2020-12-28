package services

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"gitlab.com/scalablespace/listener/app/models"
	"gitlab.com/scalablespace/listener/lib/components"
	"gitlab.com/scalablespace/listener/lib/taskflow"
	"log"
)

type consumer struct {
	messages       <-chan amqp.Delivery
	workflow       components.Workflow
	taskRepository components.TaskRepository
	taskFlowEngine *taskflow.Engine
}

func (c *consumer) Start() {
	for msg := range c.messages {
		m := msg
		go c.start(m)
	}
}

func (c *consumer) start(msg amqp.Delivery) {
	log.Println("delivery", string(msg.Body))
	m := new(models.Message)
	if err := json.Unmarshal(msg.Body, m); err != nil {
		log.Println("work is finished unexpected json", err)
		if err := msg.Ack(false); err != nil {
			panic(err)
		}
		return
	}
	task, err := c.taskRepository.Get(m.Id)
	if err != nil {
		log.Println(err)
		if err := msg.Reject(false); err != nil {
			panic(err)
		}
		return
	}
	taskData, err := json.Marshal(task)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("task is", string(taskData))
	if task.State > 0 {
		log.Println("already done")
		if err := msg.Ack(false); err != nil {
			panic(err)
		}
		return
	}
	commands := c.workflow.Flow(task)
	if err := c.taskFlowEngine.NewTaskFlow(task).Execute(commands); err != nil {
		fmt.Println(err)
	}
	if err := msg.Ack(false); err != nil {
		panic(err)
	}
}

func NewConsumer(
	router components.Workflow,
	messages <-chan amqp.Delivery,
	taskRepository components.TaskRepository,
	taskFlowEngine *taskflow.Engine,
) components.Consumer {
	return &consumer{
		messages:       messages,
		workflow:       router,
		taskRepository: taskRepository,
		taskFlowEngine: taskFlowEngine,
	}
}
