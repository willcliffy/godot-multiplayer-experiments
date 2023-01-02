package main

import (
	"net"
	"sync"

	"github.com/pion/dtls/v2"
	"github.com/rs/zerolog/log"
	"github.com/willcliffy/kilnwood-game-server/game"
	"github.com/willcliffy/kilnwood-game-server/game/objects/actions"
)

const BufSize = 8192

type MessageBroker struct {
	conns     map[string]*dtls.Conn
	playerIds map[string]string
	games     []*game.Game
	lock      sync.RWMutex
}

func NewMessageBroker() *MessageBroker {
	return &MessageBroker{
		conns:     make(map[string]*dtls.Conn),
		playerIds: make(map[string]string),
		games:     make([]*game.Game, 0, 1),
		lock:      sync.RWMutex{},
	}
}

func (self *MessageBroker) Close() {
	// have to do this first, since it also wants the lock
	// TODO - formalize disconnect message
	self.broadcastMessage([]byte("d:all"))

	self.lock.Lock()
	defer self.lock.Unlock()

	for _, conn := range self.conns {
		conn.Close()
	}
}

func (self *MessageBroker) RegisterConnection(conn *dtls.Conn) {
	self.lock.Lock()
	defer self.lock.Unlock()

	self.conns[conn.RemoteAddr().String()] = conn

	go self.clientReadLoop(conn)

	log.Info().Msgf("Connected to %s", conn.RemoteAddr().String())
}

func (self *MessageBroker) unregisterConnection(conn net.Conn) {
	self.lock.Lock()
	defer self.lock.Unlock()

	delete(self.conns, conn.RemoteAddr().String())
	delete(self.playerIds, conn.RemoteAddr().String())

	if err := conn.Close(); err != nil {
		log.Error().Msgf("Failed to disconnect from %v: %v", conn.RemoteAddr().String(), err)
		return
	}

	log.Info().Msgf("Disconnected from %v", conn.RemoteAddr().String())
}

func (self *MessageBroker) clientReadLoop(conn net.Conn) {
	buf := make([]byte, BufSize)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			self.unregisterConnection(conn)
			return
		}

		action, err := actions.ParseActionFromMessage(string(buf[:n]))
		if err != nil {
			log.Warn().Msgf("failed to parse action from message: %v", err)
			continue
		}

		addr := conn.RemoteAddr().String()

		log.Debug().Msgf("Server got message %s from %s\n", action.Id(), addr)

		playerId, playerConnected := self.playerIds[addr]
		if playerConnected && playerId != action.SourcePlayer() {
			log.Info().Msgf("Discarding action - player action coming from another address")
			continue
		}

		if action.Type() == actions.ActionType_JoinGame {
			if playerConnected {
				continue // Player already connected, fail silently? TODO
			}
			self.lock.Lock()
			self.games[0].OnPlayerJoin(action.(*actions.JoinGameAction))
			self.playerIds[addr] = action.SourcePlayer()
			self.lock.Unlock()
			continue
		}

		// todo - support multiple games/lobbies?
		self.lock.Lock()
		_ = self.games[0].QueueAction(action)
		self.lock.Unlock()
		// todo - determine whether or not to broadcast messages to game
		self.broadcastMessage(buf[:n])
	}
}

func (self *MessageBroker) broadcastMessage(msg []byte) {
	self.lock.RLock()
	defer self.lock.RUnlock()

	for _, conn := range self.conns {
		_, err := conn.Write(msg)
		if err != nil {
			log.Error().Msgf("Failed to write message to %s: %v\n", conn.RemoteAddr().String(), err)
		}
	}
}

func (self *MessageBroker) RegisterGame(game *game.Game) {
	self.lock.Lock()
	defer self.lock.Unlock()

	self.games = append(self.games, game)
}

// This satisfies the util.Broadcaster interface
func (self *MessageBroker) Broadcast(gameId string, payload []byte) error {
	// todo - support multiple games
	self.broadcastMessage(payload)
	return nil
}
