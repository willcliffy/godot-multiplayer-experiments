package broadcast

import (
	"sync"

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
	// have to do this first, since BroadcastToGame also wants the lock
	// TODO - formalize disconnect message
	for gameId := range mb.games {
		_ = mb.BroadcastToGame(gameId, []byte("d:all"))
	}

	log.Debug().Msgf("Locking")
	mb.lock.Lock()
	log.Debug().Msgf("Locked")
	defer mb.lock.Unlock()

	for _, conn := range mb.playerConns {
		_ = conn.Close()
	}
	log.Debug().Msgf("Unlocking")
}

// This satisfies the `MessageBroadcaster` interface
// Note that this blocks the thread until the connection is broken
func (mb *MessageBroker) RegisterAndHandleWebsocketConnection(conn *websocket.Conn) {
	mb.lock.Lock()
	playerId, _ := mb.playerIdGenerator.NextID()
	mb.playerConns[playerId] = conn
	mb.lock.Unlock()

	log.Info().Msgf("Connected to new player assigned id: '%d'", playerId)
	mb.clientReadLoop(playerId, conn)
}

func (mb *MessageBroker) unregisterConnection(playerId uint64) {
	mb.lock.Lock()
	defer mb.lock.Unlock()

	if err := mb.playerConns[playerId].Close(); err != nil {
		log.Error().Err(err).Msgf("Failed to disconnect from %v", playerId)
	}

	delete(mb.playerConns, playerId)

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

	if err := mb.playerConns[playerId].Close(); err != nil {
		log.Error().Err(err).Msgf("Failed to disconnect from %v", playerId)
	}

	delete(mb.playerConns, playerId)
}

func (mb *MessageBroker) clientReadLoop(playerId uint64, conn *websocket.Conn) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			mb.unregisterConnection(playerId)
			return
		}

		// TODO - support multiple games
		for _, g := range mb.games {
			err = g.OnMessageReceived(playerId, message)
			if err != nil {
				log.Warn().Err(err).Send()
			}
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

	mb.lock.Lock()
	defer mb.lock.Unlock()

	// todo - support multiple games. This blasts to all
	for _, conn := range mb.playerConns {
		if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
			return err
		}
	}

	return nil
}

// This satisfies the util.Broadcaster interface
func (mb *MessageBroker) BroadcastToPlayer(playerId uint64, payload []byte) error {
	log.Debug().Msgf("broadcasting to player '%v' payload '%v'", playerId, string(payload))

	mb.lock.Lock()
	defer mb.lock.Unlock()

	if conn, ok := mb.playerConns[playerId]; ok {
		if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
			return err
		}
	}

	return nil
}
