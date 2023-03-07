package player

import pb "github.com/willcliffy/kilnwood-game-server/proto"

const (
	PlayerMaxHp          = 10
	PlayerRegenRateTicks = 100
)

type IPlayerCombat interface {
	IsAlive() bool
	GetCurrentHp() int32
	ApplyDamage(int32) bool
}

type PlayerCombat struct {
	// Stats
	Hp int32 // in HP units
	Ad int   // in HP units

	// Health Regen
	ticksToNextRegen int

	// Death and Respawn Timers
	Alive          bool
	ticksToRespawn int
}

func NewPlayerCombat() *PlayerCombat {
	return &PlayerCombat{
		Hp:               PlayerMaxHp,
		Ad:               1,
		ticksToNextRegen: PlayerRegenRateTicks,
		Alive:            true,
	}
}

func (p PlayerCombat) IsAlive() bool {
	return p.Alive
}

func (p PlayerCombat) GetCurrentHp() int32 {
	return p.Hp
}

// Returns true if the damage kills the player
func (p *PlayerCombat) ApplyDamage(dmg int32) bool {
	p.Hp -= dmg
	if p.Hp <= 0 {
		p.Alive = false
		p.ticksToRespawn = 500
		return true
	}
	return false
}

// Returns true if
func (p *PlayerCombat) Tick() bool {
	if p.Hp <= 0 {
		p.ticksToRespawn--
		return p.ticksToRespawn == 0
	}

	// regen HP
	if p.Hp >= PlayerMaxHp {
		return false
	}

	p.ticksToNextRegen--
	if p.ticksToNextRegen <= 0 {
		p.ticksToNextRegen = PlayerRegenRateTicks
		p.Hp++
	}

	return false
}

func (p *PlayerCombat) Respawn(spawn *pb.Location) {
	p.Hp = PlayerMaxHp
	p.Alive = true
}
