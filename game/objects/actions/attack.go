package actions

import "fmt"

type AttackAction struct {
	source string
	target string
}

func NewAttackAction(sourcePlayer, targetPlayer string) *AttackAction {
	return &AttackAction{
		source: sourcePlayer,
		target: targetPlayer,
	}
}

func NewAttackActionFromMessage(msg ...string) (*AttackAction, error) {
	if len(msg) != 3 {
		return nil, fmt.Errorf("invalid AttackAction, expected 3 segments but got %d", len(msg))
	} else if ActionType(msg[0]) != ActionType_Attack {
		return nil, fmt.Errorf("incorrect ActionType, expected %s but got %s", ActionType_Attack, msg[0])
	}

	if msg[1] == "" || msg[2] == "" {
		return nil, fmt.Errorf("missing source or target for attack: %s, %s", msg[1], msg[2])
	}

	// validating that source and target are in game is not the responsibility of this function
	return NewAttackAction(msg[1], msg[2]), nil
}

func (self AttackAction) ID() string {
	return fmt.Sprintf("%v:%s:%s", self.Type(), self.source, self.target)
}

func (self AttackAction) Type() ActionType {
	return ActionType_Attack
}

func (self AttackAction) SourcePlayer() string {
	return self.source
}

func (self AttackAction) TargetPlayer() *string {
	return &self.target
}
