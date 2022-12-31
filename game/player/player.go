package player

import (
	"github.com/willcliffy/kilnwood-game-server/game/objects"
)

type Player struct {
	id        string
	Name      string
	class     CharacterClass
	Team      objects.Team
	targetPos *objects.Position
}

func NewPlayer(name, classType string, team objects.Team) *Player {
	class, ok := CharacterClassFromType(classType)
	if !ok {
		return nil
	}

	return &Player{name, name, class, team, nil}
}

func (p Player) ID() string {
	return p.id
}

func (p Player) DEBUG_Tile() [3][3]string {
	return [3][3]string{
		{"  /", "===", "\\  "},
		{" | ", p.Name[:3], " | "},
		{"  \\", "===", "/  "},
	}
}
