package actions

import (
	"fmt"
	"strconv"

	"github.com/rs/zerolog/log"
	"github.com/willcliffy/kilnwood-game-server/game/objects"
)

type MoveAction struct {
	source    string
	locationX int
	locationY int
}

func NewMoveAction(sourcePlayer string, x, y int) *MoveAction {
	return &MoveAction{
		source:    sourcePlayer,
		locationX: x,
		locationY: y,
	}
}

func NewMoveActionFromMessage(msg ...string) (*MoveAction, error) {
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

	locX, err := strconv.ParseInt(msg[2], 10, 64)
	if err != nil {
		return nil, err
	}

	locY, err := strconv.ParseInt(msg[3], 10, 64)
	if err != nil {
		return nil, err
	}

	log.Debug().Msgf("returning move action")
	return NewMoveAction(msg[1], int(locX), int(locY)), nil
}

func (self MoveAction) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", self.Id())), nil
}

func (self MoveAction) Id() string {
	return fmt.Sprintf("%v:%s:%d:%d", self.Type(), self.source, self.locationX, self.locationY)
}

func (self MoveAction) Type() ActionType {
	return ActionType_Move
}

func (self MoveAction) SourcePlayer() string {
	return self.source
}

func (self MoveAction) TargetPlayer() *string {
	return nil
}

func (self MoveAction) ToPosition() objects.Position {
	return *objects.New2DPosition(self.locationX, self.locationY)
}
