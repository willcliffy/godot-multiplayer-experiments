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
	conns map[string]*dtls.Conn
	games []*game.Game
	lock  sync.RWMutex
}

func NewMessageBroker() *MessageBroker {
	return &MessageBroker{conns: make(map[string]*dtls.Conn)}
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

		log.Debug().Msgf("Server got message: %v\n", action)

		self.broadcastMessage(buf[:n])

		if action.Type() == actions.ActionType_JoinGame {
			self.games[0].OnPlayerJoin(action.(*actions.JoinGameAction))
			continue
		}

		// todo - support multiple games/lobbies?
		self.games[0].QueueAction(action)

		// todo - determine whether or not to broadcast messages to game

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
	self.games = append(self.games, game)
}
