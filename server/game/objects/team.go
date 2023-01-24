package objects

import "math/rand"

type Team uint8

const (
	Team_Red  Team = 1
	Team_Blue Team = 2
)

type TeamColor string

const (
	TeamColor_White  TeamColor = "#eae1f0"
	TeamColor_Grey   TeamColor = "#37313b"
	TeamColor_Black  TeamColor = "#1d1c1f"
	TeamColor_Orange TeamColor = "#89423f"
	TeamColor_Yellow TeamColor = "#fdbb27"
	TeamColor_Green  TeamColor = "#8d902e"
	TeamColor_Blue   TeamColor = "#4159cb"
	TeamColor_Teal   TeamColor = "#59a7af"

	NumberOfTeamColors = 8
)

var allTeamColors = []TeamColor{
	TeamColor_White,
	TeamColor_Grey,
	TeamColor_Black,
	TeamColor_Orange,
	TeamColor_Yellow,
	TeamColor_Green,
	TeamColor_Blue,
	TeamColor_Teal,
}

func RandomTeamColor() TeamColor {
	return allTeamColors[rand.Intn(NumberOfTeamColors)]
}
