package game

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/willcliffy/kilnwood-game-server/broadcast"
	"github.com/willcliffy/kilnwood-game-server/constants"
	"github.com/willcliffy/kilnwood-game-server/game/player"
	pb "github.com/willcliffy/kilnwood-game-server/proto"
	"google.golang.org/protobuf/proto"
)

type Game struct {
	id          uint64
	broadcaster broadcast.MessageBroadcaster
	done        chan bool

	clock *time.Ticker
	tick  uint32

	players       map[uint64]*player.Player
	actionsQueued map[uint64]map[pb.ClientActionType]*pb.ClientAction
}

func NewGame(gameId uint64, broadcaster broadcast.MessageBroadcaster) *Game {
	return &Game{
		id:          gameId,
		broadcaster: broadcaster,
		done:        make(chan bool),

		players:       make(map[uint64]*player.Player),
		actionsQueued: make(map[uint64]map[pb.ClientActionType]*pb.ClientAction),
	}
}

func (g Game) Id() string {
	return fmt.Sprint(g.id)
}

func (g *Game) Start() {
	g.clock = time.NewTicker(constants.GAME_TICK_DURATION)
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

			for pId, player := range g.players {
				result := player.Tick()
				for _, result := range result {
					g.actionsQueued[pId][result.Type] = result
				}
			}

			tickResult := g.processQueue()
			if tickResult == nil || len(tickResult.Actions) == 0 {
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
		color := player.RandomTeamColor()
		p = player.NewPlayer(playerId, color)
		g.players[playerId] = p
		g.actionsQueued[playerId] = make(map[pb.ClientActionType]*pb.ClientAction)
	}

	location := p.Spawn(nil)

	var playerList []*pb.ClientAction

	for pId, p := range g.players {
		if p == nil || pId == playerId {
			continue
		}

		actionBytes, _ := proto.Marshal(&pb.Connect{
			Color: p.Color,
			Spawn: location,
		})

		playerList = append(playerList, &pb.ClientAction{
			PlayerId: pId,
			Payload:  actionBytes,
		})
	}

	msg := &pb.JoinGameResponse{
		PlayerId: p.Id,
		Color:    p.Color,
		Spawn:    location,
		Others:   playerList,
	}

	actionBytes, _ := proto.Marshal(&pb.Connect{
		Color: p.Color,
		Spawn: location,
	})

	g.actionsQueued[p.Id][pb.ClientActionType_ACTION_CONNECT] = &pb.ClientAction{
		PlayerId: p.Id,
		Payload:  actionBytes,
	}

	g.broadcaster.OnPlayerJoinGame(g.id, playerId, msg)
	return nil
}

func (g *Game) OnPlayerDisconnected(playerId uint64) {
	delete(g.players, playerId)
	actionBytes, _ := proto.Marshal(&pb.Disconnect{})
	g.actionsQueued[playerId][pb.ClientActionType_ACTION_DISCONNECT] = &pb.ClientAction{
		PlayerId: playerId,
		Payload:  actionBytes,
	}
	log.Info().Msgf("queued disconnect for player %d", playerId)
}

func (g *Game) Close() {
	for pId := range g.players {
		actionBytes, _ := proto.Marshal(&pb.Disconnect{})
		g.actionsQueued[pId][pb.ClientActionType_ACTION_DISCONNECT] = &pb.ClientAction{
			PlayerId: pId,
			Payload:  actionBytes,
		}
	}
}

func (g *Game) OnActionReceived(playerId uint64, action *pb.ClientAction) {
	if action.Type == pb.ClientActionType_ACTION_MOVE {
		var move pb.Move
		err := proto.Unmarshal([]byte(action.Payload), &move)
		if err != nil {
			log.Warn().Err(err).Msgf("failed to unmarshal move: %s", action.Payload)
		}

		g.players[playerId].HandleMovement(&move)
		// TODO - handle queued actions
		return
	}
	g.actionsQueued[playerId][action.Type] = action
	// switch action.Type {
	// case pb.ClientActionType_ACTION_PING:
	// 	// nyi
	// case pb.ClientActionType_ACTION_DISCONNECT:
	// 	var disconnect pb.Disconnect
	// 	err := proto.Unmarshal([]byte(action.Payload), &disconnect)
	// 	if err != nil {
	// 		log.Warn().Err(err).Msgf("failed to unmarshal disconnect")
	// 	}

	// 	g.disconnectsQueued[playerId] = &disconnect
	// case pb.ClientActionType_ACTION_MOVE:

	// case pb.ClientActionType_ACTION_ATTACK:
	// 	var attack pb.Attack
	// 	err := proto.Unmarshal([]byte(action.Payload), &attack)
	// 	if err != nil {
	// 		log.Warn().Err(err).Msgf("failed to unmarshal attack")
	// 	}

	// 	g.attacksQueued[playerId] = &attack
	// case pb.ClientActionType_ACTION_DAMAGE:
	// 	var damage pb.Damage
	// 	err := proto.Unmarshal([]byte(action.Payload), &damage)
	// 	if err != nil {
	// 		log.Warn().Err(err).Msgf("failed to unmarshal damage")
	// 	}

	// 	g.damageQueued[playerId] = &damage
	// default:
	// 	log.Error().Msgf("could not unmarshal unknown action type: %v", action.Type)
	// }
}

func (g *Game) processQueue() *pb.GameTick {
	// TODO - 🍝
	processedActions := make([]*pb.ClientAction, 0)

	for playerId, actions := range g.actionsQueued {
		for actionType, action := range actions {
			processedActions = append(processedActions, action)
			switch actionType {
			case pb.ClientActionType_ACTION_DISCONNECT:
				g.OnPlayerDisconnected(playerId)
			case pb.ClientActionType_ACTION_COLLECT:
				var collect pb.Collect
				_ = proto.Unmarshal(action.Payload, &collect)
				g.players[playerId].HandleCollect(&collect)
			case pb.ClientActionType_ACTION_BUILD:
				var build pb.Build
				_ = proto.Unmarshal(action.Payload, &build)
				g.players[playerId].HandleBuild(&build)
			}
		}
		g.actionsQueued[playerId] = make(map[pb.ClientActionType]*pb.ClientAction)
	}

	tick := &pb.GameTick{
		Tick:    g.tick,
		Actions: processedActions,
	}

	return tick
}
