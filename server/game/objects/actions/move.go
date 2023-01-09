package actions

import (
	"fmt"
	"strconv"

	"github.com/willcliffy/kilnwood-game-server/game/objects"
)

type MoveAction struct {
	PlayerId uint64
	Position objects.Position
}

func NewMoveAction(playerId uint64, x, z int) *MoveAction {
	return &MoveAction{
		PlayerId: playerId,
		Position: objects.Position{X: x, Z: z},
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

func (self MoveAction) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", self.Id())), nil
}

func (self MoveAction) Id() string {
	return fmt.Sprintf("%v:%d:%d:%d", self.Type(), self.PlayerId, self.Position.X, self.Position.Z)
}

func (self MoveAction) Type() ActionType {
	return ActionType_Move
}

func (self MoveAction) ToPosition() objects.Position {
	return *objects.New2DPosition(self.Position.X, self.Position.Z)
}
