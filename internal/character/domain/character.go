package domain

import "strings"

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

// ArmorClass calculates the character's AC based on armor, dexterity, and shield
// This is D&D 5e business logic and belongs in the domain layer
func (c *Character) ArmorClass() int {
	baseAC := 10
	dexMod := Modifier(c.Dex)

	// Base AC from armor (D&D 5e rules)
	switch c.Armor {
	case "leather armor", "Leather Armor":
		baseAC = 11 + dexMod
	case "studded leather", "Studded Leather":
		baseAC = 12 + dexMod
	case "chain shirt", "Chain Shirt":
		// Medium armor: max +2 dex bonus
		if dexMod > 2 {
			baseAC = 13 + 2
		} else {
			baseAC = 13 + dexMod
		}
	case "scale mail", "Scale Mail":
		// Medium armor: max +2 dex bonus
		if dexMod > 2 {
			baseAC = 14 + 2
		} else {
			baseAC = 14 + dexMod
		}
	case "chain mail", "Chain Mail":
		// Heavy armor: no dex bonus
		baseAC = 16
	case "plate", "Plate":
		// Heavy armor: no dex bonus
		baseAC = 18
	default:
		// No armor: 10 + dex modifier
		baseAC = 10 + dexMod
	}

	// Shield bonus
	if c.Shield != "" {
		baseAC += 2
	}

	return baseAC
}

// PassivePerception calculates passive perception (10 + Wis mod + proficiency if proficient)
func (c *Character) PassivePerception() int {
	wisModifier := Modifier(c.Wis)
	passive := 10 + wisModifier
	
	// Check if proficient in perception
	for _, skill := range c.SkillProficiencies {
		if skill == "perception" || skill == "Perception" {
			passive += c.ProficiencyBonus
			break
		}
	}
	
	return passive
}

// SpellcastingAbility returns the primary spellcasting ability for the character's class
func (c *Character) SpellcastingAbility() string {
	switch c.Class {
	case "wizard", "Wizard":
		return "INT"
	case "cleric", "Cleric", "druid", "Druid", "ranger", "Ranger":
		return "WIS"
	case "sorcerer", "Sorcerer", "bard", "Bard", "paladin", "Paladin", "warlock", "Warlock":
		return "CHA"
	default:
		return ""
	}
}

// SpellcastingModifier returns the ability modifier used for spellcasting
func (c *Character) SpellcastingModifier() int {
	switch c.SpellcastingAbility() {
	case "INT":
		return Modifier(c.Int)
	case "WIS":
		return Modifier(c.Wis)
	case "CHA":
		return Modifier(c.Cha)
	default:
		return 0
	}
}

// SpellSaveDC calculates spell save DC (8 + proficiency + spellcasting modifier)
func (c *Character) SpellSaveDC() int {
	if c.SpellcastingAbility() == "" {
		return 0
	}
	return 8 + c.ProficiencyBonus + c.SpellcastingModifier()
}

// SpellAttackBonus calculates spell attack bonus (proficiency + spellcasting modifier)
func (c *Character) SpellAttackBonus() int {
	if c.SpellcastingAbility() == "" {
		return 0
	}
	return c.ProficiencyBonus + c.SpellcastingModifier()
}

// IsSpellcaster checks if the character's class can cast spells
func (c *Character) IsSpellcaster() bool {
	spellcasters := map[string]bool{
		"wizard":   true,
		"sorcerer": true,
		"warlock":  true,
		"bard":     true,
		"cleric":   true,
		"druid":    true,
		"paladin":  true,
		"ranger":   true,
	}
	return spellcasters[c.Class] || spellcasters[strings.ToLower(c.Class)]
}
