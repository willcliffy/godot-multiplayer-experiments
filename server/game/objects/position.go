package objects

type Location struct {
	X int
	Z int
	//z int
}

func New2DLocation(x, z int) *Location {
	return &Location{x, z}
}
