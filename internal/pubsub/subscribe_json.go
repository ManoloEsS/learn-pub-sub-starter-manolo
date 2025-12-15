package pubsub

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
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
