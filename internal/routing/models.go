// Package routing defines message types and routing constants for the Peril game.
package routing

import "time"

// PlayingState represents the pause/resume state of the game.
type PlayingState struct {
	IsPaused bool // Whether the game is currently paused
}

// GameLog represents a game event log entry with timestamp, message, and user.
type GameLog struct {
	CurrentTime time.Time // When the log entry was created
	Message     string    // The log message content
	Username    string    // The player who generated the log
}
