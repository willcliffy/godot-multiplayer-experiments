package broadcast

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"github.com/sony/sonyflake"
	pb "github.com/willcliffy/kilnwood-game-server/proto"
	"google.golang.org/protobuf/proto"
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
	for _, g := range mb.games {
		err := g.OnPlayerConnected(playerId)
		if err != nil {
			log.Error().Err(err).Msgf("failed to onplayerconnected")
			return
		}
	}
	log.Info().Msg("Starting read loop")
	mb.clientReadLoop(playerId, conn)
	log.Info().
		Msgf("Disconnected from player '%d'. Connection duration %v", playerId, time.Since(start))
}

func (mb *MessageBroker) clientReadLoop(playerId uint64, conn *websocket.Conn) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Warn().Err(err).Msgf("failed to read message from client, disconnecting")
			mb.unregisterConnection(playerId)
			return
		}

		var action pb.ClientAction
		err = proto.Unmarshal(message, &action)
		if err != nil {
			log.Warn().Err(err).Msgf("failed to unmarshal client action")
			continue
		}

		for _, g := range mb.games {
			g.OnActionReceived(playerId, &action)
		}
	}
}

func (mb *MessageBroker) unregisterConnection(playerId uint64) {
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
func (mb *MessageBroker) OnPlayerJoinGame(gameId, playerId uint64, response *pb.JoinGameResponse) {
	if len(response.Others) == 0 {
		response.Others = []*pb.Connect{}
	}

	responseBytes, err := proto.Marshal(response)
	if err != nil {
		log.Error().Err(err).Msgf("failed to marshal join game response: %+v", response)
		return
	}

	message := &pb.ServerMessage{
		Type:    pb.ServerMessageType_MESSAGE_JOIN,
		Payload: responseBytes,
	}

	messageBytes, err := proto.Marshal(message)
	if err != nil {
		log.Error().Err(err).Msgf("failed to marshal player join response: %v", message)
		return
	}

	mb.broadcastToPlayer(playerId, messageBytes)
}

// This satifies the MessageBroadcaster interface
func (mb *MessageBroker) OnPlayerLeftGame(gameId, playerId uint64) {
	mb.lock.Lock()
	defer mb.lock.Unlock()

	err := mb.playerConns[playerId].Close()
	if err != nil {
		log.Error().Err(err).Msgf("Failed to disconnect from %v", playerId)
	}

	delete(mb.playerConns, playerId)
}

// This satifies the MessageBroadcaster interface
func (mb *MessageBroker) OnGameTick(gameId uint64, tick *pb.GameTick) {
	tickBytes, err := proto.Marshal(tick)
	if err != nil {
		log.Error().Err(err).Msgf("failed to marshal game tick: %v", tick)
		return
	}

	message := &pb.ServerMessage{
		Type:    pb.ServerMessageType_MESSAGE_TICK,
		Payload: tickBytes,
	}

	messageBytes, err := proto.Marshal(message)
	if err != nil {
		log.Error().Err(err).Msgf("failed to marshal game tick server message: %v", tick)
		return
	}

	mb.broadcastToGame(gameId, messageBytes)
}

func (mb *MessageBroker) RegisterMessageReceiver(game MessageReceiver) uint64 {
	mb.lock.Lock()
	defer mb.lock.Unlock()

	gameId, _ := mb.playerIdGenerator.NextID()
	mb.games[gameId] = game

	return gameId
}

func (mb *MessageBroker) broadcastToGame(gameId uint64, payload []byte) {
	mb.lock.Lock()
	defer mb.lock.Unlock()

	// todo - support multiple games. This blasts to all
	for _, conn := range mb.playerConns {
		err := conn.WriteMessage(websocket.BinaryMessage, payload)
		if err != nil {
			log.Warn().Err(err).Msgf("Failed to broadcast to player")
		}
	}

	log.Debug().Msgf("broadcasted tick")
}

func (mb *MessageBroker) broadcastToPlayer(playerId uint64, payload []byte) {
	mb.lock.Lock()
	defer mb.lock.Unlock()

	conn, ok := mb.playerConns[playerId]
	if !ok {
		log.Warn().Msgf("Tried to broadcast to player that is not connected")
		return
	}

	err := conn.WriteMessage(websocket.BinaryMessage, payload)
	if err != nil {
		log.Warn().Err(err).Msgf("Failed to broadcast to player")
	}
}
