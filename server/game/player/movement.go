package player

import (
	"math"

	pb "github.com/willcliffy/kilnwood-game-server/proto"
)

type IPlayerMovement interface {
	GetLocation() *pb.Location
	SetPath([]*pb.Location)
}

type PlayerMovement struct {
	location            *pb.Location
	path                []*pb.Location
	ticksToNextLocation int32
}

func NewPlayerMovement() *PlayerMovement {
	return &PlayerMovement{}
}

func (p *PlayerMovement) Tick() *pb.Location {
	if len(p.path) == 0 {
		return nil
	}

	if p.ticksToNextLocation > 0 {
		p.ticksToNextLocation--
		return nil
	}

	p.moveToNextLocationOnPath()

	return p.location
}

func (p PlayerMovement) GetLocation() *pb.Location {
	return p.location
}

func (p *PlayerMovement) SetPath(path []*pb.Location) {
	if len(path) == 0 {
		return
	}

	p.path = path
	p.ticksToNextLocation = 0
}

func (p *PlayerMovement) moveToNextLocationOnPath() {
	p.location = p.path[0]
	if len(p.path) == 0 {
		return
	}

	p.path = p.path[1:]
	if len(p.path) > 0 {
		p.ticksToNextLocation = p.calculateTicksToNextLocation()
	}
}

func (p PlayerMovement) calculateTicksToNextLocation() int32 {
	dist := DistanceBetweenLocations(p.location, p.path[0])
	return int32(math.Round(dist / PLAYER_MOVE_SPEED))
}

func DistanceBetweenLocations(loc1, loc2 *pb.Location) float64 {
	return math.Sqrt(math.Pow(float64(loc2.X-loc1.X), 2) + math.Pow(float64(loc2.Z-loc1.Z), 2))
}
