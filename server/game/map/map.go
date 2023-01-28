package gamemap

import (
	"errors"
	"math/rand"

	"github.com/willcliffy/kilnwood-game-server/game/objects"
	"github.com/willcliffy/kilnwood-game-server/game/player"
)

const (
	spawn_range_x = 11
	spawn_range_z = 11
)

var (
	ErrPlayerNotInGame = errors.New("player not in game")

	Spawn_RedOne  = objects.Location{X: 1, Z: 1}
	Spawn_RedTwo  = objects.Location{X: 9, Z: 1}
	Spawn_BlueOne = objects.Location{X: 1, Z: 9}
	Spawn_BlueTwo = objects.Location{X: 9, Z: 9}
)

type GameMap struct {
	players map[uint64]*player.Player
}

func NewGameMap() *GameMap {
	return &GameMap{
		players: make(map[uint64]*player.Player),
	}
}

func (m *GameMap) AddPlayer(p *player.Player) error {
	// if len(m.players) >= 2 {
	// 	return errors.New("game is full")
	// }

	for _, player := range m.players {
		if player.Id() == p.Id() {
			//return errors.New("player already in game")
			return nil
		}
	}

	m.players[p.Id()] = p
	p.SetPlayerState(objects.PlayerState_Vibing)
	return nil
}

func (m *GameMap) RemovePlayer(playerId uint64) error {
	delete(m.players, playerId)

	return nil
}

func (m *GameMap) SpawnPlayer(p *player.Player) (objects.Location, error) {
	x := rand.Intn(spawn_range_x)
	z := rand.Intn(spawn_range_z)

	spawn := objects.New2DLocation(x, z)

	p.SetTargetLocation(spawn)
	p.SetPlayerState(objects.PlayerState_Alive)

	return spawn, nil
}

func (m *GameMap) DespawnPlayer() error {
	return errors.New("nyi")
}

func (m *GameMap) ApplyMovement(movement *MoveAction) error {
	player, ok := m.players[movement.PlayerId]
	if !ok {
		return ErrPlayerNotInGame
	}

	player.Location = movement.ToLocation()
	return nil
}

func (m *GameMap) ApplyAttack(attack AttackAction) (int, error) {
	target, ok := m.players[attack.TargetPlayerId]
	if !ok {
		return 0, ErrPlayerNotInGame
	}

	source, ok := m.players[attack.SourcePlayerId]
	if !ok {
		return 0, ErrPlayerNotInGame
	}

	return target.Combat.ApplyDamage(source.Combat.AD), nil
}
