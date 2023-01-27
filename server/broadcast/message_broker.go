package broadcast

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"github.com/sony/sonyflake"
)

type MessageBroker struct {
	lock sync.Mutex

	games map[uint64]MessageReceiver

	playerIdGenerator *sonyflake.Sonyflake
	playerConns       map[uint64]*websocket.Conn
}

func NewMessageBroker() *MessageBroker {
	return &MessageBroker{
		playerConns: make(map[uint64]*websocket.Conn),
		games:       make(map[uint64]MessageReceiver),
		lock:        sync.Mutex{},

		playerIdGenerator: sonyflake.NewSonyflake(sonyflake.Settings{}),
	}
}

func (mb *MessageBroker) Close() {
	for _, game := range mb.games {
		game.Close()
	}

	mb.lock.Lock()
	defer mb.lock.Unlock()

	for _, conn := range mb.playerConns {
		_ = conn.Close()
	}
}

// This satisfies the `MessageBroadcaster` interface
// Note that this blocks the thread until the connection is closed
func (mb *MessageBroker) RegisterAndHandleWebsocketConnection(conn *websocket.Conn) {
	mb.lock.Lock()
	playerId, _ := mb.playerIdGenerator.NextID()
	mb.playerConns[playerId] = conn
	mb.lock.Unlock()

	start := time.Now()
	log.Info().Msgf("Connected to new player assigned id: '%d'", playerId)
	mb.clientReadLoop(playerId, conn)
	log.Info().
		Msgf("Disconnected from player '%d'. Connection duration %v", playerId, time.Since(start))
}

func (mb *MessageBroker) unregisterConnection_LOCK(playerId uint64) {
	mb.lock.Lock()

	err := mb.playerConns[playerId].Close()
	if err != nil {
		log.Error().Err(err).Msgf("Failed to disconnect from %v", playerId)
	}

	delete(mb.playerConns, playerId)
	mb.lock.Unlock()

	for _, game := range mb.games {
		game.OnPlayerDisconnected(playerId)
	}

	log.Debug().Msgf("Disconnected from %v", playerId)
}

// This satifies the MessageBroadcaster interface
// This is the only allowed communication from the games to the MessageBroker
func (mb *MessageBroker) OnPlayerLeft(playerId uint64) {
	mb.lock.Lock()
	defer mb.lock.Unlock()

	err := mb.playerConns[playerId].Close()
	if err != nil {
		log.Error().Err(err).Msgf("Failed to disconnect from %v", playerId)
	}

	delete(mb.playerConns, playerId)
}

func (mb *MessageBroker) clientReadLoop(playerId uint64, conn *websocket.Conn) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			mb.unregisterConnection_LOCK(playerId)
			return
		}

		// TODO - support multiple games
		for _, g := range mb.games {
			g.OnMessageReceived(playerId, message)
		}
	}
}

func (mb *MessageBroker) RegisterMessageReceiver(game MessageReceiver) uint64 {
	mb.lock.Lock()
	defer mb.lock.Unlock()

	gameId, _ := mb.playerIdGenerator.NextID()
	mb.games[gameId] = game

	return gameId
}

// This satisfies the util.Broadcaster interface
func (mb *MessageBroker) BroadcastToGame(gameId uint64, payload []byte) error {
	log.Debug().Msgf("broadcasting to game '%v' payload '%v'", gameId, string(payload))
	return mb.broadcastToGame_LOCK(gameId, payload)
}

func (mb *MessageBroker) broadcastToGame_LOCK(gameId uint64, payload []byte) error {
	mb.lock.Lock()
	defer mb.lock.Unlock()

	// todo - support multiple games. This blasts to all
	for _, conn := range mb.playerConns {
		err := conn.WriteMessage(websocket.TextMessage, payload)
		if err != nil {
			return err
		}
	}

	return nil
}

// This satisfies the util.Broadcaster interface
func (mb *MessageBroker) BroadcastToPlayer(playerId uint64, payload []byte) error {
	log.Debug().Msgf("broadcasting to player '%v' payload '%v'", playerId, string(payload))
	return mb.broadcastToPlayer_LOCK(playerId, payload)
}

func (mb *MessageBroker) broadcastToPlayer_LOCK(playerId uint64, payload []byte) error {
	log.Debug().Msgf("broadcasting to player '%v' payload '%v'", playerId, string(payload))

	mb.lock.Lock()
	defer mb.lock.Unlock()

	conn, ok := mb.playerConns[playerId]
	if !ok {
		log.Warn().Msgf("")
	}

	err := conn.WriteMessage(websocket.TextMessage, payload)
	if err != nil {
		return err
	}

	return nil
}
