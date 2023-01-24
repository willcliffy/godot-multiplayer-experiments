package player

import (
	"github.com/rs/zerolog/log"
	"github.com/willcliffy/kilnwood-game-server/game/objects"
)

type Player struct {
	id    uint64
	Name  string
	class *CharacterClass
	//Team     objects.Team
	Color    objects.TeamColor
	Location objects.Location
	Combat   PlayerCombatStats
}

func NewPlayer(playerId uint64, name, classType string, color objects.TeamColor) *Player {
	class, ok := CharacterClassFromType(classType)
	if !ok {
		log.Warn().Msgf("Could not parse character class: %v", classType)
		return nil
	}

	return &Player{
		id:    playerId,
		Name:  name,
		class: class,
		//Team:     team,
		Color:    color,
		Location: objects.Location{},
		Combat:   NewPlayerCombatStats(),
	}
}

func (p Player) Id() uint64 {
	return p.id
}

func (p *Player) SetTargetLocation(location objects.Location) {
	p.Location = location
}

func (p Player) GetTargetLocation() objects.Location {
	return p.Location
}

func (p *Player) Tick() {
	p.Combat.Tick()
}

func (p *Player) SetPlayerState(state objects.PlayerState) {}

func (p Player) DEBUG_Tile() [3][3]string {
	return [3][3]string{
		{" / ", " = ", " \\ "},
		{"|  ", p.Name, "  |"},
		{" \\ ", " = ", " / "},
	}
}
