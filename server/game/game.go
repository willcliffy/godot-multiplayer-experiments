package game

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/willcliffy/kilnwood-game-server/broadcast"
	gamemap "github.com/willcliffy/kilnwood-game-server/game/map"
	"github.com/willcliffy/kilnwood-game-server/game/objects"
	"github.com/willcliffy/kilnwood-game-server/game/player"
	"github.com/willcliffy/kilnwood-game-server/pb"
	"google.golang.org/protobuf/proto"
)

// 240 bpm or 4 bps
const gameTick = 1000 / 4 * time.Millisecond

type Game struct {
	id          uint64
	broadcaster broadcast.MessageBroadcaster
	done        chan bool

	clock *time.Ticker
	tick  uint32

	gameMap *gamemap.GameMap
	players map[uint64]*player.Player

	connectsQueued    map[uint64]*pb.ConnectAction
	disconnectsQueued map[uint64]*pb.DisconnectAction
	movementsQueued   map[uint64]*pb.MoveAction
	attacksQueued     map[uint64]*pb.AttackAction
}

func NewGame(gameId uint64, broadcaster broadcast.MessageBroadcaster) *Game {
	return &Game{
		id:          gameId,
		broadcaster: broadcaster,
		done:        make(chan bool),

		gameMap: gamemap.NewGameMap(),
		players: make(map[uint64]*player.Player),

		connectsQueued:    make(map[uint64]*pb.ConnectAction),
		disconnectsQueued: make(map[uint64]*pb.DisconnectAction),
		movementsQueued:   make(map[uint64]*pb.MoveAction),
		attacksQueued:     make(map[uint64]*pb.AttackAction),
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
			processed := g.processQueue()
			if len(processed) == 0 {
				continue
			}

			body := &pb.GameTick{
				Tick:        g.tick,
				Connects:    nil,
				Disconnects: nil,
				Moves:       nil,
				Attacks:     nil,
			}

			payload, _ := proto.Marshal(body)

			err := g.broadcaster.BroadcastToGame(g.id, payload)
			if err != nil {
				log.Warn().Err(err).Msgf("failed to broadcast")
			}
		}
	}
}

// This satisfies the `MessageReceiver` interface, which the MessageBroker uses
func (g *Game) OnMessageReceived(playerId uint64, message []byte) {
	log.Debug().Msgf("message received from player '%v': %v", playerId, string(message))

	action, err := actions.ParseActionFromMessage(playerId, string(message))
	if err != nil {
		log.Warn().Err(err).Msgf("failed to parse action from message")
	}

	switch action.Type() {
	case actions.ActionType_JoinGame:
		err := g.onPlayerJoin(playerId, action.(*actions.JoinGameAction))
		if err != nil {
			log.Warn().Err(err).Msgf("failed to join game")
		}
	case actions.ActionType_Move:
		g.movementsQueued[playerId] = action.(*actions.MoveAction)
	case actions.ActionType_Disconnect:
		g.actionsQueued = append(g.actionsQueued, action)
	default:
		g.actionsQueued = append(g.actionsQueued, action)
	}
}

// This satisfies the `MessageReceiver` interface, which the MessageBroker uses
func (g *Game) OnPlayerDisconnected(playerId uint64) {
	delete(g.players, playerId)
	action := actions.DisconnectAction{PlayerId: playerId}
	g.actionsQueued = append(g.actionsQueued, action)
	log.Info().Msgf("queued disconnectw")
}

func (g *Game) Close() {
	action := actions.DisconnectAction{All: true}
	g.actionsQueued = append(g.actionsQueued, action)
}

func (g *Game) onPlayerJoin(playerId uint64, a *actions.JoinGameAction) error {
	p, playerInGame := g.players[playerId]
	if !playerInGame {
		// TODO - allow specifying team. for now, give everyone a random color
		color := objects.RandomTeamColor()
		p = player.NewPlayer(playerId, "", a.Class, color)
		g.players[playerId] = p
	}

	err := g.gameMap.AddPlayer(p)
	if err != nil {
		delete(g.players, p.Id())
		return err
	}

	_, err = g.gameMap.SpawnPlayer(p)
	if err != nil {
		delete(g.players, p.Id())
		return err
	}

	playerList := make([]PlayerListEntry, 0, 2)

	for pId, p := range g.players {
		if p == nil || pId == playerId {
			continue
		}

		playerList = append(playerList, PlayerListEntry{
			PlayerId: fmt.Sprint(pId),
			Color:    p.Color,
			Spawn:    p.Location,
			//Team:     p.Team,
		})

	}

	msg := PlayerJoinResponse{
		Type:     "join-response",
		PlayerId: fmt.Sprint(p.Id()),
		Color:    p.Color,
		Spawn:    p.Location,
		Others:   playerList,
		// Team:     team
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = g.broadcaster.BroadcastToPlayer(playerId, payload)
	if err != nil {
		return err
	}

	msg.Type = "join-broadcast"
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
func (g *Game) processQueue() []actions.Action {
	processed := make([]actions.Action, 0)

	for _, action := range g.actionsQueued {
		switch action.Type() {
		case actions.ActionType_Attack:
			attackAction, ok := action.(actions.AttackAction)
			if !ok {
				log.Error().Msgf("Discarded action: could not cast %s to AttackAction", action.Id())
				continue
			}

			dmg, err := g.gameMap.ApplyAttack(attackAction)
			if err != nil {
				log.Error().Err(err).Msgf("Failed action: could not apply AttackAction")
				continue
			}

			attackAction.SetDamageDealt(dmg)

			processed = append(processed, attackAction)
		case actions.ActionType_Disconnect:
			disconnectAction, ok := action.(actions.DisconnectAction)
			if !ok {
				log.Error().Msgf("Discarded action: could not cast %s to DisconnectAction", action.Id())
				continue
			}

			log.Info().Msgf("calling player disconnected")
			processed = append(processed, disconnectAction)
		default:
			log.Error().Msgf("got bad action in process queue: %s", action.Id())
			continue
		}
	}

	for _, move := range g.movementsQueued {
		err := g.gameMap.ApplyMovement(move)
		if err != nil {
			log.Error().Err(err).Msgf("Failed action: could not apply MoveAction")
			continue
		}

		processed = append(processed, move)
	}

	// reset actionQueue for the next tick
	g.actionsQueued = make([]actions.Action, 0, 16)
	g.movementsQueued = make(map[uint64]*actions.MoveAction)

	return processed
}
