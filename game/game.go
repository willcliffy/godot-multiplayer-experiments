package game

import (
	"time"

	"github.com/rs/zerolog/log"
	"github.com/willcliffy/kilnwood-game-server/game/actions"
	"github.com/willcliffy/kilnwood-game-server/util"
)

const gameTick = 2500 * time.Millisecond

type Game struct {
	clock       *time.Ticker
	done        chan bool
	actionQueue []actions.Action
}

func NewGame() *Game {
	return &Game{}
}

func (g *Game) Start() {
	g.clock = time.NewTicker(gameTick)
	g.done = make(chan bool)
	g.actionQueue = make([]actions.Action, 0, 16)
	go g.run()
}

func (g *Game) Stop() {
	g.clock.Stop()
	g.done <- true
}

func (g *Game) run() {
	for {
		select {
		case <-g.done:
			return
		case _ = <-g.clock.C:
			log.Debug().Msgf("Tick. Queue length: %v", len(g.actionQueue))
			g.ProcessQueue()
		}
	}
}

func (g *Game) QueueAction(a actions.Action) {
	g.actionQueue = append(g.actionQueue, a)
}

func (g *Game) DequeueAction(a actions.Action) {
	// inefficient, but simple and preserves order
	for i, action := range g.actionQueue {
		if action.ID() == a.ID() {
			util.RemoveElementFromSlice(g.actionQueue, i)
		}
	}
}

func (g *Game) ProcessQueue() {
	for _, action := range g.actionQueue {
		log.Info().Msgf("Processed action: %v", action)
		// TODO - actually process event
	}

	g.actionQueue = make([]actions.Action, 0, 16)
}
