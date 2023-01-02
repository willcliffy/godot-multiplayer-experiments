package player

type CharacterType string

const (
	CharacterType_Fighter = "f"
	CharacterType_Ranged  = "r"
	CharacterType_Tank    = "t"
)

type CharacterClass struct {
	Type         CharacterType
	AttackRange  int
	AttackSpeed  int
	AttackDamage int
	HpTotal      int
	HpRemaining  int
}

func NewCharacterClass_Fighter() *CharacterClass {
	return &CharacterClass{
		Type:         CharacterType_Fighter,
		AttackRange:  1,
		AttackSpeed:  1,
		AttackDamage: 5,
		HpTotal:      15,
		HpRemaining:  15,
	}
}

func NewCharacterClass_Ranged() *CharacterClass {
	return &CharacterClass{
		Type:         CharacterType_Ranged,
		AttackRange:  1,
		AttackSpeed:  1,
		AttackDamage: 4,
		HpTotal:      10,
		HpRemaining:  10,
	}
}

func NewCharacterClass_Tank() *CharacterClass {
	return &CharacterClass{
		Type:         CharacterType_Tank,
		AttackRange:  2,
		AttackSpeed:  2,
		AttackDamage: 2,
		HpTotal:      20,
		HpRemaining:  20,
	}
}

func CharacterClassFromType(t string) (*CharacterClass, bool) {
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
