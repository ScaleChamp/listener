package components

import (
	"github.com/streadway/amqp"
	"gitlab.com/scalablespace/listener/app/models"
)

type ConsumerChannel interface {
	Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)
}

type ProducerChannel interface {
	Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
	Request(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) (*amqp.Delivery, error)
}

type ConsumerStopper interface {
	Stop()
}

type Consumer interface {
	Start()
}

type Workflow interface {
	Flow(*models.Task) models.Steps
}
