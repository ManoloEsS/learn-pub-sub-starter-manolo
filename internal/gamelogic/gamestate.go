package gamelogic

import (
	"sync"
)

// GameState represents the current state of the game for a player.
// It includes player information, pause status, and thread-safe access controls.
type GameState struct {
	Player Player        // Current player data
	Paused bool          // Game pause state
	mu     *sync.RWMutex // Mutex for thread-safe operations
}

// NewGameState creates a new game state for the specified username.
func NewGameState(username string) *GameState {
	return &GameState{
		Player: Player{
			Username: username,
			Units:    map[int]Unit{},
		},
		Paused: false,
		mu:     &sync.RWMutex{},
	}
}

// resumeGame sets the pause state to false, allowing game actions.
func (gs *GameState) resumeGame() {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	gs.Paused = false
}

// pauseGame sets the pause state to true, preventing game actions.
func (gs *GameState) pauseGame() {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	gs.Paused = true
}

// isPaused returns the current pause state of the game.
func (gs *GameState) isPaused() bool {
	gs.mu.RLock()
	defer gs.mu.RUnlock()
	return gs.Paused
}

// addUnit adds a new unit to the player's army.
func (gs *GameState) addUnit(u Unit) {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	gs.Player.Units[u.ID] = u
}

// removeUnitsInLocation removes all units from the specified location.
func (gs *GameState) removeUnitsInLocation(loc Location) {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	for k, v := range gs.Player.Units {
		if v.Location == loc {
			delete(gs.Player.Units, k)
		}
	}
}

// UpdateUnit updates an existing unit's information in the player's army.
func (gs *GameState) UpdateUnit(u Unit) {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	gs.Player.Units[u.ID] = u
}

// GetUsername returns the current player's username.
func (gs *GameState) GetUsername() string {
	return gs.Player.Username
}

// getUnitsSnap returns a snapshot of all units as a slice.
func (gs *GameState) getUnitsSnap() []Unit {
	gs.mu.RLock()
	defer gs.mu.RUnlock()
	Units := []Unit{}
	for _, v := range gs.Player.Units {
		Units = append(Units, v)
	}
	return Units
}

// GetUnit retrieves a unit by its ID, returning the unit and a boolean indicating if found.
func (gs *GameState) GetUnit(id int) (Unit, bool) {
	gs.mu.RLock()
	defer gs.mu.RUnlock()
	u, ok := gs.Player.Units[id]
	return u, ok
}

// GetPlayerSnap returns a deep copy of the current player state.
func (gs *GameState) GetPlayerSnap() Player {
	gs.mu.RLock()
	defer gs.mu.RUnlock()
	Units := map[int]Unit{}
	for k, v := range gs.Player.Units {
		Units[k] = v
	}
	return Player{
		Username: gs.Player.Username,
		Units:    Units,
	}
}
