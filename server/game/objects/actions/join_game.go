package actions

import "fmt"

type JoinGameAction struct {
	PlayerId uint64
	Class    string
	//gameId string
}

func NewJoinGameAction(playerId uint64, class string) *JoinGameAction {
	return &JoinGameAction{
		PlayerId: playerId,
		Class:    class,
	}
}

func NewJoinGameActionFromMessage(playerId uint64, msg ...string) (*JoinGameAction, error) {
	if len(msg) != 4 {
		return nil, fmt.Errorf("invalid JoinGameAction, expected 4 segments but got %d", len(msg))
	} else if ActionType(msg[0]) != ActionType_JoinGame {
		return nil, fmt.Errorf("incorrect ActionType, expected %s but got %s", ActionType_JoinGame, msg[0])
	}

	if msg[3] == "" {
		return nil, fmt.Errorf("missing class: %s", msg[3])
	}

	// validating that source and target are in game is not the responsibility of this function
	return NewJoinGameAction(playerId, msg[3]), nil
}

func (self JoinGameAction) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", self.Id())), nil
}

func (self JoinGameAction) Id() string {
	return fmt.Sprintf("%v:%d::%s", self.Type(), self.PlayerId, self.Class)
}

func (self JoinGameAction) Type() ActionType {
	return ActionType_JoinGame
}
