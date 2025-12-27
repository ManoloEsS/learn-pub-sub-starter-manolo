package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	rmqServer = "amqp://guest:guest@localhost:5672/"
)

func main() {
	// Start server and connect to RabbitMQ
	fmt.Println("Starting Peril server...")

	conn, err := amqp.Dial(rmqServer)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ server %s: %s", rmqServer, err)
	}
	defer conn.Close()
	fmt.Printf("Peril game server connected to RabbitMq server %s\n", rmqServer)

	// Create channel
	pubCh, err := conn.Channel()
	if err != nil {
		log.Fatalf("could not create channel")
	}

	// Bind to topic pause exchange
	err = pubsub.SubscribeGob(
		conn,
		routing.ExchangePerilTopic,
		routing.GameLogSlug,
		routing.GameLogSlug+".*",
		pubsub.SimpleQueueDurable,
		handlerGameLogs(),
	)
	if err != nil {
		log.Fatalf("could not start consuming logs: %v", err)
	}

	// Start server REPL
	gamelogic.PrintServerHelp()
	for {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}
		switch input[0] {
		case "pause":
			fmt.Println("Publishing paused game state...")
			err = pubsub.PublishJSON(
				pubCh,
				routing.ExchangePerilDirect,
				string(routing.PauseKey),
				routing.PlayingState{IsPaused: true},
			)
			if err != nil {
				log.Printf("could not publish time: %s", err)
			}
			fmt.Println("Pause message sent!")
		case "resume":
			fmt.Println("Sending resume message...")
			err = pubsub.PublishJSON(
				pubCh,
				routing.ExchangePerilDirect,
				string(routing.PauseKey),
				routing.PlayingState{IsPaused: false},
			)
			if err != nil {
				log.Printf("could not publish time: %s", err)
			}
			fmt.Println("Resume message sent!")
		case "quit":
			log.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid command...")
		}
	}
}
