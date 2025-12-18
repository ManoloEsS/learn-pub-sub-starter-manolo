// Package gamelogic contains the core game logic for the Peril strategy game.
// It handles player management, unit spawning, movement, combat, and game state.
package gamelogic

// Player represents a game participant with their username and military units.
type Player struct {
	Username string
	Units    map[int]Unit
}

// UnitRank represents the type of military unit with different power levels.
type UnitRank string

const (
	// RankInfantry represents basic infantry units with power level 1
	RankInfantry = "infantry"
	// RankCavalry represents mounted units with power level 5
	RankCavalry = "cavalry"
	// RankArtillery represents heavy artillery units with power level 10
	RankArtillery = "artillery"
)

// Unit represents a military unit with ID, type, and location on the game map.
type Unit struct {
	ID       int      // Unique identifier for the unit
	Rank     UnitRank // Type/rank of the unit (infantry, cavalry, artillery)
	Location Location // Current location of the unit on the map
}

// ArmyMove represents a movement order containing the player, units being moved, and destination.
type ArmyMove struct {
	Player     Player   // Player initiating the move
	Units      []Unit   // List of units being moved
	ToLocation Location // Destination location for the units
}

// RecognitionOfWar represents a declaration of war between two players.
type RecognitionOfWar struct {
	Attacker Player // Player initiating the war
	Defender Player // Player being attacked
}

// Location represents a geographic area on the game map.
type Location string

// getAllRanks returns a set of all valid unit types in the game.
func getAllRanks() map[UnitRank]struct{} {
	return map[UnitRank]struct{}{
		RankInfantry:  {},
		RankCavalry:   {},
		RankArtillery: {},
	}
}

// getAllLocations returns a set of all valid geographic locations in the game.
func getAllLocations() map[Location]struct{} {
	return map[Location]struct{}{
		"americas":   {},
		"europe":     {},
		"africa":     {},
		"asia":       {},
		"australia":  {},
		"antarctica": {},
	}
}
