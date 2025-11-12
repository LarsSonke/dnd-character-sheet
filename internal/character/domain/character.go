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
		"wizard":           true,
		"sorcerer":         true,
		"warlock":          true,
		"bard":             true,
		"cleric":           true,
		"druid":            true,
		"paladin":          true,
		"ranger":           true,
		"artificer":        true,
		"eldritch knight":  true,
		"arcane trickster": true,
	}
	classLower := strings.ToLower(c.Class)
	return spellcasters[classLower]
}

// IsPreparedCaster returns true if the class prepares spells (vs learning them)
// D&D 5e rule: Some classes learn spells permanently, others prepare daily
func (c *Character) IsPreparedCaster() bool {
	// Known caster classes that learn spells permanently
	knownCasters := map[string]bool{
		"sorcerer":         true,
		"warlock":          true,
		"bard":             true,
		"eldritch knight":  true,
		"arcane trickster": true,
		"ranger":           true, // Rangers know spells in 5e
	}

	classLower := strings.ToLower(c.Class)
	// If it's a known caster, it's NOT a prepared caster
	if knownCasters[classLower] {
		return false
	}

	// If it's a spellcaster but not a known caster, it's a prepared caster
	return c.IsSpellcaster()
}

// GetSpellSlots returns spell slots for the character based on class and level
// D&D 5e rule: Different classes have different spell slot progressions
func (c *Character) GetSpellSlots() map[int]int {
	classLower := strings.ToLower(c.Class)

	switch classLower {
	case "wizard", "cleric", "druid", "bard", "sorcerer":
		slots := FullCasterSpellSlots(c.Level)
		// Add cantrips (Level 0) for full casters
		slots[0] = FullCasterCantrips(c.Level)
		return slots
	case "paladin", "ranger":
		return HalfCasterSpellSlots(c.Level)
	case "warlock":
		return PactMagicSpellSlots(c.Level)
	default:
		return map[int]int{}
	}
}

// Initiative calculates initiative bonus (Dex modifier + class bonuses)
// D&D 5e rule: Initiative = Dex modifier, with class-specific bonuses
func (c *Character) Initiative() int {
	initiative := Modifier(c.Dex)

	// Class-specific initiative bonuses (D&D 5e rules)
	classLower := strings.ToLower(c.Class)
	switch classLower {
	case "bard":
		// Jack of All Trades: add half proficiency to initiative (from level 2)
		if c.Level >= 2 {
			initiative += c.ProficiencyBonus / 2
		}
		// Future: could add other class features like Feral Instinct for Barbarian
	}

	return initiative
}

// MaxHitPoints calculates maximum hit points based on class and level
// D&D 5e rule: Class hit die + Con modifier per level
func (c *Character) MaxHitPoints() int {
	conMod := Modifier(c.Con)

	// Class hit dice (D&D 5e rules)
	classLower := strings.ToLower(c.Class)
	var baseHP int

	switch classLower {
	case "barbarian":
		// d12 hit die: max at 1st level, average (7) afterwards
		baseHP = 12 + (c.Level-1)*7
	case "fighter", "paladin", "ranger":
		// d10 hit die: max at 1st level, average (6) afterwards
		baseHP = 10 + (c.Level-1)*6
	case "bard", "cleric", "druid", "monk", "rogue", "warlock":
		// d8 hit die: max at 1st level, average (5) afterwards
		baseHP = 8 + (c.Level-1)*5
	case "artificer", "sorcerer", "wizard":
		// d6 hit die: max at 1st level, average (4) afterwards
		baseHP = 6 + (c.Level-1)*4
	default:
		// Default to d8 if class unknown
		baseHP = 8 + (c.Level-1)*5
	}

	// Add Constitution modifier for each level
	totalHP := baseHP + (conMod * c.Level)

	// Minimum 1 HP per level
	if totalHP < c.Level {
		totalHP = c.Level
	}

	return totalHP
}

// Race represents a D&D 5e character race and its mechanical effects
type Race struct {
	Name string
}

// NewRace creates a new Race instance
func NewRace(name string) *Race {
	return &Race{Name: name}
}

// GetAbilityBonuses returns the ability score bonuses for this race according to D&D 5e rules
func (r *Race) GetAbilityBonuses() map[string]int {
	bonuses := make(map[string]int)

	switch strings.ToLower(r.Name) {
	case "dwarf":
		bonuses["con"] = 2
	case "elf":
		bonuses["dex"] = 2
	case "halfling":
		bonuses["dex"] = 2
	case "lightfoot halfling":
		bonuses["dex"] = 2
		bonuses["cha"] = 1
	case "stout halfling":
		bonuses["dex"] = 2
		bonuses["con"] = 1
	case "human":
		bonuses["str"] = 1
		bonuses["dex"] = 1
		bonuses["con"] = 1
		bonuses["int"] = 1
		bonuses["wis"] = 1
		bonuses["cha"] = 1
	case "dragonborn":
		bonuses["str"] = 2
		bonuses["cha"] = 1
	case "gnome":
		bonuses["int"] = 2
	case "half elf":
		bonuses["cha"] = 2
	case "half orc":
		bonuses["str"] = 2
		bonuses["con"] = 1
	case "tiefling":
		bonuses["int"] = 1
		bonuses["cha"] = 2
	case "hill dwarf":
		bonuses["con"] = 2
		bonuses["wis"] = 1
	}

	return bonuses
}

// Background represents a D&D 5e character background
type Background struct {
	Name string
}

// NewBackground creates a new Background instance
func NewBackground(name string) *Background {
	return &Background{Name: name}
}

// GetSkillProficiencies returns the skill proficiencies for this background according to D&D 5e rules
func (b *Background) GetSkillProficiencies() []string {
	backgroundSkills := map[string][]string{
		"acolyte":       {"insight", "religion"},
		"charlatan":     {"deception", "sleight of hand"},
		"criminal":      {"deception", "stealth"},
		"entertainer":   {"acrobatics", "performance"},
		"folk hero":     {"animal handling", "survival"},
		"guild artisan": {"insight", "persuasion"},
		"hermit":        {"medicine", "religion"},
		"noble":         {"history", "persuasion"},
		"outlander":     {"athletics", "survival"},
		"sage":          {"arcana", "history"},
		"sailor":        {"athletics", "perception"},
		"soldier":       {"athletics", "intimidation"},
		"urchin":        {"sleight of hand", "stealth"},
	}

	return backgroundSkills[strings.ToLower(b.Name)]
}

// Class represents a D&D 5e character class
type Class struct {
	Name string
}

// NewClass creates a new Class instance
func NewClass(name string) *Class {
	return &Class{Name: name}
}

// GetAvailableSkills returns the skills that this class can choose from according to D&D 5e rules
func (cl *Class) GetAvailableSkills() []string {
	classSkills := map[string][]string{
		"barbarian": {"animal handling", "athletics", "intimidation", "nature", "perception", "survival"},
		"bard":      {"acrobatics", "animal handling", "arcana", "athletics", "deception", "history", "insight", "intimidation", "investigation", "medicine", "nature", "perception", "performance", "persuasion", "religion", "sleight of hand", "stealth", "survival"},
		"cleric":    {"history", "insight", "medicine", "persuasion", "religion"},
		"druid":     {"arcana", "animal handling", "insight", "medicine", "nature", "perception", "religion", "survival"},
		"fighter":   {"acrobatics", "animal handling", "athletics", "history", "insight", "intimidation", "perception", "survival"},
		"monk":      {"acrobatics", "athletics", "history", "insight", "religion", "stealth"},
		"paladin":   {"athletics", "insight", "intimidation", "medicine", "persuasion", "religion"},
		"ranger":    {"animal handling", "athletics", "insight", "investigation", "nature", "perception", "stealth", "survival"},
		"rogue":     {"acrobatics", "athletics", "deception", "insight", "intimidation", "investigation", "perception", "performance", "persuasion", "sleight of hand", "stealth"},
		"sorcerer":  {"arcana", "deception", "insight", "intimidation", "persuasion", "religion"},
		"warlock":   {"arcana", "deception", "history", "intimidation", "investigation", "nature", "religion"},
		"wizard":    {"arcana", "history", "insight", "investigation", "medicine", "religion"},
	}

	return classSkills[strings.ToLower(cl.Name)]
}

// GetSkillCount returns the number of skills this class can choose according to D&D 5e rules
func (cl *Class) GetSkillCount() int {
	classSkillCount := map[string]int{
		"barbarian": 2,
		"bard":      3,
		"cleric":    2,
		"druid":     2,
		"fighter":   2,
		"monk":      2,
		"paladin":   2,
		"ranger":    3,
		"rogue":     4,
		"sorcerer":  2,
		"warlock":   2,
		"wizard":    2,
	}

	return classSkillCount[strings.ToLower(cl.Name)]
}
