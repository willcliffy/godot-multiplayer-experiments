package gamemap

// TODO - placeholder
type Tile [3][3]string

func (self *Tile) Copy() Tile {
	return *self
}
