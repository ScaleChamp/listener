package initializers

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/streadway/amqp"
	"gitlab.com/scalablespace/listener/app/models"
	"gitlab.com/scalablespace/listener/lib/components"
	"log"
	"os"
	"sync"
	"time"
)

type consumer struct {
	url        string
	queue      string
	stop       chan struct{}
	wg         sync.WaitGroup
	deliveries chan<- amqp.Delivery
}

func (e *consumer) Stop() {
	close(e.stop)
	e.wg.Wait()
	close(e.deliveries)
}

func (e *consumer) Start() {
	e.wg.Add(1)
	defer e.wg.Done()
	try := 0
	for {
		if try == 5 {
			os.Exit(1)
		}
		conn, err := amqp.Dial(e.url)
		if err != nil {
			try += 1
			time.Sleep(5 * time.Second)
			log.Println("dial", err)
			continue
		}
		ch, err := conn.Channel()
		if err != nil {
			log.Println("channel", err, conn.Close())
			continue
		}
		if err := ch.Qos(100, 0, false); err != nil {
			log.Println("qos", err, conn.Close())
			continue
		}
		const durable = true
		const noWait = false
		const autoDelete = false
		const exclusive = false
		const internal = false
		queue, err := ch.QueueDeclare(e.queue, durable, autoDelete, exclusive, noWait, nil)
		if err != nil {
			log.Println("qos", err, conn.Close())
			continue
		}
		if err := ch.ExchangeDeclare(queue.Name, "direct", durable, autoDelete, internal, noWait, nil); err != nil {
			log.Println("exchange", err, conn.Close())
			continue
		}
		if err := ch.QueueBind(queue.Name, "", queue.Name, noWait, nil); err != nil {
			log.Println("queue-bind", err, conn.Close())
			continue
		}
		id := uuid.NewV4()
		msgs, err := ch.Consume(
			e.queue,
			fmt.Sprintf("listener-%s", id),
			false,
			exclusive,
			false,
			noWait,
			nil,
		)
		if err != nil {
			log.Println("consume", err, conn.Close())
			continue
		}

		next, stopped := make(chan struct{}), make(chan struct{})

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()

			select {
			case <-e.stop:
				conn.Close()
				close(stopped)
			case <-next:
			}
		}()

		for msg := range msgs {
			e.deliveries <- msg
		}

		close(next)
		wg.Wait()

		select {
		case <-stopped:
			return
		default:
		}
	}
}

func NewDeliveries(e models.Environment) (<-chan amqp.Delivery, components.ConsumerStopper) {
	deliveries := make(chan amqp.Delivery, 100)
	c := &consumer{
		url:        e.RabbitURL,
		queue:      e.TasksQueue,
		deliveries: deliveries,
		stop:       make(chan struct{}),
	}
	go c.Start()
	return deliveries, c
}
