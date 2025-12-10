package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

const rmqServer = "amqp://guest:guest@localhost:5672/"

func main() {
	fmt.Println("Starting Peril server...")

	conn, err := amqp.Dial(rmqServer)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ server %s: %s", rmqServer, err)
	}
	defer conn.Close()
	fmt.Printf("Peril game server connecte to RabbitMq server %s\n", rmqServer)

	pubCh, err := conn.Channel()
	if err != nil {
		log.Fatalf("could not create channel")
	}

	err = pubsub.PublishJson(pubCh, routing.ExchangePerilDirect, string(routing.PauseKey), routing.PlayingState{IsPaused: true})
	if err != nil {
		log.Printf("could not publish time: %s", err)
	}
	fmt.Println("Pause message sent!")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("RabbitMQ connection closed.")

}
