package player

import (
	pb "github.com/willcliffy/kilnwood-game-server/proto"
	"google.golang.org/protobuf/proto"
)

const (
	PlayerMaxHp          = 10
	PlayerRegenRateTicks = 100
)

type IPlayerCombat interface {
	IsAlive() bool
	GetCurrentHp() int32
	ApplyDamage(int32) bool
	Tick(uint64) []*pb.ClientAction
}

type PlayerCombat struct {
	super *Player
	// Stats
	Hp int32 // in HP units
	Ad int   // in HP units

	// State
	IsInCombat bool

	// Health Regen
	ticksToNextRegen int

	// Death and Respawn Timers
	isAlive        bool
	ticksToRespawn int
}

func NewPlayerCombat(player *Player) *PlayerCombat {
	return &PlayerCombat{
		super:            player,
		Hp:               PlayerMaxHp,
		Ad:               1,
		IsInCombat:       false,
		ticksToNextRegen: PlayerRegenRateTicks,
		isAlive:          true,
	}
}

func (p PlayerCombat) IsAlive() bool {
	return p.isAlive
}

func (p PlayerCombat) GetCurrentHp() int32 {
	return p.Hp
}

// Returns true if the damage kills the player
func (p *PlayerCombat) ApplyDamage(dmg int32) bool {
	p.Hp -= dmg
	if p.Hp <= 0 {
		p.isAlive = false
		p.ticksToRespawn = 500
		return true
	}
	return false
}

func (p *PlayerCombat) Tick(playerId uint64) []*pb.ClientAction {
	if p.Hp <= 0 {
		p.ticksToRespawn--
		if p.ticksToRespawn == 0 {
			actionBytes, _ := proto.Marshal(&pb.Respawn{
				Spawn: p.super.Spawn(nil),
			})
			return []*pb.ClientAction{
				{
					PlayerId: playerId,
					Payload:  actionBytes,
				},
			}
		}
	}

	// regen HP
	if p.Hp >= PlayerMaxHp {
		return nil
	}

	p.ticksToNextRegen--
	if p.ticksToNextRegen <= 0 {
		p.ticksToNextRegen = PlayerRegenRateTicks
		p.Hp++
		// todo - maybe one day add regen as clientaction?
	}

	return nil
}

func (p *PlayerCombat) Respawn(spawn *pb.Location) {
	p.Hp = PlayerMaxHp
	p.isAlive = true
}
