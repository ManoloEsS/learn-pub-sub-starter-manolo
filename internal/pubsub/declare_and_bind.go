package pubsub

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

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

type SimpleQueueType int

const (
	SimpleQueueDurable SimpleQueueType = iota
	SimpleQueueTransient
)

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
