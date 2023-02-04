package broadcast

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"github.com/sony/sonyflake"
	pb "github.com/willcliffy/kilnwood-game-server/proto"
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
	mb.clientReadLoop(playerId, conn)
	log.Info().
		Msgf("Disconnected from player '%d'. Connection duration %v", playerId, time.Since(start))
}

func (mb *MessageBroker) clientReadLoop(playerId uint64, conn *websocket.Conn) {
	for {
		// buf := make([]byte, 1024)
		// n, err := conn.UnderlyingConn().Read(buf)
		// if err != nil {
		// 	log.Warn().Err(err).Msgf("failed to read message from client, disconnecting")
		// 	return
		// }
		// log.Info().Msgf("n '%v', buf '%v'", n, string(buf[:n]))

		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Warn().Err(err).Msgf("failed to read message from client, disconnecting")
			mb.unregisterConnection(playerId)
			return
		}

		var action pb.ClientAction
		err = json.Unmarshal(message, &action)
		if err != nil {
			log.Warn().Err(err).Msgf("failed to unmarshal client action")
			continue
		}

		for _, g := range mb.games {
			if action.Type == int32(pb.ClientActionType_ACTION_CONNECT) {
				err := g.OnPlayerConnected(playerId)
				if err != nil {
					log.Error().Msgf("failed to onplayerconnected")
				}
				continue
			}
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

	responseBytes, err := json.Marshal(response)
	if err != nil {
		log.Error().Err(err).Msgf("failed to marshal join game response: %+v", response)
		return
	}

	message := &pb.ServerMessage{
		Type:    int32(pb.ServerMessageType_MESSAGE_JOIN),
		Payload: string(responseBytes),
	}

	messageBytes, err := json.Marshal(&message)
	if err != nil {
		log.Error().Err(err).Msgf("failed to marshal player join response: %v", message)
		return
	}
	log.Debug().Msg(string(responseBytes))
	log.Debug().Msg(string(messageBytes))

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

// TODO - this better
func tickIsEmpty(tick *pb.GameTick) bool {
	return len(tick.Connects) == 0 &&
		len(tick.Disconnects) == 0 &&
		len(tick.Moves) == 0 &&
		len(tick.Attacks) == 0
}

// This satifies the MessageBroadcaster interface
func (mb *MessageBroker) OnGameTick(gameId uint64, tick *pb.GameTick) {
	if tickIsEmpty(tick) {
		return
	}

	tickBytes, err := json.Marshal(tick)
	if err != nil {
		log.Error().Err(err).Msgf("failed to marshal game tick: %v", tick)
		return
	}

	message := &pb.ServerMessage{
		Type:    int32(pb.ServerMessageType_MESSAGE_TICK),
		Payload: string(tickBytes),
	}

	messageBytes, err := json.Marshal(&message)
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
		err := conn.WriteMessage(websocket.TextMessage, payload)
		if err != nil {
			log.Warn().Err(err).Msgf("Failed to broadcast to player")
		}
	}
}

func (mb *MessageBroker) broadcastToPlayer(playerId uint64, payload []byte) {
	mb.lock.Lock()
	defer mb.lock.Unlock()

	conn, ok := mb.playerConns[playerId]
	if !ok {
		log.Warn().Msgf("Tried to broadcast to player that is not connected")
		return
	}

	err := conn.WriteMessage(websocket.TextMessage, payload)
	if err != nil {
		log.Warn().Err(err).Msgf("Failed to broadcast to player")
	}
}
