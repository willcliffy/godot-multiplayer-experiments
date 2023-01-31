package game

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/willcliffy/kilnwood-game-server/broadcast"
	pb "github.com/willcliffy/kilnwood-game-server/proto"
	"google.golang.org/protobuf/proto"
)

// 600 bpm or 10 bps
const gameTick = 100 * time.Millisecond

type Game struct {
	id          uint64
	broadcaster broadcast.MessageBroadcaster
	done        chan bool

	clock *time.Ticker
	tick  uint32

	gameMap *GameMap
	players map[uint64]*Player

	disconnectsQueued map[uint64]*pb.Disconnect
	movementsQueued   map[uint64]*pb.Move
	attacksQueued     map[uint64]*pb.Attack
}

func NewGame(gameId uint64, broadcaster broadcast.MessageBroadcaster) *Game {
	return &Game{
		id:          gameId,
		broadcaster: broadcaster,
		done:        make(chan bool),

		gameMap: NewGameMap(),
		players: make(map[uint64]*Player),

		disconnectsQueued: make(map[uint64]*pb.Disconnect),
		movementsQueued:   make(map[uint64]*pb.Move),
		attacksQueued:     make(map[uint64]*pb.Attack),
	}
}

func (g Game) Id() string {
	return fmt.Sprint(g.id)
}

func (g *Game) Start() {
	g.clock = time.NewTicker(gameTick)
	go g.run()
}

func (g *Game) Stop() {
	g.done <- true
	g.clock.Stop()
}

func (g *Game) run() {
	for {
		select {
		case <-g.done:
			break
		case <-g.clock.C:
			g.tick += 1
			tickResult := g.processQueue()
			if tickResult == nil {
				continue
			}

			g.broadcaster.OnGameTick(g.id, tickResult)
		}
	}
}

func (g *Game) OnPlayerConnected(playerId uint64) error {
	p, playerInGame := g.players[playerId]
	if !playerInGame {
		// TODO - allow specifying team. for now, give everyone a random color
		color := RandomTeamColor()
		p = NewPlayer(playerId, color)
		g.players[playerId] = p
	}

	err := g.gameMap.AddPlayer(p)
	if err != nil {
		delete(g.players, p.Id)
		return err
	}

	_, err = g.gameMap.SpawnPlayer(p)
	if err != nil {
		delete(g.players, p.Id)
		return err
	}

	var playerList []*pb.JoinGameResponse

	for pId, p := range g.players {
		if p == nil || pId == playerId {
			continue
		}

		playerList = append(playerList, &pb.JoinGameResponse{
			PlayerId: pId,
			Color:    p.Color,
			Spawn:    p.Location,
		})

	}

	msg := &pb.JoinGameResponse{
		PlayerId: p.Id,
		Color:    p.Color,
		Spawn:    p.Location,
		Others:   playerList,
	}

	g.broadcaster.OnPlayerJoinGame(g.id, playerId, msg)
	return nil
}

func (g *Game) OnPlayerDisconnected(playerId uint64) {
	delete(g.players, playerId)
	g.disconnectsQueued[playerId] = &pb.Disconnect{PlayerId: playerId}
	log.Info().Msgf("queued disconnect for player %d", playerId)
}

func (g *Game) Close() {
	for pId := range g.players {
		g.disconnectsQueued[pId] = &pb.Disconnect{PlayerId: pId}
	}
}

func (g *Game) OnActionReceived(playerId uint64, action *pb.Action) {
	switch action.Type {
	case pb.ActionType_ACTION_PING:
		// nyi
	case pb.ActionType_ACTION_DISCONNECT:
		var disconnect *pb.Disconnect
		err := proto.Unmarshal(action.Payload, disconnect)
		if err != nil {
			log.Warn().Err(err).Msgf("failed to unmarshal disconnect")
		}

		g.disconnectsQueued[playerId] = disconnect
	case pb.ActionType_ACTION_MOVE:
		var disconnect *pb.Disconnect
		err := proto.Unmarshal(action.Payload, disconnect)
		if err != nil {
			log.Warn().Err(err).Msgf("failed to unmarshal move")
		}

		g.disconnectsQueued[playerId] = disconnect
	case pb.ActionType_ACTION_ATTACK:
		var attack *pb.Attack
		err := proto.Unmarshal(action.Payload, attack)
		if err != nil {
			log.Warn().Err(err).Msgf("failed to unmarshal attack")
		}

		g.attacksQueued[playerId] = attack
	default:
		log.Error().Msgf("could not unmarshal unknown action type: %v", action.Type)
	}
}

func (g *Game) processQueue() *pb.GameTick {
	var disconnectsProcessed []*pb.Disconnect
	var movementsProcessed []*pb.Move
	var attacksProcessed []*pb.Attack

	for _, disconnect := range g.disconnectsQueued {
		g.OnPlayerDisconnected(disconnect.PlayerId)
		disconnectsProcessed = append(disconnectsProcessed, disconnect)
	}

	for _, move := range g.movementsQueued {
		err := g.gameMap.ApplyMovement(move)
		if err != nil {
			log.Error().Err(err).Msgf("Failed action: could not apply MoveAction")
			continue
		}

		movementsProcessed = append(movementsProcessed, move)
	}

	for _, attack := range g.attacksQueued {
		_, err := g.gameMap.ApplyAttack(attack)
		if err != nil {
			log.Error().Err(err).Msgf("Failed action: could not apply MoveAction")
			continue
		}

		attacksProcessed = append(attacksProcessed, attack)
	}

	// reset actionQueue for the next tick
	g.disconnectsQueued = make(map[uint64]*pb.Disconnect)
	g.movementsQueued = make(map[uint64]*pb.Move)
	g.attacksQueued = make(map[uint64]*pb.Attack)

	return &pb.GameTick{
		Tick:        g.tick,
		Disconnects: disconnectsProcessed,
		Moves:       movementsProcessed,
		Attacks:     attacksProcessed,
	}
}
