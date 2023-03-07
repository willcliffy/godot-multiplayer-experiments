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

	players map[uint64]*player.Player

	deathTimers map[uint64]int32

	connectsQueued    map[uint64]*pb.Connect
	disconnectsQueued map[uint64]*pb.Disconnect
	movementsQueued   map[uint64]*pb.Move
	attacksQueued     map[uint64]*pb.Attack
	damageQueued      map[uint64]*pb.Damage
	respawnsQueued    map[uint64]*pb.Respawn
}

func NewGame(gameId uint64, broadcaster broadcast.MessageBroadcaster) *Game {
	return &Game{
		id:          gameId,
		broadcaster: broadcaster,
		done:        make(chan bool),

		players: make(map[uint64]*player.Player),

		deathTimers: make(map[uint64]int32),

		connectsQueued:    make(map[uint64]*pb.Connect),
		disconnectsQueued: make(map[uint64]*pb.Disconnect),
		movementsQueued:   make(map[uint64]*pb.Move),
		attacksQueued:     make(map[uint64]*pb.Attack),
		damageQueued:      make(map[uint64]*pb.Damage),
		respawnsQueued:    make(map[uint64]*pb.Respawn),
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
				if result.Respawn {
					g.respawnsQueued[pId] = &pb.Respawn{
						PlayerId: pId,
						Spawn:    player.Movement.GetLocation(),
					}
				} else if result.Location != nil {
					g.movementsQueued[pId] = &pb.Move{
						PlayerId: pId,
						Path:     []*pb.Location{result.Location},
					}
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
	}

	p.Spawn(nil)

	var playerList []*pb.Connect

	for pId, p := range g.players {
		if p == nil || pId == playerId {
			continue
		}

		playerList = append(playerList, &pb.Connect{
			PlayerId: pId,
			Color:    p.Color,
			Spawn:    p.Movement.GetLocation(),
		})
	}

	msg := &pb.JoinGameResponse{
		PlayerId: p.Id,
		Color:    p.Color,
		Spawn:    p.Movement.GetLocation(),
		Others:   playerList,
	}

	g.connectsQueued[p.Id] = &pb.Connect{
		PlayerId: p.Id,
		Color:    p.Color,
		Spawn:    p.Movement.GetLocation(),
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

func (g *Game) OnActionReceived(playerId uint64, action *pb.ClientAction) {
	switch action.Type {
	case pb.ClientActionType_ACTION_PING:
		// nyi
	case pb.ClientActionType_ACTION_DISCONNECT:
		var disconnect pb.Disconnect
		err := proto.Unmarshal([]byte(action.Payload), &disconnect)
		if err != nil {
			log.Warn().Err(err).Msgf("failed to unmarshal disconnect")
		}

		g.disconnectsQueued[playerId] = &disconnect
	case pb.ClientActionType_ACTION_MOVE:
		var move pb.Move
		err := proto.Unmarshal([]byte(action.Payload), &move)
		if err != nil {
			log.Warn().Err(err).Msgf("failed to unmarshal move: %s", action.Payload)
		}

		g.players[playerId].Movement.SetPath(move.Path)
	case pb.ClientActionType_ACTION_ATTACK:
		var attack pb.Attack
		err := proto.Unmarshal([]byte(action.Payload), &attack)
		if err != nil {
			log.Warn().Err(err).Msgf("failed to unmarshal attack")
		}

		g.attacksQueued[playerId] = &attack
	case pb.ClientActionType_ACTION_DAMAGE:
		var damage pb.Damage
		err := proto.Unmarshal([]byte(action.Payload), &damage)
		if err != nil {
			log.Warn().Err(err).Msgf("failed to unmarshal damage")
		}

		g.damageQueued[playerId] = &damage
	default:
		log.Error().Msgf("could not unmarshal unknown action type: %v", action.Type)
	}
}

func (g *Game) processQueue() *pb.GameTick {
	// TODO - 🍝
	actions := make([]*pb.GameTickAction, 0)

	for _, connect := range g.connectsQueued {
		connectBytes, _ := proto.Marshal(connect)
		actions = append(actions, &pb.GameTickAction{
			Type:  pb.ClientActionType_ACTION_CONNECT,
			Value: connectBytes,
		})
	}

	for _, disconnect := range g.disconnectsQueued {
		g.OnPlayerDisconnected(disconnect.PlayerId)
		disconnectBytes, _ := proto.Marshal(disconnect)
		actions = append(actions, &pb.GameTickAction{
			Type:  pb.ClientActionType_ACTION_DISCONNECT,
			Value: disconnectBytes,
		})
	}

	for _, move := range g.movementsQueued {
		moveBytes, _ := proto.Marshal(move)
		actions = append(actions, &pb.GameTickAction{
			Type:  pb.ClientActionType_ACTION_MOVE,
			Value: moveBytes,
		})
	}

	// Temporarily disable combat
	// for _, attack := range g.attacksQueued {
	// 	err := g.gameMap.ApplyMovement(attack.TargetPlayerId, attack.TargetPlayerLocation)
	// 	if err != nil {
	// 		log.Error().Err(err).Msgf("Failed action: could not apply MoveAction")
	// 		continue
	// 	}

	// 	attackBytes, _ := proto.Marshal(attack)
	// 	actions = append(actions, &pb.GameTickAction{
	// 		Type:  pb.ClientActionType_ACTION_ATTACK,
	// 		Value: attackBytes,
	// 	})
	// }

	for _, damage := range g.damageQueued {
		killedTarget := g.players[damage.TargetPlayerId].Combat.ApplyDamage(damage.DamageDealt)
		damageBytes, _ := proto.Marshal(damage)
		actions = append(actions, &pb.GameTickAction{
			Type:  pb.ClientActionType_ACTION_DAMAGE,
			Value: damageBytes,
		})

		if killedTarget {
			g.deathTimers[damage.TargetPlayerId] = 50
			deathBytes, _ := proto.Marshal(&pb.Death{
				PlayerId: damage.TargetPlayerId,
				Location: g.players[damage.TargetPlayerId].Movement.GetLocation(),
			})
			actions = append(actions, &pb.GameTickAction{
				Type:  pb.ClientActionType_ACTION_DEATH,
				Value: deathBytes,
			})
		}
	}

	for pId, respawn := range g.respawnsQueued {
		delete(g.deathTimers, pId)
		player := g.players[pId]
		player.Spawn(nil)
		respawnBytes, _ := proto.Marshal(respawn)
		actions = append(actions, &pb.GameTickAction{
			Type:  pb.ClientActionType_ACTION_RESPAWN,
			Value: respawnBytes,
		})
	}

	// reset actionQueue for the next tick
	g.connectsQueued = make(map[uint64]*pb.Connect)
	g.disconnectsQueued = make(map[uint64]*pb.Disconnect)
	g.movementsQueued = make(map[uint64]*pb.Move)
	g.attacksQueued = make(map[uint64]*pb.Attack)
	g.damageQueued = make(map[uint64]*pb.Damage)
	g.respawnsQueued = make(map[uint64]*pb.Respawn)

	tick := &pb.GameTick{
		Tick:    g.tick,
		Actions: actions,
	}

	return tick
}
