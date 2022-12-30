package actions

import (
	"errors"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
)

type ActionType string

const (
	ActionType_CancelAction = "C"
	ActionType_Move         = "m"
	ActionType_Attack       = "a"
)

type Action interface {
	ID() string
	Type() ActionType
	SourcePlayer() string
	TargetPlayer() *string
}

func ParseActionFromMessage(msg string) (Action, error) {
	split := strings.Split(strings.TrimSpace(msg), ":")
	if len(split) == 1 {
		return nil, errors.New("empty message")
	}

	log.Debug().Msgf("parsing message: %v", split)
	switch split[0] {
	case ActionType_CancelAction:
		return nil, fmt.Errorf("cancelling actions nyi")
	case ActionType_Move:
		return NewMoveActionFromMessage(split...)
	case ActionType_Attack:
		return NewAttackActionFromMessage(split...)
	}

	return nil, nil
}
