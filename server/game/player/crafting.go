package player

import pb "github.com/willcliffy/kilnwood-game-server/proto"

type IPlayerCrafting interface {
	AddResource(pb.ResourceType, int)
	Craft(*pb.Location) bool
	Tick()
}

type PlayerCrafting struct {
	resources map[pb.ResourceType]int
	copies    []*pb.Location
}

func NewPlayerCrafting(player *Player) *PlayerCrafting {
	return &PlayerCrafting{
		resources: make(map[pb.ResourceType]int),
	}
}

func (p *PlayerCrafting) Tick() {

}

func (p *PlayerCrafting) AddResource(resourceType pb.ResourceType, amount int) {
	p.resources[resourceType] += amount
}

func (p *PlayerCrafting) Craft(location *pb.Location) bool {
	red, ok := p.resources[pb.ResourceType_RED]
	if !ok {
		return false
	}
	blue, ok := p.resources[pb.ResourceType_BLUE]
	if !ok {
		return false
	}

	if red < 1 || blue < 2 {
		return false
	}

	p.resources[pb.ResourceType_RED] -= 1
	p.resources[pb.ResourceType_BLUE] -= 2
	p.copies = append(p.copies, location)

	return true
}
