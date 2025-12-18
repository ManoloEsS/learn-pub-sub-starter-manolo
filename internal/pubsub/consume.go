// Package pubsub provides functionality for RabbitMQ message publishing and consumption.
package pubsub

import (
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// SimpleQueueType defines the persistence characteristics of a RabbitMQ queue.
type SimpleQueueType int

const (
	// SimpleQueueDurable creates a queue that survives broker restarts
	SimpleQueueDurable SimpleQueueType = iota
	// SimpleQueueTransient creates a queue that is deleted when connection closes
	SimpleQueueTransient
)

// SubscribeJSON subscribes to a RabbitMQ queue and handles JSON messages of type T.
func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange, // Exchange name to bind to
	queueName, // Queue name to create/consume from
	key string, // Routing key for binding
	queueType SimpleQueueType, // Queue persistence type
	handler func(T), // Message handler function
) error {
	amqpChann, _, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return err
	}

	newChann, err := amqpChann.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for m := range newChann {
			var target T
			if err := json.Unmarshal(m.Body, &target); err != nil {
				log.Printf("could not unmarshall %s: %s", m.Body, err)
				continue
			}
			handler(target)
			if err := m.Ack(false); err != nil {
				log.Printf("could not acknowledge %s: %s", m.Body, err)
			}

		}
	}()

	return nil
}

// DeclareAndBind creates a RabbitMQ channel, declares a queue, and binds it to an exchange.
func DeclareAndBind(
	conn *amqp.Connection,
	exchange, // Exchange name
	queueName, // Queue name
	key string, // Routing key
	queueType SimpleQueueType, // Queue persistence type
) (*amqp.Channel, amqp.Queue, error) {
	chann, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("could not create channel: %w", err)
	}

	var newQueue amqp.Queue
	switch queueType {
	case SimpleQueueDurable:
		newQueue, err = chann.QueueDeclare(queueName, true, false, false, false, nil)
	case SimpleQueueTransient:
		newQueue, err = chann.QueueDeclare(queueName, false, true, true, false, nil)
	default:
		return nil, amqp.Queue{}, fmt.Errorf("invalid queue type: %v", queueType)
	}

	if err != nil {
		chann.Close()
		return nil, amqp.Queue{}, fmt.Errorf("could not declare queue (%v): %w", queueType, err)
	}

	err = chann.QueueBind(queueName, key, exchange, false, nil)
	if err != nil {
		chann.Close()
		return nil, amqp.Queue{}, fmt.Errorf("could not bind queue to exchange: %w", err)
	}

	return chann, newQueue, nil
}

// String returns a human-readable representation of the queue type.
func (s SimpleQueueType) String() string {
	switch s {
	case SimpleQueueDurable:
		return "Durable"
	case SimpleQueueTransient:
		return "Transient"
	default:
		return fmt.Sprintf("Unknown(%d)", s)
	}
}
