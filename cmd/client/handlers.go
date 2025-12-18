package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)

// Handler for pause/resume messages from the direct pause exchange
func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) {
	return func(playingState routing.PlayingState) {
		defer fmt.Print("> ")
		gs.HandlePause(playingState)
	}
}

// Handler for player's army move messages from the topic exchange
func handlerMove(gs *gamelogic.GameState) func(gamelogic.ArmyMove) {
	return func(move gamelogic.ArmyMove) {
		defer fmt.Print("> ")
		gs.HandleMove(move)
	}
}
