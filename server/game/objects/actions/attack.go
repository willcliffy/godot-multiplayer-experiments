package actions

import (
	"fmt"
	"strconv"

	"github.com/rs/zerolog/log"
)

type AttackAction struct {
	sourcePlayerId uint64
	targetPlayerId uint64
}

func NewAttackAction(sourcePlayerId, targetPlayerId uint64) *AttackAction {
	return &AttackAction{
		sourcePlayerId: sourcePlayerId,
		targetPlayerId: targetPlayerId,
	}
}

func NewAttackActionFromMessage(playerId uint64, msg ...string) (*AttackAction, error) {
	if len(msg) != 3 {
		return nil, fmt.Errorf("invalid AttackAction, expected 3 segments but got %d", len(msg))
	} else if ActionType(msg[0]) != ActionType_Attack {
		return nil, fmt.Errorf("incorrect ActionType, expected %s but got %s", ActionType_Attack, msg[0])
	}

	if msg[1] == "" || msg[2] == "" {
		return nil, fmt.Errorf("missing source or target for attack: %s, %s", msg[1], msg[2])
	}

	sourceId, err := strconv.ParseUint(msg[1], 10, 64)
	if err != nil {
		return nil, err
	}

	if sourceId != playerId {
		log.Warn().Msgf("address does not match player Id. address tied to '%d' but message claims to be from '%d'", playerId, sourceId)
		return nil, fmt.Errorf("hax")
	}

	targetId, err := strconv.ParseUint(msg[2], 10, 64)
	if err != nil {
		return nil, err
	}

	// validating that source and target are in game is not the responsibility of this function
	return NewAttackAction(sourceId, targetId), nil
}

func (self AttackAction) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", self.Id())), nil
}

func (self AttackAction) Id() string {
	return fmt.Sprintf("%v:%d:%d", self.Type(), self.sourcePlayerId, self.targetPlayerId)
}

func (self AttackAction) Type() ActionType {
	return ActionType_Attack
}
