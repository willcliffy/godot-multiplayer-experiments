package gamemap

import (
	"errors"

	"github.com/willcliffy/kilnwood-game-server/game/objects"
	"github.com/willcliffy/kilnwood-game-server/game/objects/actions"
	"github.com/willcliffy/kilnwood-game-server/game/player"
	"github.com/willcliffy/kilnwood-game-server/util"
)

const (
	gamemap_x = 11
	gamemap_y = 11
)

var (
	Spawn_RedOne  = objects.Location{X: 1, Z: 1}
	Spawn_RedTwo  = objects.Location{X: 9, Z: 1}
	Spawn_BlueOne = objects.Location{X: 1, Z: 9}
	Spawn_BlueTwo = objects.Location{X: 9, Z: 9}
)

type GameMap struct {
	tiles           [][]Tile
	players         []*player.Player
	playerLocations map[uint64]objects.Location
}

func NewGameMap() *GameMap {
	tiles := make([][]Tile, gamemap_x)
	for i := 0; i < gamemap_x; i++ {
		tiles[i] = make([]Tile, gamemap_y)
		for j := 0; j < gamemap_y; j++ {
			tiles[i][j] = Tile{}
		}
	}

	return &GameMap{
		tiles:           tiles,
		players:         make([]*player.Player, 0, 4),
		playerLocations: make(map[uint64]objects.Location),
	}
}

func (self *GameMap) AddPlayer(p *player.Player) error {
	// if len(self.players) >= 2 {
	// 	return errors.New("game is full")
	// }

	for _, player := range self.players {
		if player.Id() == p.Id() {
			//return errors.New("player already in game")
			return nil
		}
	}

	self.players = append(self.players, p)
	p.SetPlayerState(objects.PlayerState_Vibing)
	return nil
}

func (self *GameMap) RemovePlayer(playerId uint64) error {
	for i, player := range self.players {
		if player.Id() == playerId {
			util.RemoveElementFromSlice(self.players, i)
			delete(self.playerLocations, playerId)
		}
	}

	return nil
}

func (self *GameMap) SpawnPlayer(p *player.Player) (objects.Location, error) {
	var spawn objects.Location
	if p.Team == objects.Team_Red {
		spawn = Spawn_RedOne
	} else {
		spawn = Spawn_BlueOne
	}

	for _, player := range self.players {
		if player.Team == p.Team {
			if p.Team == objects.Team_Red {
				spawn = Spawn_RedTwo
			} else {
				spawn = Spawn_BlueTwo
			}
		}
	}

	self.playerLocations[p.Id()] = spawn
	p.SetTargetLocation(spawn)
	p.SetPlayerState(objects.PlayerState_Alive)

	return spawn, nil
}

func (self *GameMap) DespawnPlayer() error {
	return errors.New("nyi")
}

func (self *GameMap) ApplyMovement(movement *actions.MoveAction) error {
	self.playerLocations[movement.PlayerId] = movement.ToLocation()
	return nil
}

func (self *GameMap) ApplyAttack(attack actions.AttackAction) error {
	return errors.New("nyi")
}
