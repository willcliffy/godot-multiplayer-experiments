package game

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	gamemap "github.com/willcliffy/kilnwood-game-server/game/map"
	"github.com/willcliffy/kilnwood-game-server/game/objects"
	"github.com/willcliffy/kilnwood-game-server/game/objects/actions"
	"github.com/willcliffy/kilnwood-game-server/game/player"
	"github.com/willcliffy/kilnwood-game-server/util"
)

const gameTick = 10000 * time.Millisecond

type Game struct {
	clock       *time.Ticker
	tick        int
	done        chan bool
	actionQueue []actions.Action
	gameMap     *gamemap.GameMap
	players     []*player.Player
	broadcaster util.Broadcaster
}

func NewGame(b util.Broadcaster) *Game {
	return &Game{
		broadcaster: b,
	}
}

func (self *Game) Start() {
	self.clock = time.NewTicker(gameTick)
	self.done = make(chan bool)
	self.actionQueue = make([]actions.Action, 0, 16)
	self.gameMap = gamemap.NewGameMap()
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
				Tick   int
				Events []actions.Action
			}{
				Tick:   self.tick,
				Events: processed,
			})

			err := self.broadcaster.Broadcast("TODO - gameId", payload)
			if err != nil {
				log.Warn().Err(err).Msgf("failed to broadcast")
			}
			mapText := self.gameMap.DEBUG_DisplayGameMapText()
			for _, row := range mapText {
				fmt.Println(row)
			}

			self.tick += 1
		}
	}
}

func (self *Game) OnPlayerJoin(a *actions.JoinGameAction) {
	// TODO - allow specifying team
	var team objects.Team
	if len(self.players) < 2 {
		team = objects.Team_Red
	} else {
		team = objects.Team_Blue
	}

	player := player.NewPlayer(a.SourcePlayer(), a.Class, team)

	_ = self.gameMap.AddPlayer(player)
	self.players = append(self.players, player)

	err := self.gameMap.SpawnPlayer(player)
	if err != nil {
		_ = self.gameMap.RemovePlayer(player.Id())
		util.RemoveElementFromSlice(self.players, len(self.players)-1)
		return
	}
}

func (self *Game) QueueAction(a actions.Action) error {
	log.Debug().Msgf("queuing action: %v", a)
	// TODO - validate the action in the context of the game before appending?
	// This is a design decision, not a requirement - the action can just fail on the next tick
	self.actionQueue = append(self.actionQueue, a)
	return nil
}

func (self *Game) DequeueAction(a actions.Action) {
	// inefficient but simple and preserves order
	for i, action := range self.actionQueue {
		if action.Id() == a.Id() {
			util.RemoveElementFromSlice(self.actionQueue, i)
		}
	}
}

func (self *Game) processQueue() []actions.Action {
	processed := make([]actions.Action, 0)

	for _, action := range self.actionQueue {
		switch action.Type() {
		case actions.ActionType_Move:
			moveAction, ok := action.(*actions.MoveAction)
			if !ok {
				log.Error().Msgf("Discarded action: could not cast %s to MoveAction", action.Id())
				continue
			}

			err := self.gameMap.ApplyMovement(moveAction)
			if err != nil {
				log.Error().Err(err).Msgf("Failed action: could not apply MoveAction")
				continue
			}

			processed = append(processed, moveAction)
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

	// reset actionQueue for the next tick
	self.actionQueue = make([]actions.Action, 0, 16)

	return processed
}
