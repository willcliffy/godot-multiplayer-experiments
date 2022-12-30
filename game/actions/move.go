package actions

import (
	"fmt"
	"strconv"
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
	} else if msg[0] != ActionType_Move {
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

	return NewMoveAction(msg[1], int(locX), int(locY)), nil
}

func (m MoveAction) ID() string {
	return fmt.Sprintf("%v:%s:%d:%d", m.Type(), m.source, m.locationX, m.locationY)
}

func (m MoveAction) Type() ActionType {
	return ActionType_Move
}

func (m MoveAction) SourcePlayer() string {
	return m.source
}

func (m MoveAction) TargetPlayer() *string {
	return nil
}
