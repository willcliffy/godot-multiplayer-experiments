package player

type PlayerCombatStats struct {
	HP    int // in HP units
	AD    int // in HP units
	Regen int // in ticks

	currentHP        int
	ticksToNextRegen int
}

func NewPlayerCombatStats() PlayerCombatStats {
	return PlayerCombatStats{
		HP:               10,
		AD:               1,
		Regen:            10,
		ticksToNextRegen: 10,
	}
}

func (p *PlayerCombatStats) Tick() {
	if p.currentHP == 0 || p.currentHP >= p.HP {
		return
	}

	p.ticksToNextRegen--
	if p.ticksToNextRegen <= 0 {
		p.ticksToNextRegen = p.Regen
		p.HP++
	}
}

func (p *PlayerCombatStats) ApplyDamage(dmg int) int {
	p.HP -= dmg
	return dmg
}
