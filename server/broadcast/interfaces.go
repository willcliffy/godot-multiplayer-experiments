package broadcast

type MessageBroadcaster interface {
	BroadcastToGame(gameId uint64, payload []byte) error
	BroadcastToPlayer(playerId uint64, payload []byte) error
}

type MessageReceiver interface {
	OnMessageReceived(uint64, []byte) error
	OnPlayerDisconnected(uint64)
}
