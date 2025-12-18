package routing

// Message routing keys and exchange names for RabbitMQ.

const (
	// ArmyMovesPrefix is the routing key prefix for army movement messages
	ArmyMovesPrefix = "army_moves"

	// WarRecognitionsPrefix is the routing key prefix for war declaration messages
	WarRecognitionsPrefix = "war"

	// PauseKey is the routing key for pause/resume game state messages
	PauseKey = "pause"

	// GameLogSlug is the routing key for game log messages
	GameLogSlug = "game_logs"
)

// Exchange names used in the Peril game messaging system.

const (
	// ExchangePerilDirect is the direct exchange for targeted message delivery
	ExchangePerilDirect = "peril_direct"
	// ExchangePerilTopic is the topic exchange for pattern-based message routing
	ExchangePerilTopic = "peril_topic"
)
