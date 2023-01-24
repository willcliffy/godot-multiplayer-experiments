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
)

const gameTick = 100 * time.Millisecond

type Game struct {
	id              uint64
	clock           *time.Ticker
	tick            uint8
	done            chan bool
	actionsQueued   []actions.Action
	movementsQueued map[uint64]*actions.MoveAction
	gameMap         *gamemap.GameMap
	players         map[uint64]*player.Player
	broadcaster     broadcast.MessageBroadcaster
}

func NewGame(gameId uint64, broadcaster broadcast.MessageBroadcaster) *Game {
	return &Game{
		id:              gameId,
		done:            make(chan bool),
		actionsQueued:   make([]actions.Action, 0, 16),
		movementsQueued: make(map[uint64]*actions.MoveAction),
		gameMap:         gamemap.NewGameMap(),
		players:         make(map[uint64]*player.Player),
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
			self.tick += 1
			processed := self.processQueue()
			if len(processed) == 0 {
				continue
			}

			payload, _ := json.Marshal(struct {
				Type   string
				Tick   uint8
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
		}
	}
}

// This satisfies the `MessageReceiver` interface, which the MessageBroker uses
func (self *Game) OnMessageReceived(playerId uint64, message []byte) error {
	log.Debug().Msgf("message received from player '%v': %v", playerId, string(message))

	action, err := actions.ParseActionFromMessage(playerId, string(message))
	if err != nil {
		return err
	}

	if action.Type() == actions.ActionType_JoinGame {
		return self.onPlayerJoin(playerId, action.(*actions.JoinGameAction))
	}

	return self.QueueAction(playerId, action)
}

// This satisfies the `MessageReceiver` interface, which the MessageBroker uses
func (self *Game) OnPlayerDisconnected(playerId uint64) {
	delete(self.players, playerId)
	err := self.broadcaster.BroadcastToGame(self.id, []byte(fmt.Sprintf("d:%d", playerId)))
	if err != nil {
		log.Warn().Err(err).Msgf("failed to broadcast")
	}
}

func (self *Game) onPlayerJoin(playerId uint64, a *actions.JoinGameAction) error {
	// TODO - allow specifying team
	// team := objects.Team_Red
	// if len(self.players) == 0 {
	// 	team = objects.Team_Blue
	// }

	// for now, give everyone a random color
	color := objects.RandomTeamColor()

	p, playerInGame := self.players[playerId]
	if !playerInGame {
		p = player.NewPlayer(playerId, "", a.Class, color)
		self.players[playerId] = p
	}

	err := self.gameMap.AddPlayer(p)
	if err != nil {
		delete(self.players, p.Id())
		return err
	}

	_, err = self.gameMap.SpawnPlayer(p)
	if err != nil {
		delete(self.players, p.Id())
		return err
	}

	type PlayerListEntry struct {
		PlayerId string
		Color    objects.TeamColor
		Spawn    objects.Location
		//Team     objects.Team
	}

	playerList := make([]PlayerListEntry, 0, 2)

	for pId, p := range self.players {
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

	msg := struct {
		Type     string
		PlayerId string
		Color    objects.TeamColor
		Spawn    objects.Location
		Others   []PlayerListEntry
		// Team     objects.Team
	}{
		Type:     "join-response",
		PlayerId: fmt.Sprint(playerId),
		Color:    color,
		Spawn:    p.Location,
		Others:   playerList,
		// Team:     team
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
	msg.Others = nil
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

func (self *Game) QueueAction(playerId uint64, a actions.Action) error {
	log.Debug().Msgf("queuing action: %v", a)

	if a.Type() == actions.ActionType_Move {
		self.movementsQueued[playerId] = a.(*actions.MoveAction)
		return nil
	}

	self.actionsQueued = append(self.actionsQueued, a)
	return nil
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

			dmg, err := self.gameMap.ApplyAttack(attackAction)
			if err != nil {
				log.Error().Err(err).Msgf("Failed action: could not apply AttackAction")
				continue
			}

			attackAction.SetDamageDealt(dmg)

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
		err := self.gameMap.ApplyMovement(move)
		if err != nil {
			log.Error().Err(err).Msgf("Failed action: could not apply MoveAction")
			continue
		}

		processed = append(processed, move)
	}

	// reset actionQueue for the next tick
	self.actionsQueued = make([]actions.Action, 0, 16)
	self.movementsQueued = make(map[uint64]*actions.MoveAction)

	return processed
}
