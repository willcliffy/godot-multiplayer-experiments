package player

type CharacterType string

const (
	CharacterType_Fighter = "f"
	CharacterType_Ranged  = "r"
	CharacterType_Tank    = "t"
)

type CharacterClass interface {
	Type() CharacterType
	Range() int
}

func CharacterClassFromType(t string) (CharacterClass, bool) {
	switch CharacterType(t) {
	case CharacterType_Fighter:
		return nil, true
	case CharacterType_Ranged:
		return nil, true
	case CharacterType_Tank:
		return nil, true
	default:
		return nil, false
	}
}
