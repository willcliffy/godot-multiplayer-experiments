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

func (self *MessageBroker) Close() {
	// have to do this first, since BroadcastToGame also wants the lock
	// TODO - formalize disconnect message
	for gameId := range self.games {
		_ = self.BroadcastToGame(gameId, []byte("d:all"))
	}

	log.Debug().Msgf("Locking")
	self.lock.Lock()
	log.Debug().Msgf("Locked")
	defer self.lock.Unlock()

	for _, conn := range self.playerConns {
		_ = conn.Close()
	}
	log.Debug().Msgf("Unlocking")
}

// This satisfies the `MessageBroadcaster` interface
// Note that this blocks the thread until the connection is broken
func (self *MessageBroker) RegisterAndHandleWebsocketConnection(conn *websocket.Conn) {
	self.lock.Lock()
	playerId, _ := self.playerIdGenerator.NextID()
	self.playerConns[playerId] = conn
	self.lock.Unlock()

	log.Info().Msgf("Connected to new player assigned id: '%d'", playerId)
	self.clientReadLoop(playerId, conn)
}

func (self *MessageBroker) unregisterConnection(playerId uint64) {
	self.lock.Lock()
	defer self.lock.Unlock()

	if err := self.playerConns[playerId].Close(); err != nil {
		log.Error().Err(err).Msgf("Failed to disconnect from %v", playerId)
	}

	delete(self.playerConns, playerId)

	for _, game := range self.games {
		game.OnPlayerDisconnected(playerId)
	}

	log.Debug().Msgf("Disconnected from %v", playerId)
}

// This satifies the MessageBroadcaster interface
// This is the only allowed communication from the games to the MessageBroker
func (self *MessageBroker) OnPlayerLeft(playerId uint64) {
	if err := self.playerConns[playerId].Close(); err != nil {
		log.Error().Err(err).Msgf("Failed to disconnect from %v", playerId)
	}

	delete(self.playerConns, playerId)
}

func (self *MessageBroker) clientReadLoop(playerId uint64, conn *websocket.Conn) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			self.unregisterConnection(playerId)
			return
		}

		// TODO - support multiple games
		for _, g := range self.games {
			err = g.OnMessageReceived(playerId, message)
			if err != nil {
				log.Warn().Err(err).Send()
			}
		}
	}
}

func (self *MessageBroker) RegisterMessageReceiver(game MessageReceiver) uint64 {
	self.lock.Lock()
	defer self.lock.Unlock()

	gameId, _ := self.playerIdGenerator.NextID()
	self.games[gameId] = game

	return gameId
}

// This satisfies the util.Broadcaster interface
func (self *MessageBroker) BroadcastToGame(gameId uint64, payload []byte) error {
	log.Debug().Msgf("broadcasting to game '%v' payload '%v'", gameId, string(payload))

	self.lock.Lock()
	defer self.lock.Unlock()

	// todo - support multiple games. This blasts to all
	for _, conn := range self.playerConns {
		if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
			return err
		}
	}

	return nil
}

// This satisfies the util.Broadcaster interface
func (self *MessageBroker) BroadcastToPlayer(playerId uint64, payload []byte) error {
	log.Debug().Msgf("broadcasting to player '%v' payload '%v'", playerId, string(payload))

	self.lock.Lock()
	defer self.lock.Unlock()

	if conn, ok := self.playerConns[playerId]; ok {
		if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
			return err
		}
	}

	return nil
}
