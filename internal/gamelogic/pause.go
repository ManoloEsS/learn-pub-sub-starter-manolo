package gamelogic

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)

// HandlePause processes pause/resume state changes from the server.
func (gs *GameState) HandlePause(ps routing.PlayingState) {
	defer fmt.Println("------------------------")
	fmt.Println()
	if ps.IsPaused {
		fmt.Println("==== Pause Detected ====")
		gs.pauseGame()
	} else {
		fmt.Println("==== Resume Detected ====")
		gs.resumeGame()
	}
}
