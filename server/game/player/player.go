package player

import (
	"github.com/rs/zerolog/log"
	"github.com/willcliffy/kilnwood-game-server/game/objects"
)

type Player struct {
	id        uint64
	Name      string
	class     *CharacterClass
	Team      objects.Team
	targetPos *objects.Position
}

func NewPlayer(playerId uint64, name, classType string, team objects.Team) *Player {
	class, ok := CharacterClassFromType(classType)
	if !ok {
		log.Warn().Msgf("Could not parse character class: %v", classType)
		return nil
	}

	return &Player{playerId, name, class, team, nil}
}

func (p Player) Id() uint64 {
	return p.id
}

func (p *Player) SetTargetLocation(location objects.Position) {
	p.targetPos = &location
}

func (p *Player) SetPlayerState(state objects.PlayerState) {

}

func (p Player) DEBUG_Tile() [3][3]string {
	return [3][3]string{
		{" / ", " = ", " \\ "},
		{"|  ", p.Name, "  |"},
		{" \\ ", " = ", " / "},
	}
}