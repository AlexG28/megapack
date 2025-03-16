package message

import (
	"fmt"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"

	amqp "github.com/rabbitmq/amqp091-go"
)

func connectToRabbitMq() (*amqp.Connection, error) {
	numberOfRetries := 30
	var err error

	for range numberOfRetries {
		conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")

		if err == nil {
			return conn, nil
		}

		log.Println("failed to connect to RabbitMQ, trying again")
		time.Sleep(2 * time.Second)

	}

	return nil, fmt.Errorf("failed to dial rabbitmq after %d attempts: %v", numberOfRetries, err)
}

func OpenRabbitMQConnection(queueName string) (*amqp.Channel, *amqp.Queue, error) {

	conn, err := connectToRabbitMq()

	if err != nil {
		return nil, nil, fmt.Errorf("failed to dial rabbitmq: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create channel: %v", err)
	}

	q, err := ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to declare queue: %v", err)
	}

	err = ch.Qos(
		1, 0, false,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to set QoS: %v", err)
	}
	return ch, &q, nil
}

func GetMessages(ch *amqp091.Channel, q *amqp091.Queue) (<-chan amqp091.Delivery, error) {
	return ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
}
