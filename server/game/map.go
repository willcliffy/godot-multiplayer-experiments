package game

import (
	"errors"
	"math"
	"math/rand"

	pb "github.com/willcliffy/kilnwood-game-server/proto"
)

const (
	PLAYER_ATTACK_RANGE = 2

	spawn_range_x = 10
	spawn_range_z = 10
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

func (m *GameMap) addPlayer(p *Player) {
	for _, player := range m.players {
		if player.Id == p.Id {
			return
		}
	}

	m.players[p.Id] = p
}

func (m *GameMap) RemovePlayer(playerId uint64) {
	delete(m.players, playerId)
}

func (m *GameMap) SpawnPlayer(p *Player) {
	m.addPlayer(p)

	x := uint32(rand.Intn(spawn_range_x)) + 1
	z := uint32(rand.Intn(spawn_range_z)) + 1

	p.Location = &pb.Location{X: x, Z: z}
}

func (m *GameMap) DespawnPlayer() error {
	return errors.New("nyi")
}

func (m *GameMap) ApplyMovement(playerId uint64, target *pb.Location) error {
	player, ok := m.players[playerId]
	if !ok {
		return ErrPlayerNotInGame
	}

	player.Location = target

	return nil
}

func DistanceBetweenLocations(loc1, loc2 *pb.Location) float64 {
	return math.Sqrt(math.Pow(float64(loc2.X-loc1.X), 2) + math.Pow(float64(loc2.Z-loc1.Z), 2))
}

func (m *GameMap) InRangeToAttack(damage *pb.Damage) bool {
	dist := DistanceBetweenLocations(
		m.players[damage.SourcePlayerId].Location,
		m.players[damage.TargetPlayerId].Location)
	return dist < PLAYER_ATTACK_RANGE
}

func (m *GameMap) ApplyDamage(damage *pb.Damage) (killedTarget bool, err error) {
	target, ok := m.players[damage.TargetPlayerId]
	if !ok {
		return false, ErrPlayerNotInGame
	}

	source, ok := m.players[damage.SourcePlayerId]
	if !ok {
		return false, ErrPlayerNotInGame
	}

	return target.ApplyDamage(source.AttackDamage()), nil
}
