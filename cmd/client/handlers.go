package main

import (
	"fmt"
	"time"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Handler for pause/resume messages from the direct pause exchange
func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) pubsub.AckType {
	return func(playingState routing.PlayingState) pubsub.AckType {
		defer fmt.Print("> ")
		gs.HandlePause(playingState)
		return pubsub.Ack
	}
}

// Handler for player's army move messages from the topic exchange
func handlerMove(gs *gamelogic.GameState, ch *amqp.Channel) func(gamelogic.ArmyMove) pubsub.AckType {
	return func(move gamelogic.ArmyMove) pubsub.AckType {
		defer fmt.Print("> ")

		moveOutcome := gs.HandleMove(move)

		switch moveOutcome {
		case gamelogic.MoveOutComeSafe:
			return pubsub.Ack
		case gamelogic.MoveOutcomeMakeWar:
			err := pubsub.PublishJSON(
				ch,
				routing.ExchangePerilTopic,
				routing.WarRecognitionsPrefix+"."+gs.GetUsername(),
				gamelogic.RecognitionOfWar{
					Attacker: move.Player,
					Defender: gs.GetPlayerSnap(),
				})
			if err != nil {
				fmt.Printf("error: %s\n", err)
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		case gamelogic.MoveOutcomeSamePlayer:
			return pubsub.NackDiscard
		default:
			fmt.Println("error: unknown war outcome")
			return pubsub.NackDiscard
		}
	}
}

func handlerWar(gs *gamelogic.GameState, ch *amqp.Channel) func(gamelogic.RecognitionOfWar) pubsub.AckType {
	return func(war gamelogic.RecognitionOfWar) pubsub.AckType {
		defer fmt.Print("> ")

		warOutcome, winner, loser := gs.HandleWar(war)

		var ackType pubsub.AckType
		var logMessage string
		publishLog := false

		switch warOutcome {
		case gamelogic.WarOutcomeNotInvolved:
			ackType = pubsub.NackRequeue
		case gamelogic.WarOutcomeNoUnits:
			ackType = pubsub.NackDiscard
		case gamelogic.WarOutcomeOpponentWon:
			ackType = pubsub.Ack
			publishLog = true
			logMessage = fmt.Sprintf("%s won a war against %s", winner, loser)
		case gamelogic.WarOutcomeYouWon:
			ackType = pubsub.Ack
			publishLog = true
			logMessage = fmt.Sprintf("%s won a war against %s", winner, loser)
		case gamelogic.WarOutcomeDraw:
			ackType = pubsub.Ack
			publishLog = true
			logMessage = fmt.Sprintf("A war between %s and %s resulted in a draw", winner, loser)
		default:
			fmt.Println("error: unknown move outcome")
			ackType = pubsub.NackDiscard
		}

		if publishLog {
			err := publishGameLog(ch, gs.GetUsername(), logMessage)
			if err != nil {
				fmt.Printf("error publishing game log: %v\n", err)
				return pubsub.NackRequeue
			}
		}

		return ackType

	}
}

func publishGameLog(ch *amqp.Channel, username, message string) error {
	gameLog := routing.GameLog{
		Message:     message,
		CurrentTime: time.Now(),
		Username:    username,
	}
	err := pubsub.PublishGob(ch, routing.ExchangePerilTopic, routing.GameLogSlug+"."+username, gameLog)
	if err != nil {
		return fmt.Errorf("could not publish game log: %w", err)
	}
	return nil
}
