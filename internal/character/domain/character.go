package domain

// Character represents a D&D 5e character with all their attributes and abilities.
type Character struct {
	Name               string      `json:"name"`
	Race               string      `json:"race"`
	Class              string      `json:"class"`
	Level              int         `json:"level"`
	Str                int         `json:"str"`
	Dex                int         `json:"dex"`
	Con                int         `json:"con"`
	Int                int         `json:"int"`
	Wis                int         `json:"wis"`
	Cha                int         `json:"cha"`
	Background         string      `json:"background"`
	ProficiencyBonus   int         `json:"proficiencyBonus"`
	SkillProficiencies []string    `json:"skillProficiencies"`
	SpellSlots         map[int]int `json:"spell_slots"` // key: spell level, value: slots available
	Weapon             string      `json:"weapon"`
	WeaponSlot         string      `json:"weapon_slot"`
	Armor              string      `json:"armor,omitempty"`
	Shield             string      `json:"shield,omitempty"`
	KnownSpells        []string    `json:"knownSpells,omitempty"`
	PreparedSpells     []string    `json:"preparedSpells,omitempty"`
}

// NewCharacter creates a new Character instance with proper spell slot calculation.
func NewCharacter(name, race, class string, level, str, dex, con, int_, wis, cha int, background string, skills []string) *Character {
	var spellSlots map[int]int
	switch class {
	case "Wizard", "Cleric", "Druid", "Bard", "Sorcerer":
		spellSlots = FullCasterSpellSlots(level)
	case "Paladin", "Ranger":
		spellSlots = HalfCasterSpellSlots(level)
	default:
		spellSlots = map[int]int{}
	}
	return &Character{
		Name:               name,
		Race:               race,
		Class:              class,
		Level:              level,
		Str:                str,
		Dex:                dex,
		Con:                con,
		Int:                int_,
		Wis:                wis,
		Cha:                cha,
		Background:         background,
		ProficiencyBonus:   ProficiencyBonus(level),
		SkillProficiencies: skills,
		SpellSlots:         spellSlots,
	}
}

// ProficiencyBonus calculates proficiency bonus based on level (D&D 5e rules)
func ProficiencyBonus(level int) int {
	return 2 + (level-1)/4
}
