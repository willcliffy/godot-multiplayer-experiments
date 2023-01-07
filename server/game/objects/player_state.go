package objects

type PlayerState uint8

const (
	PlayerState_Vibing       PlayerState = 1 // i.e. at beginning of game before spawn
	PlayerState_Alive        PlayerState = 2
	PlayerState_Dead         PlayerState = 3
	PlayerState_Invulnerable PlayerState = 4
)
