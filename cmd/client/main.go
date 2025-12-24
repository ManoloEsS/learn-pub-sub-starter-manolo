package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// Connect to RabbitMQ server with default credentials
	const rabbitConnString = "amqp://guest:guest@localhost:5672/"
	fmt.Println("Starting Peril client...")

	conn, err := amqp.Dial(rabbitConnString)
	if err != nil {
		log.Fatalf(
			"could not connect to RabbitMQ server %s: %s", rabbitConnString, err,
		)
	}
	defer conn.Close()
	fmt.Printf(
		"Peril game server connected to RabbitMQ server %s\n", rabbitConnString,
	)

	// Create channel for publishing messages
	publishCh, err := conn.Channel()
	if err != nil {
		log.Fatalf("could not create channel: %v", err)
	}

	// Get player username from interactive welcome prompt
	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("could not get username: %v", err)
	}

	// Initialize game state for the player
	gameState := gamelogic.NewGameState(username)

	// Subscribe to army moves topic exchange for this player's moves
	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilTopic,
		routing.ArmyMovesPrefix+"."+gameState.GetUsername(),
		routing.ArmyMovesPrefix+".*",
		pubsub.SimpleQueueTransient,
		handlerMove(gameState, publishCh),
	)
	if err != nil {
		log.Fatalf("could not subscribe to army moves: %v", err)
	}

	// Subscribe to pause/resume messages via direct exchange
	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilDirect,
		routing.PauseKey+"."+gameState.GetUsername(),
		routing.PauseKey,
		pubsub.SimpleQueueTransient,
		handlerPause(gameState),
	)
	if err != nil {
		log.Fatalf("could not subscribe to pause: %v", err)
	}

	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilTopic,
		string(routing.WarRecognitionsPrefix),
		routing.WarRecognitionsPrefix+".*",
		pubsub.SimpleQueueDurable,
		handlerWar(gameState, publishCh),
	)
	// game loop REPL
	for {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}
		switch input[0] {
		case "spawn":
			err := gameState.CommandSpawn(input)
			if err != nil {
				fmt.Printf("could not spawn unit: %v\n", err)
				continue
			}
		case "move":
			mv, err := gameState.CommandMove(input)
			if err != nil {
				fmt.Printf("could not move unit: %v\n", err)
				continue
			}

			// publish move to publish channel
			err = pubsub.PublishJSON(publishCh,
				string(routing.ExchangePerilTopic),
				routing.ArmyMovesPrefix+"."+mv.Player.Username,
				mv,
			)
			if err != nil {
				fmt.Printf("could not publish move: %s", err)
			}
			fmt.Printf("Moved %v unit(s) to %s\n", len(mv.Units), mv.ToLocation)

		case "status":
			gameState.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			// TODO: implement spam command to publish logs
			fmt.Println("Spamming not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
			return
		default:
			fmt.Println("Invalid command...")
		}

	}

	fmt.Println("RabbitMQ connection closed...")

}
