package main

import (
	"net"
	"sync"

	"github.com/pion/dtls/v2"
	"github.com/rs/zerolog/log"
	"github.com/willcliffy/kilnwood-game-server/game/actions"
)

const BufSize = 8192

type MessageBroker struct {
	conns map[string]*dtls.Conn
	lock  sync.RWMutex
}

func NewMessageBroker() *MessageBroker {
	return &MessageBroker{conns: make(map[string]*dtls.Conn)}
}

func (h *MessageBroker) RegisterClient(conn *dtls.Conn) {
	h.lock.Lock()
	defer h.lock.Unlock()

	log.Info().Msgf("Connected to %s\n", conn.RemoteAddr())

	h.conns[conn.RemoteAddr().String()] = conn

	go h.clientReadLoop(conn)
}

func (h *MessageBroker) unregisterClient(conn net.Conn) {
	h.lock.Lock()
	defer h.lock.Unlock()

	delete(h.conns, conn.RemoteAddr().String())

	if err := conn.Close(); err != nil {
		log.Error().Msgf("Failed to disconnect from %v: %v", conn.RemoteAddr(), err)
		return
	}

	log.Info().Msgf("Disconnected from %v", conn.RemoteAddr())
}

func (h *MessageBroker) clientReadLoop(conn net.Conn) {
	b := make([]byte, BufSize)
	for {
		n, err := conn.Read(b)
		if err != nil {
			h.unregisterClient(conn)
			return
		}

		action, err := actions.ParseActionFromMessage(string(b[:n]))
		if err != nil {
			log.Warn().Msgf("failed to parse action from message: %v", err)
			continue
		}

		log.Debug().Msgf("Server got message: %v\n", action)
		// todo - perform action

		h.broadcastMessage(b[:n])
	}
}

func (h *MessageBroker) broadcastMessage(msg []byte) {
	h.lock.RLock()
	defer h.lock.RUnlock()
	for _, conn := range h.conns {
		_, err := conn.Write(msg)
		if err != nil {
			log.Error().Msgf("Failed to write message to %s: %v\n", conn.RemoteAddr(), err)
		}
	}
}
