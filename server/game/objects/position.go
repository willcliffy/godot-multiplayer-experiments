package objects

type Position struct {
	X int
	Z int
	//z int
}

func New2DPosition(x, z int) *Position {
	return &Position{x, z}
}
