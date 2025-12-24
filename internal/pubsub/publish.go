package pubsub

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// PublishJSON marshals a value to JSON and publishes it to a RabbitMQ exchange.
func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	jsonVal, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("could not marshal value to json: %w", err)
	}

	return ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        jsonVal,
	})
}

func PublishGob[T any](ch *amqp.Channel, exchange, key string, val T) error {
	var gobVal bytes.Buffer
	encoder := gob.NewEncoder(&gobVal)

	err := encoder.Encode(val)
	if err != nil {
		return fmt.Errorf("could not encode value to gob: %w", err)
	}

	return ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{
		ContentType: "application/gob",
		Body:        gobVal.Bytes(),
	})
}
