package broadcast

import pb "github.com/willcliffy/kilnwood-game-server/proto"

type MessageBroadcaster interface {
	OnPlayerJoinGame(gameId uint64, playerId uint64, response *pb.JoinGameResponse)
	OnPlayerLeftGame(gameId uint64, playerId uint64)
	OnGameTick(gameId uint64, tick *pb.GameTick)
}

type MessageReceiver interface {
	OnPlayerConnected(playerId uint64) error
	OnPlayerDisconnected(playerId uint64)
	OnActionReceived(playerId uint64, action *pb.Action)
}
