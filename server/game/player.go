package game

import (
	"math/rand"

	pb "github.com/willcliffy/kilnwood-game-server/proto"
)

const (
	TeamColor_White  string = "#eae1f0"
	TeamColor_Grey   string = "#37313b"
	TeamColor_Black  string = "#1d1c1f"
	TeamColor_Orange string = "#89423f"
	TeamColor_Yellow string = "#fdbb27"
	TeamColor_Green  string = "#8d902e"
	TeamColor_Blue   string = "#4159cb"
	TeamColor_Teal   string = "#59a7af"

	NumberOfTeamColors = 8
)

var allTeamColors = []string{
	TeamColor_White,
	TeamColor_Grey,
	TeamColor_Black,
	TeamColor_Orange,
	TeamColor_Yellow,
	TeamColor_Green,
	TeamColor_Blue,
	TeamColor_Teal,
}

func RandomTeamColor() string {
	return allTeamColors[rand.Intn(NumberOfTeamColors)]
}

type Player struct {
	Id       uint64
	Color    string
	Location *pb.Location
	combat   *PlayerCombatStats
}

func NewPlayer(playerId uint64, color string) *Player {
	return &Player{
		Id:       playerId,
		Color:    color,
		Location: &pb.Location{},
		combat:   NewPlayerCombatStats(),
	}
}

func (p *Player) Tick() {
	p.combat.Tick()
}

func (p *Player) AttackDamage() int {
	return p.combat.AD
}

func (p *Player) ApplyDamage(dmg int) int {
	return p.combat.ApplyDamage(dmg)
}

type PlayerCombatStats struct {
	HP    int // in HP units
	AD    int // in HP units
	Regen int // in ticks

	currentHP        int
	ticksToNextRegen int
}

func NewPlayerCombatStats() *PlayerCombatStats {
	return &PlayerCombatStats{
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
