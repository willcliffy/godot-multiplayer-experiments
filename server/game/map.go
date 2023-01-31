package game

import (
	"errors"
	"math/rand"

	pb "github.com/willcliffy/kilnwood-game-server/proto"
)

const (
	spawn_range_x = 11
	spawn_range_z = 11
)

var (
	ErrPlayerNotInGame = errors.New("player not in game")

	Spawn_RedOne  = pb.Location{X: 1, Z: 1}
	Spawn_RedTwo  = pb.Location{X: 9, Z: 1}
	Spawn_BlueOne = pb.Location{X: 1, Z: 9}
	Spawn_BlueTwo = pb.Location{X: 9, Z: 9}
)

type GameMap struct {
	players map[uint64]*Player
}

func NewGameMap() *GameMap {
	return &GameMap{
		players: make(map[uint64]*Player),
	}
}

func (m *GameMap) AddPlayer(p *Player) error {
	// if len(m.players) >= 2 {
	// 	return errors.New("game is full")
	// }

	for _, player := range m.players {
		if player.Id == p.Id {
			//return errors.New("player already in game")
			return nil
		}
	}

	m.players[p.Id] = p
	return nil
}

func (m *GameMap) RemovePlayer(playerId uint64) error {
	delete(m.players, playerId)

	return nil
}

func (m *GameMap) SpawnPlayer(p *Player) (*pb.Location, error) {
	x := uint32(rand.Intn(spawn_range_x))
	z := uint32(rand.Intn(spawn_range_z))

	spawn := &pb.Location{X: x, Z: z}

	p.Location = spawn

	return spawn, nil
}

func (m *GameMap) DespawnPlayer() error {
	return errors.New("nyi")
}

func (m *GameMap) ApplyMovement(movement *pb.Move) error {
	player, ok := m.players[movement.PlayerId]
	if !ok {
		return ErrPlayerNotInGame
	}

	player.Location = movement.Target

	return nil
}

func (m *GameMap) ApplyAttack(attack *pb.Attack) (int, error) {
	target, ok := m.players[attack.TargetPlayerId]
	if !ok {
		return 0, ErrPlayerNotInGame
	}

	source, ok := m.players[attack.SourcePlayerId]
	if !ok {
		return 0, ErrPlayerNotInGame
	}

	return target.ApplyDamage(source.AttackDamage()), nil
}
