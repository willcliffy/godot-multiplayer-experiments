package game

import (
	"encoding/json"
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

	connectsQueued    map[uint64]*pb.Connect
	disconnectsQueued map[uint64]*pb.Disconnect
	movementsQueued   map[uint64]*pb.Move
	attacksQueued     map[uint64]*pb.Attack
}

func NewGame(gameId uint64) *Game {
	return &Game{
		id:   gameId,
		done: make(chan bool),

		gameMap: NewGameMap(),
		players: make(map[uint64]*Player),

		connectsQueued:    make(map[uint64]*pb.Connect),
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

			payload, _ := proto.Marshal(tickResult)

			err := g.broadcaster.BroadcastToGame(g.id, payload)
			if err != nil {
				log.Warn().Err(err).Msgf("failed to broadcast")
			}
		}
	}
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

func (g *Game) OnPlayerConnected(playerId uint64, a *pb.Connect) error {
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
			Spawn:    &p.Location,
		})

	}

	msg := pb.JoinGameResponse{
		PlayerId: p.Id,
		Color:    p.Color,
		Spawn:    &p.Location,
		Others:   playerList,
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = g.broadcaster.BroadcastToPlayer(playerId, payload)
	if err != nil {
		return err
	}

	msg.Others = nil
	payload, err = json.Marshal(msg)
	if err != nil {
		return err
	}

	err = g.broadcaster.BroadcastToGame(g.id, payload)
	if err != nil {
		return err
	}

	return nil
}
func (g *Game) processQueue() *pb.GameTick {
	var connectsProcessed []*pb.Connect
	var disconnectsProcessed []*pb.Disconnect
	var movementsProcessed []*pb.Move
	var attacksProcessed []*pb.Attack

	for _, connect := range g.connectsQueued {
		err := g.OnPlayerConnected(connect.PlayerId, connect)
		if err != nil {
			log.Err(err).Msgf("")
		}

		connectsProcessed = append(connectsProcessed, connect)
	}

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
	g.connectsQueued = make(map[uint64]*pb.Connect)
	g.disconnectsQueued = make(map[uint64]*pb.Disconnect)
	g.movementsQueued = make(map[uint64]*pb.Move)
	g.attacksQueued = make(map[uint64]*pb.Attack)

	return &pb.GameTick{
		Tick:        g.tick,
		Connects:    nil,
		Disconnects: nil,
		Moves:       nil,
		Attacks:     nil,
	}
}
