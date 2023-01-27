package actions

import "fmt"

type DisconnectAction struct {
	All      bool
	PlayerId uint64
}

func NewDisconnectAction(playerId uint64) (*DisconnectAction, error) {
	return &DisconnectAction{PlayerId: playerId}, nil
}

func NewDisconnectAllAction() *DisconnectAction {
	return &DisconnectAction{
		All: true,
	}
}

func (d DisconnectAction) Id() string {
	if d.All {
		return "d:all"
	}

	return fmt.Sprintf("d:%d", d.PlayerId)
}

func (d DisconnectAction) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", d.Id())), nil
}

func (DisconnectAction) Type() ActionType {
	return ActionType_Disconnect
}
