package player

import "math/rand"

const (
	TeamColor_White  string = "#eae1f0"
	TeamColor_Grey   string = "#37313b"
	TeamColor_Black  string = "#1d1c1f"
	TeamColor_Orange string = "#89423f"
	TeamColor_Yellow string = "#fdbb27"
	TeamColor_Green  string = "#8d902e"
	TeamColor_Blue   string = "#4159cb"
	TeamColor_Teal   string = "#59a7af"

	NumberOfTeamColors = 8
)

var allTeamColors = []string{
	TeamColor_White,
	TeamColor_Grey,
	TeamColor_Black,
	TeamColor_Orange,
	TeamColor_Yellow,
	TeamColor_Green,
	TeamColor_Blue,
	TeamColor_Teal,
}

func RandomTeamColor() string {
	return allTeamColors[rand.Intn(NumberOfTeamColors)]
}
