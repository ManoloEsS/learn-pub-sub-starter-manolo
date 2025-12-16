package pubsub

import (
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType int

const (
	SimpleQueueDurable SimpleQueueType = iota
	SimpleQueueTransient
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
	handler func(T),
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

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
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
