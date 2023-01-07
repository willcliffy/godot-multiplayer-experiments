package objects

type Position struct {
	X int
	Y int
	//z int
}

func New2DPosition(x, y int) *Position {
	return &Position{x, y}
}
