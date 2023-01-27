package actions

import (
	"errors"
	"strings"
)

type ActionType string

const (
	ActionType_JoinGame   ActionType = "J"
	ActionType_Disconnect ActionType = "D"

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
	case ActionType_Move:
		return NewMoveActionFromMessage(playerId, split...)
	case ActionType_Attack:
		return NewAttackActionFromMessage(playerId, split...)
	case ActionType_JoinGame:
		return NewJoinGameActionFromMessage(playerId, split...)
	case ActionType_Disconnect:
		return NewDisconnectAction(playerId)
	default:
		return nil, errors.New("invalid message")
	}
}
