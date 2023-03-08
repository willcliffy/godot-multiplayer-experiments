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
	combat   IPlayerCombat
	movement IPlayerMovement
	crafting IPlayerCrafting
}

type PlayerTickResult struct {
	Respawn  bool
	Location *pb.Location
}

func NewPlayer(playerId uint64, color string) *Player {
	player := &Player{
		Id:    playerId,
		Color: color,
	}

	combat := NewPlayerCombat(player)
	player.combat = combat

	movement := NewPlayerMovement(player)
	player.movement = movement

	return player
}

func (p *Player) Tick() []*pb.ClientAction {
	return append(p.movement.Tick(p.Id), p.combat.Tick(p.Id)...)
}

func (p *Player) Spawn(location *pb.Location) *pb.Location {
	if location == nil {
		location = &pb.Location{
			X: int32(rand.Intn(SPAWN_RANGE_X)) + 1,
			Z: int32(rand.Intn(SPAWN_RANGE_Z)) + 1,
		}
	}

	p.movement.JumpToLocation(location)
	return location
}

func (p *Player) HandleMovement(move *pb.Move) {
	p.movement.SetPath(move.Path)
	p.movement.QueueAction(move.Queued)
}

func (p *Player) HandleCollection(collect *pb.Collect) {
	p.crafting.AddResource(collect.Type, 1)
}

func (p *Player) HandleBuild(build *pb.Build) bool {
	return p.crafting.Craft(build.Location)
}
