package gamemap

import (
	"errors"
	"fmt"

	"github.com/willcliffy/kilnwood-game-server/game/objects"
	"github.com/willcliffy/kilnwood-game-server/game/objects/actions"
	"github.com/willcliffy/kilnwood-game-server/game/player"
)

const (
	gamemap_x = 10
	gamemap_y = 10
)

var (
	Spawn_RedOne  = objects.Position{X: 10, Y: 10}
	Spawn_RedTwo  = objects.Position{X: 90, Y: 10}
	Spawn_BlueOne = objects.Position{X: 10, Y: 90}
	Spawn_BlueTwo = objects.Position{X: 90, Y: 90}
)

type GameMap struct {
	tiles           [][]Tile
	players         []*player.Player
	playerLocations map[string]objects.Position
}

func NewGameMap() *GameMap {
	tiles := make([][]Tile, gamemap_x)
	for i := 0; i < gamemap_x; i++ {
		tiles[i] = make([]Tile, gamemap_y)
		for j := 0; j < gamemap_y; j++ {
			tiles[i][j] = Tile{
				{" NW", " N ", "NE "},
				{" W ", "   ", " E "},
				{" SW", " S ", "SE "},
			}
		}
	}

	fmt.Printf("%v\n", len(tiles))

	return &GameMap{
		tiles:           tiles,
		players:         make([]*player.Player, 0, 4),
		playerLocations: make(map[string]objects.Position),
	}
}

func (self *GameMap) AddPlayer(p *player.Player) error {
	// TODO - check if in here already. check if too many playrs
	self.players = append(self.players, p)
	return nil
}

func (self *GameMap) RemovePlayer() error {
	return errors.New("nyi")
}

func (self *GameMap) SpawnPlayer(p *player.Player) error {
	var spawn objects.Position
	if p.Team == objects.Team_Red {
		spawn = Spawn_RedOne
	} else {
		spawn = Spawn_BlueOne
	}

	for _, player := range self.players {
		if player.Name == p.Name {
			return errors.New("player already spawned")
		}

		if player.Team == p.Team {
			if p.Team == objects.Team_Red {
				spawn = Spawn_RedTwo
			} else {
				spawn = Spawn_BlueTwo
			}
		}
	}

	self.playerLocations[p.ID()] = spawn

	return nil
}

func (self *GameMap) DespawnPlayer() error {
	return errors.New("nyi")
}

func (self *GameMap) ApplyMovement(movement actions.MoveAction) error {
	return errors.New("nyi")
}

func (self *GameMap) ApplyAttack(attack actions.AttackAction) error {
	return errors.New("nyi")
}

func (self GameMap) DEBUG_CopyTiles() [][]Tile {
	tiles := make([][]Tile, gamemap_x)
	for i := 0; i < gamemap_x; i++ {
		tiles[i] = make([]Tile, gamemap_y)
	}

	for i, tileRow := range self.tiles {
		for j, tile := range tileRow {
			tiles[i][j] = tile.Copy()
		}
	}

	return tiles
}

func (self GameMap) DEBUG_DisplayGameMapText() []string {
	tiles := self.DEBUG_CopyTiles()
	for _, player := range self.players {
		loc := self.playerLocations[player.ID()]
		tiles[loc.X][loc.Y] = player.DEBUG_Tile()
	}

	var res []string
	for _, tileRow := range tiles {
		var one, two, three string
		for _, tile := range tileRow {
			for i := 0; i < 3; i++ {
				one += tile[0][i]
				two += tile[1][i]
				three += tile[2][i]
			}
		}

		res = append(res, one)
		res = append(res, two)
		res = append(res, three)
	}

	return res
}
