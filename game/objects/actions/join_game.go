package actions

import "fmt"

type JoinGameAction struct {
	source string
	//gameId string
	Class string
}

func NewJoinGameAction(name, class string) *JoinGameAction {
	return &JoinGameAction{
		source: name,
		Class:  class,
	}
}

func NewJoinGameActionFromMessage(msg ...string) (*JoinGameAction, error) {
	if len(msg) != 4 {
		return nil, fmt.Errorf("invalid JoinGameAction, expected 4 segments but got %d", len(msg))
	} else if ActionType(msg[0]) != ActionType_JoinGame {
		return nil, fmt.Errorf("incorrect ActionType, expected %s but got %s", ActionType_JoinGame, msg[0])
	}

	if msg[1] == "" || msg[3] == "" {
		return nil, fmt.Errorf("missing name or class: %s, %s", msg[1], msg[3])
	}

	// validating that source and target are in game is not the responsibility of this function
	return NewJoinGameAction(msg[1], msg[3]), nil
}

func (self JoinGameAction) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", self.Id())), nil
}

func (self JoinGameAction) Id() string {
	return fmt.Sprintf("%v:%s:nil:%s", self.Type(), self.source, self.Class)
}

func (self JoinGameAction) Type() ActionType {
	return ActionType_JoinGame
}

func (self JoinGameAction) SourcePlayer() string {
	return self.source
}

func (self JoinGameAction) TargetPlayer() *string {
	return nil
}
