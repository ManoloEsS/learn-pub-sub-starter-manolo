// Package pubsub provides functionality for RabbitMQ message publishing and consumption.
package pubsub

import (
	"bytes"
	"encoding/gob"
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

type AckType int

const (
	Ack AckType = iota
	NackRequeue
	NackDiscard
)

func subscribe[T any](
	conn *amqp.Connection,
	exchange, // Exchange name to bind to
	queueName, // Queue name to create/consume from
	key string, // Routing key for binding
	queueType SimpleQueueType, // Queue persistence type
	handler func(T) AckType, // Message handler function
	unmarshaller func([]byte) (T, error),
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
			target, err := unmarshaller(m.Body)
			if err != nil {
				log.Printf("could not unmarshall %s: %s", m.Body, err)
				continue
			}

			ackType := handler(target)
			switch ackType {
			case Ack:
				err = m.Ack(false)
				log.Printf("Positive acknowledgement of type %s", ackType.String())
			case NackRequeue:
				err = m.Nack(false, true)
				log.Printf("Negative acknowledgement of type %s...requeueing", ackType.String())
			case NackDiscard:
				err = m.Nack(false, false)
				log.Printf("Negative acknowledgement of type %s...discarding", ackType.String())
			default:
				log.Printf("Invalid acknowledge type %v: %v", ackType, err)
			}
			if err != nil {
				log.Printf("Could not acknowledge message %s: %v", m.Body, err)
			}
		}
	}()

	return nil
}

// SubscribeJSON subscribes to a RabbitMQ queue and handles JSON messages of type T.
func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange, // Exchange name to bind to
	queueName, // Queue name to create/consume from
	key string, // Routing key for binding
	queueType SimpleQueueType, // Queue persistence type
	handler func(T) AckType, // Message handler function
) error {
	err := subscribe(conn, exchange, queueName, key, queueType, handler, unmarshalJSON)
	if err != nil {
		return err
	}
	return nil
}

func SubscribeGob[T any](
	conn *amqp.Connection,
	exchange, // Exchange name to bind to
	queueName, // Queue name to create/consume from
	key string, // Routing key for binding
	queueType SimpleQueueType, // Queue persistence type
	handler func(T) AckType, // Message handler function
) error {
	err := subscribe(conn, exchange, queueName, key, queueType, handler, unmarshalGob)
	if err != nil {
		return err
	}
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
		newQueue, err = chann.QueueDeclare(queueName, true, false, false, false, amqp.Table{
			"x-dead-letter-exchange": "peril_dlx",
		})
	case SimpleQueueTransient:
		newQueue, err = chann.QueueDeclare(queueName, false, true, true, false, amqp.Table{
			"x-dead-letter-exchange": "peril_dlx",
		})
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

func (a AckType) String() string {
	switch a {
	case Ack:
		return "Ack"
	case NackRequeue:
		return "NackRequeue"
	case NackDiscard:
		return "NackDiscard"
	default:
		return fmt.Sprintf("Unknown(%d)", a)
	}
}

func unmarshalJSON[T any](source []byte) (T, error) {
	var target T
	if err := json.Unmarshal(source, &target); err != nil {
		return target, fmt.Errorf("could not unmarshall %s: %s", source, err)
	}
	return target, nil
}

func unmarshalGob[T any](source []byte) (T, error) {
	var target T
	sourceStream := bytes.NewReader(source)
	if err := gob.NewDecoder(sourceStream).Decode(&target); err != nil {
		return target, fmt.Errorf("could not decode %s: %w", source, err)
	}

	return target, nil
}
