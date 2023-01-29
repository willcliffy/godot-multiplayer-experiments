package broadcast

import "github.com/gorilla/websocket"

type MessageBroadcaster interface {
	// from the gameserver to MessageBroadcaster
	RegisterAndHandleWebsocketConnection(conn *websocket.Conn)
	RegisterMessageReceiver(receiver MessageReceiver) uint64

	// From MessageBroadcaster to all registered MessageReceivers
	BroadcastToGame(gameId uint64, payload []byte) error
	BroadcastToPlayer(playerId uint64, payload []byte) error

	// From MessageReceivers to their MessageBroadcaster
	OnPlayerLeft(playerId uint64)
}

type MessageReceiver interface {
	OnPlayerDisconnected(uint64)
	Close()
}
