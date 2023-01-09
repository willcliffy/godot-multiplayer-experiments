package game

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/willcliffy/kilnwood-game-server/broadcast"
	gamemap "github.com/willcliffy/kilnwood-game-server/game/map"
	"github.com/willcliffy/kilnwood-game-server/game/objects"
	"github.com/willcliffy/kilnwood-game-server/game/objects/actions"
	"github.com/willcliffy/kilnwood-game-server/game/player"
	"github.com/willcliffy/kilnwood-game-server/util"
)

const gameTick = 3000 * time.Millisecond

type Game struct {
	id              uint64
	clock           *time.Ticker
	tick            int
	done            chan bool
	actionsQueued   []actions.Action
	movementsQueued map[uint64]actions.MoveAction
	gameMap         *gamemap.GameMap
	players         map[uint64]*player.Player
	broadcaster     broadcast.MessageBroadcaster
}

func NewGame(gameId uint64, broadcaster broadcast.MessageBroadcaster) *Game {
	return &Game{
		id:              gameId,
		done:            make(chan bool),
		actionsQueued:   make([]actions.Action, 0, 16),
		movementsQueued: make(map[uint64]actions.MoveAction),
		gameMap:         gamemap.NewGameMap(),
		broadcaster:     broadcaster,
	}
}

func (self Game) Id() string {
	return fmt.Sprint(self.id)
}

func (self *Game) Start() {
	self.clock = time.NewTicker(gameTick)
	go self.run()
}

func (self *Game) Stop() {
	self.done <- true
	self.clock.Stop()
}

func (self *Game) run() {
	for {
		select {
		case <-self.done:
			break
		case <-self.clock.C:
			processed := self.processQueue()

			payload, _ := json.Marshal(struct {
				Type   string
				Tick   int
				Events []actions.Action
			}{
				Type:   "tick",
				Tick:   self.tick,
				Events: processed,
			})

			err := self.broadcaster.BroadcastToGame(self.id, payload)
			if err != nil {
				log.Warn().Err(err).Msgf("failed to broadcast")
			}

			self.tick += 1
		}
	}
}

func (self *Game) OnMessageReceived(playerId uint64, message []byte) error {
	action, err := actions.ParseActionFromMessage(playerId, string(message))
	if err != nil {
		return err
	}

	if action.Type() != actions.ActionType_JoinGame {
		return self.QueueAction(playerId, action)
	}

	spawn, err := self.OnPlayerJoin(playerId, action.(*actions.JoinGameAction))
	if err != nil {
		log.Warn().Err(err).Msgf("err on player join")
	}

	msg := struct {
		Type     string
		PlayerId string
		Spawn    objects.Position
	}{
		Type:     "join-response",
		PlayerId: fmt.Sprint(playerId),
		Spawn:    spawn,
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = self.broadcaster.BroadcastToPlayer(playerId, payload)
	if err != nil {
		return err
	}

	msg.Type = "join-broadcast"
	payload, err = json.Marshal(msg)
	if err != nil {
		return err
	}

	err = self.broadcaster.BroadcastToGame(self.id, payload)
	if err != nil {
		return err
	}

	return nil
}

func (self *Game) OnPlayerJoin(playerId uint64, a *actions.JoinGameAction) (objects.Position, error) {
	// TODO - allow specifying team
	var team objects.Team
	if len(self.players) < 2 {
		team = objects.Team_Red
	} else {
		team = objects.Team_Blue
	}

	if _, ok := self.players[playerId]; ok {
		if playerPosition, ok := self.gameMap.GetPlayerPosition(playerId); ok {
			return playerPosition, nil
		}
	}

	player := player.NewPlayer(playerId, "", a.Class, team)

	err := self.gameMap.AddPlayer(player)
	if err != nil {
		return objects.Position{}, err
	}

	self.players[playerId] = player

	spawn, err := self.gameMap.SpawnPlayer(player)
	if err != nil {
		_ = self.gameMap.RemovePlayer(player.Id())
		delete(self.players, player.Id())
		return objects.Position{}, err
	}

	return spawn, nil
}

func (self *Game) QueueAction(playerId uint64, a actions.Action) error {
	log.Debug().Msgf("queuing action: %v", a)

	if a.Type() == actions.ActionType_Move {
		self.movementsQueued[playerId] = a.(actions.MoveAction)
		return nil
	}

	self.actionsQueued = append(self.actionsQueued, a)
	return nil
}

func (self *Game) DequeueAction(a actions.Action) {
	// inefficient but simple and preserves order
	for i, action := range self.actionsQueued {
		if action.Id() == a.Id() {
			self.actionsQueued = util.RemoveElementFromSlice(self.actionsQueued, i)
		}
	}
}

func (self *Game) processQueue() []actions.Action {
	processed := make([]actions.Action, 0)

	for _, action := range self.actionsQueued {
		switch action.Type() {
		case actions.ActionType_Attack:
			attackAction, ok := action.(actions.AttackAction)
			if !ok {
				log.Error().Msgf("Discarded action: could not cast %s to AttackAction", action.Id())
				continue
			}

			err := self.gameMap.ApplyAttack(attackAction)
			if err != nil {
				log.Error().Err(err).Msgf("Failed action: could not apply AttackAction")
				continue
			}

			processed = append(processed, attackAction)
		case actions.ActionType_CancelAction:
			log.Error().Msg("cancel action nyi")
			continue
		default:
			log.Error().Msgf("got bad action in process queue: %s", action.Id())
			continue
		}
	}

	for _, move := range self.movementsQueued {
		err := self.gameMap.ApplyMovement(&move)
		if err != nil {
			log.Error().Err(err).Msgf("Failed action: could not apply MoveAction")
			continue
		}

		processed = append(processed, move)
	}

	// reset actionQueue for the next tick
	self.actionsQueued = make([]actions.Action, 0, 16)

	return processed
}
