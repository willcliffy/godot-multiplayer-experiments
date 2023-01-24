package actions

import (
	"fmt"
	"strconv"

	"github.com/willcliffy/kilnwood-game-server/game/objects"
)

type MoveAction struct {
	PlayerId uint64
	Location objects.Location
}

func NewMoveAction(playerId uint64, x, z int) *MoveAction {
	return &MoveAction{
		PlayerId: playerId,
		Location: objects.Location{X: x, Z: z},
	}
}

func NewMoveActionFromMessage(playerId uint64, msg ...string) (*MoveAction, error) {
	if len(msg) != 4 {
		return nil, fmt.Errorf("invalid MoveAction, expected 4 segments but got %d", len(msg))
	} else if ActionType(msg[0]) != ActionType_Move {
		return nil, fmt.Errorf("incorrect ActionType, expected %s but got %s", ActionType_Move, msg[0])
	}

	if msg[1] == "" {
		return nil, fmt.Errorf("no source player provided for ")
	} else if msg[2] == "" || msg[3] == "" {
		return nil, fmt.Errorf("x or y coordinate not provided: %s, %s", msg[2], msg[3])
	}

	id, err := strconv.ParseUint(msg[1], 10, 64)
	if err != nil {
		return nil, err
	}

	locX, err := strconv.ParseInt(msg[2], 10, 64)
	if err != nil {
		return nil, err
	}

	locZ, err := strconv.ParseInt(msg[3], 10, 64)
	if err != nil {
		return nil, err
	}

	return NewMoveAction(id, int(locX), int(locZ)), nil
}

func (a MoveAction) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", a.Id())), nil
}

func (a MoveAction) Id() string {
	return fmt.Sprintf("%v:%d:%d:%d", a.Type(), a.PlayerId, a.Location.X, a.Location.Z)
}

func (a MoveAction) Type() ActionType {
	return ActionType_Move
}

func (a MoveAction) ToLocation() objects.Location {
	return objects.New2DLocation(a.Location.X, a.Location.Z)
}
