package actions

import (
	"errors"
	"fmt"
	"strings"
)

type ActionType string

const (
	ActionType_CancelAction ActionType = "C"
	ActionType_JoinGame     ActionType = "J"

	ActionType_Move   ActionType = "m"
	ActionType_Attack ActionType = "a"
)

type Action interface {
	MarshalJSON() ([]byte, error)
	Id() string
	Type() ActionType
}

func ParseActionFromMessage(playerId uint64, msg string) (Action, error) {
	split := strings.Split(strings.TrimSpace(msg), ":")
	if len(split) == 1 {
		return nil, errors.New("empty message")
	}

	switch ActionType(split[0]) {
	case ActionType_CancelAction:
		return nil, fmt.Errorf("cancelling actions nyi")
	case ActionType_Move:
		return NewMoveActionFromMessage(playerId, split...)
	case ActionType_Attack:
		return NewAttackActionFromMessage(playerId, split...)
	case ActionType_JoinGame:
		return NewJoinGameActionFromMessage(playerId, split...)
	default:
		return nil, errors.New("invalid message")
	}
}
