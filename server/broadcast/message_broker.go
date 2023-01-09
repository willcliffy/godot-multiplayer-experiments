package broadcast

import (
	"net"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/sony/sonyflake"
)

const BufSize = 8192

type MessageBroker struct {
	conns     map[string]net.Conn
	playerIds map[string]uint64

	games map[uint64]MessageReceiver
	lock  sync.RWMutex

	playerIdGenerator *sonyflake.Sonyflake
}

func NewMessageBroker() *MessageBroker {
	return &MessageBroker{
		conns:     make(map[string]net.Conn),
		playerIds: make(map[string]uint64),
		games:     make(map[uint64]MessageReceiver),
		lock:      sync.RWMutex{},

		playerIdGenerator: sonyflake.NewSonyflake(sonyflake.Settings{}),
	}
}

func (self *MessageBroker) Close() {
	// have to do this first, since it also wants the lock
	// TODO - formalize disconnect message
	_ = self.broadcastMessage([]byte("d:all"))

	self.lock.Lock()
	defer self.lock.Unlock()

	for _, conn := range self.conns {
		conn.Close()
	}
}

func (self *MessageBroker) RegisterConnection(conn net.Conn) {
	self.lock.Lock()
	defer self.lock.Unlock()

	addr := conn.RemoteAddr().String()
	self.conns[addr] = conn
	playerId, _ := self.playerIdGenerator.NextID()
	self.playerIds[addr] = playerId

	go self.clientReadLoop(conn)

	log.Info().Msgf("Connected to %s", addr)
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

		addr := conn.RemoteAddr().String()

		for _, g := range self.games {
			err = g.OnMessageReceived(self.playerIds[addr], buf[:n])
			if err != nil {
				log.Warn().Err(err).Send()
			}
		}
	}
}

func (self *MessageBroker) broadcastMessage(msg []byte) error {
	self.lock.RLock()
	defer self.lock.RUnlock()

	for _, conn := range self.conns {
		_, err := conn.Write(msg)
		if err != nil {
			return err
		}
	}

	return nil
}

func (self *MessageBroker) broadcastMessageTo(addr string, msg []byte) error {
	self.lock.RLock()
	defer self.lock.RUnlock()

	if conn, ok := self.conns[addr]; ok {
		_, err := conn.Write(msg)
		if err != nil {
			return err
		}
	}

	return nil
}

func (self *MessageBroker) RegisterGame(gameId uint64, game MessageReceiver) {
	self.lock.Lock()
	defer self.lock.Unlock()

}

// This satisfies the util.Broadcaster interface
func (self *MessageBroker) BroadcastToGame(gameId uint64, payload []byte) error {
	// todo - support multiple games
	return self.broadcastMessage(payload)
}

// This satisfies the util.Broadcaster interface
func (self *MessageBroker) BroadcastToPlayer(playerId uint64, payload []byte) error {
	// todo - support multiple games
	for addr, pId := range self.playerIds {
		if pId == playerId {
			return self.broadcastMessageTo(addr, payload)
		}
	}

	return nil
}
