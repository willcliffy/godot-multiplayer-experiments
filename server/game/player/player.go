package player

import (
	"math/rand"

	pb "github.com/willcliffy/kilnwood-game-server/proto"
)

const (
	PLAYER_MOVE_SPEED   = 5
	PLAYER_ATTACK_RANGE = 2

	SPAWN_RANGE_X = 10
	SPAWN_RANGE_Z = 10
)

type Player struct {
	Id    uint64
	Color string

	// Expose these interfaces
	Combat   IPlayerCombat
	Movement IPlayerMovement

	// Allow calling underlying methods here though
	combat   *PlayerCombat
	movement *PlayerMovement
}

type PlayerTickResult struct {
	Respawn  bool
	Location *pb.Location
}

func NewPlayer(playerId uint64, color string) *Player {
	combat := NewPlayerCombat()
	movement := NewPlayerMovement()
	return &Player{
		Id:       playerId,
		Color:    color,
		Combat:   combat,
		Movement: movement,
		combat:   combat,
		movement: movement,
	}
}

func (p *Player) Tick() *PlayerTickResult {
	return &PlayerTickResult{
		Respawn:  p.combat.Tick(),
		Location: p.movement.Tick(),
	}
}

func (p *Player) Spawn(location *pb.Location) {
	if location == nil {
		location = &pb.Location{
			X: int32(rand.Intn(SPAWN_RANGE_X)) + 1,
			Z: int32(rand.Intn(SPAWN_RANGE_Z)) + 1,
		}
	}

	p.movement.location = location
}
