package util

type Broadcaster interface {
	Broadcast(gameId string, payload []byte) error
}
