package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
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
	fmt.Printf("Peril game server connected to RabbitMq server %s\n", rmqServer)

	pubCh, err := conn.Channel()
	if err != nil {
		log.Fatalf("could not create channel")
	}

	gamelogic.PrintServerHelp()
	for {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}
		switch input[0] {
		case "pause":
			log.Println("Publishing paused game state...")
			err = pubsub.PublishJson(
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
			log.Println("Sending resume message...")
			err = pubsub.PublishJson(
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
