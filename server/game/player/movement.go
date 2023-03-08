package player

import (
	"math"

	"github.com/willcliffy/kilnwood-game-server/constants"
	pb "github.com/willcliffy/kilnwood-game-server/proto"
	"google.golang.org/protobuf/proto"
)

type IPlayerMovement interface {
	GetLocation() *pb.Location
	JumpToLocation(*pb.Location)
	SetPath([]*pb.Location)
	QueueAction(*pb.ClientAction)
	Tick(uint64) []*pb.ClientAction
}

type PlayerMovement struct {
	super               *Player
	location            *pb.Location
	path                []*pb.Location
	ticksToNextLocation int32
	queuedAction        *pb.ClientAction
}

func NewPlayerMovement(player *Player) *PlayerMovement {
	return &PlayerMovement{
		super: player,
	}
}

func (p *PlayerMovement) Tick(playerId uint64) []*pb.ClientAction {
	if len(p.path) == 0 {
		return nil
	}

	if p.ticksToNextLocation > 0 {
		p.ticksToNextLocation--
		return nil
	}

	var result []*pb.ClientAction

	p.location = p.path[0]
	p.path = p.path[1:]

	actionBytes, _ := proto.Marshal(&pb.Move{
		Path: []*pb.Location{p.location},
	})

	result = append(result, &pb.ClientAction{
		Type:     pb.ClientActionType_ACTION_MOVE,
		PlayerId: playerId,
		Payload:  actionBytes,
	})

	if len(p.path) != 0 {
		p.ticksToNextLocation = p.calculateTicksToNextLocation()
		return result
	}

	if p.queuedAction != nil {
		result = append(result, p.queuedAction)
	}

	return result
}

func (p PlayerMovement) GetLocation() *pb.Location {
	return p.location
}

func (p *PlayerMovement) JumpToLocation(location *pb.Location) {
	p.location = location
}

func (p *PlayerMovement) SetPath(path []*pb.Location) {
	if len(path) == 0 {
		return
	}

	p.path = path
	p.ticksToNextLocation = 0
}

func (p *PlayerMovement) QueueAction(action *pb.ClientAction) {
	p.queuedAction = action
}

func (p PlayerMovement) calculateTicksToNextLocation() int32 {
	dist := DistanceBetweenLocations(p.location, p.path[0])
	return int32(math.Round(dist/PLAYER_MOVE_SPEED) / constants.GAME_TICK_DURATION.Seconds())
}

func DistanceBetweenLocations(loc1, loc2 *pb.Location) float64 {
	return math.Sqrt(math.Pow(float64(loc2.X-loc1.X), 2) + math.Pow(float64(loc2.Z-loc1.Z), 2))
}
