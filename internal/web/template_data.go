package web

import (
	"strings"

	"DnD-sheet/internal/character/domain"
)

// WeaponAttack represents weapon attack information for display
type WeaponAttack struct {
	Name        string
	AttackBonus int
	Damage      string
	DamageType  string
	Range       string
	Properties  []string
	IsMelee     bool
	IsRanged    bool
	IsFinesse   bool
	IsTwoHanded bool
}

// CharacterTemplateData holds all data needed for the HTML character sheet template
type CharacterTemplateData struct {
	// Basic Character Info
	Name       string
	Race       string
	Class      string
	Level      int
	Background string

	// Ability Scores
	Str int
	Dex int
	Con int
	Int int
	Wis int
	Cha int

	// Ability Modifiers
	StrMod int
	DexMod int
	ConMod int
	IntMod int
	WisMod int
	ChaMod int

	// Proficiency and Combat Stats
	ProficiencyBonus  int
	ArmorClass        int
	Initiative        int
	Speed             int
	HitPointMax       int
	CurrentHP         int
	PassivePerception int

	// Equipment
	Weapon     string
	WeaponSlot string
	Armor      string
	Shield     string

	// Attack Information (for Attacks & Spellcasting section)
	WeaponAttacks []WeaponAttack

	// Skills (as a formatted string for display)
	SkillProficiencies []string
	SkillsDisplay      string

	// Spellcasting (if applicable)
	CanCastSpells        bool
	SpellcastingAbility  string
	SpellcastingModifier int
	SpellSaveDC          int
	SpellAttackBonus     int
	SpellSlots           map[int]int
	KnownSpells          []string
	PreparedSpells       []string

	// Saving Throws
	StrSave int
	DexSave int
	ConSave int
	IntSave int
	WisSave int
	ChaSave int

	// Individual Skill Modifiers for detailed display
	Acrobatics     int
	AnimalHandling int
	Arcana         int
	Athletics      int
	Deception      int
	History        int
	Insight        int
	Intimidation   int
	Investigation  int
	Medicine       int
	Nature         int
	Perception     int
	Performance    int
	Persuasion     int
	Religion       int
	SleightOfHand  int
	Stealth        int
	Survival       int
}

// NewCharacterTemplateData creates template data from a character domain object
func NewCharacterTemplateData(char *domain.Character) *CharacterTemplateData {
	data := &CharacterTemplateData{
		Name:       char.Name,
		Race:       char.Race,
		Class:      char.Class,
		Level:      char.Level,
		Background: char.Background,
		Str:        char.Str,
		Dex:        char.Dex,
		Con:        char.Con,
		Int:        char.Int,
		Wis:        char.Wis,
		Cha:        char.Cha,

		// Calculate ability modifiers using domain logic
		StrMod: domain.Modifier(char.Str),
		DexMod: domain.Modifier(char.Dex),
		ConMod: domain.Modifier(char.Con),
		IntMod: domain.Modifier(char.Int),
		WisMod: domain.Modifier(char.Wis),
		ChaMod: domain.Modifier(char.Cha),

		ProficiencyBonus:  char.ProficiencyBonus,
		ArmorClass:        char.ArmorClass(),
		Initiative:        char.Initiative(),
		Speed:             30, // Default speed, could be race-dependent
		PassivePerception: char.PassivePerception(),

		// Equipment
		Weapon:     char.Weapon,
		WeaponSlot: char.WeaponSlot,
		Armor:      char.Armor,
		Shield:     char.Shield,

		// Calculate weapon attacks
		WeaponAttacks: calculateWeaponAttacks(char),

		SkillProficiencies: char.SkillProficiencies,
		SkillsDisplay:      strings.Join(char.SkillProficiencies, ", "),

		// Spellcasting
		SpellSlots:     char.SpellSlots,
		KnownSpells:    char.KnownSpells,
		PreparedSpells: char.PreparedSpells,

		// Calculate HP
		HitPointMax: char.MaxHitPoints(),
		CurrentHP:   char.MaxHitPoints(), // Assume full HP for now
	}

	// Calculate spellcasting stats if applicable
	if char.IsSpellcaster() {
		data.CanCastSpells = true
		// Get ability name from domain (INT/WIS/CHA) and convert to full name
		spellAbility := char.SpellcastingAbility()
		switch spellAbility {
		case "INT":
			data.SpellcastingAbility = "Intelligence"
		case "WIS":
			data.SpellcastingAbility = "Wisdom"
		case "CHA":
			data.SpellcastingAbility = "Charisma"
		}
		data.SpellcastingModifier = char.SpellcastingModifier()
		data.SpellSaveDC = char.SpellSaveDC()
		data.SpellAttackBonus = char.SpellAttackBonus()
	}

	// Calculate saving throws
	data.StrSave = data.StrMod
	data.DexSave = data.DexMod
	data.ConSave = data.ConMod
	data.IntSave = data.IntMod
	data.WisSave = data.WisMod
	data.ChaSave = data.ChaMod

	// Add proficiency bonus to saving throws (simplified - assumes all classes get prof in 2 saves)
	// This would need to be expanded for proper class-based saving throw proficiencies

	// Calculate individual skill modifiers
	data.calculateSkillModifiers(char)

	return data
}

// calculateSkillModifiers calculates all skill modifiers with proficiency bonuses
func (data *CharacterTemplateData) calculateSkillModifiers(char *domain.Character) {
	// Base skill modifiers (ability modifier only)
	data.Acrobatics = data.DexMod
	data.AnimalHandling = data.WisMod
	data.Arcana = data.IntMod
	data.Athletics = data.StrMod
	data.Deception = data.ChaMod
	data.History = data.IntMod
	data.Insight = data.WisMod
	data.Intimidation = data.ChaMod
	data.Investigation = data.IntMod
	data.Medicine = data.WisMod
	data.Nature = data.IntMod
	data.Perception = data.WisMod
	data.Performance = data.ChaMod
	data.Persuasion = data.ChaMod
	data.Religion = data.IntMod
	data.SleightOfHand = data.DexMod
	data.Stealth = data.DexMod
	data.Survival = data.WisMod

	// Add proficiency bonus for proficient skills
	for _, skill := range char.SkillProficiencies {
		skillLower := strings.ToLower(skill)
		switch skillLower {
		case "acrobatics":
			data.Acrobatics += char.ProficiencyBonus
		case "animal handling":
			data.AnimalHandling += char.ProficiencyBonus
		case "arcana":
			data.Arcana += char.ProficiencyBonus
		case "athletics":
			data.Athletics += char.ProficiencyBonus
		case "deception":
			data.Deception += char.ProficiencyBonus
		case "history":
			data.History += char.ProficiencyBonus
		case "insight":
			data.Insight += char.ProficiencyBonus
		case "intimidation":
			data.Intimidation += char.ProficiencyBonus
		case "investigation":
			data.Investigation += char.ProficiencyBonus
		case "medicine":
			data.Medicine += char.ProficiencyBonus
		case "nature":
			data.Nature += char.ProficiencyBonus
		case "perception":
			data.Perception += char.ProficiencyBonus
		case "performance":
			data.Performance += char.ProficiencyBonus
		case "persuasion":
			data.Persuasion += char.ProficiencyBonus
		case "religion":
			data.Religion += char.ProficiencyBonus
		case "sleight of hand":
			data.SleightOfHand += char.ProficiencyBonus
		case "stealth":
			data.Stealth += char.ProficiencyBonus
		case "survival":
			data.Survival += char.ProficiencyBonus
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// calculateWeaponAttacks calculates weapon attack information for display
func calculateWeaponAttacks(char *domain.Character) []WeaponAttack {
	var attacks []WeaponAttack

	// If character has no weapon, return empty attacks
	if char.Weapon == "" {
		return attacks
	}

	// Calculate base attack bonuses using domain logic
	strMod := domain.Modifier(char.Str)
	dexMod := domain.Modifier(char.Dex)

	// Try to enrich weapon with API data (simplified for now)
	attack := WeaponAttack{
		Name: char.Weapon,
	}

	// Determine attack type and calculate bonuses based on weapon name
	weaponLower := strings.ToLower(char.Weapon)

	// Simple weapon categorization (would be better with API enrichment)
	switch {
	case strings.Contains(weaponLower, "bow") ||
		strings.Contains(weaponLower, "crossbow") ||
		strings.Contains(weaponLower, "sling") ||
		strings.Contains(weaponLower, "dart") ||
		strings.Contains(weaponLower, "javelin"):
		// Ranged weapons use Dex
		attack.AttackBonus = dexMod + char.ProficiencyBonus
		attack.IsRanged = true
		attack.Range = "150/600" // Default ranged weapon range

	case strings.Contains(weaponLower, "dagger") ||
		strings.Contains(weaponLower, "rapier") ||
		strings.Contains(weaponLower, "scimitar") ||
		strings.Contains(weaponLower, "shortsword"):
		// Finesse weapons can use Dex or Str (choose better)
		if dexMod > strMod {
			attack.AttackBonus = dexMod + char.ProficiencyBonus
		} else {
			attack.AttackBonus = strMod + char.ProficiencyBonus
		}
		attack.IsFinesse = true
		attack.IsMelee = true
		attack.Range = "5 ft"

	case strings.Contains(weaponLower, "greataxe") ||
		strings.Contains(weaponLower, "greatsword") ||
		strings.Contains(weaponLower, "maul") ||
		strings.Contains(weaponLower, "pike"):
		// Two-handed melee weapons use Str
		attack.AttackBonus = strMod + char.ProficiencyBonus
		attack.IsTwoHanded = true
		attack.IsMelee = true
		attack.Range = "5 ft"

	default:
		// Default melee weapons use Str
		attack.AttackBonus = strMod + char.ProficiencyBonus
		attack.IsMelee = true
		attack.Range = "5 ft"
	}

	// Set damage based on weapon (simplified - would be better with API data)
	attack.Damage = getDamageForWeapon(weaponLower)
	attack.DamageType = getDamageTypeForWeapon(weaponLower)

	attacks = append(attacks, attack)

	// Add spell attacks if character can cast spells
	if char.IsSpellcaster() {
		spellAttack := WeaponAttack{
			Name:        "Spell Attack",
			AttackBonus: char.SpellAttackBonus(),
			Damage:      "Varies",
			DamageType:  "Varies",
			Range:       "Varies",
			IsRanged:    true,
		}
		attacks = append(attacks, spellAttack)
	}

	return attacks
}

// getDamageForWeapon returns damage dice for common weapons
func getDamageForWeapon(weaponLower string) string {
	switch {
	case strings.Contains(weaponLower, "dagger"):
		return "1d4"
	case strings.Contains(weaponLower, "shortsword") || strings.Contains(weaponLower, "scimitar"):
		return "1d6"
	case strings.Contains(weaponLower, "longsword") || strings.Contains(weaponLower, "rapier"):
		return "1d8"
	case strings.Contains(weaponLower, "greatsword"):
		return "2d6"
	case strings.Contains(weaponLower, "greataxe"):
		return "1d12"
	case strings.Contains(weaponLower, "bow"):
		return "1d8"
	case strings.Contains(weaponLower, "crossbow"):
		if strings.Contains(weaponLower, "heavy") {
			return "1d10"
		}
		return "1d8"
	case strings.Contains(weaponLower, "club") || strings.Contains(weaponLower, "quarterstaff"):
		return "1d4"
	case strings.Contains(weaponLower, "mace") || strings.Contains(weaponLower, "spear"):
		return "1d6"
	case strings.Contains(weaponLower, "warhammer") || strings.Contains(weaponLower, "battleaxe"):
		return "1d8"
	case strings.Contains(weaponLower, "maul"):
		return "2d6"
	default:
		return "1d6" // Default damage
	}
}

// getDamageTypeForWeapon returns damage type for common weapons
func getDamageTypeForWeapon(weaponLower string) string {
	switch {
	case strings.Contains(weaponLower, "sword") ||
		strings.Contains(weaponLower, "axe") ||
		strings.Contains(weaponLower, "scimitar"):
		return "Slashing"
	case strings.Contains(weaponLower, "dagger") ||
		strings.Contains(weaponLower, "spear") ||
		strings.Contains(weaponLower, "rapier") ||
		strings.Contains(weaponLower, "arrow") ||
		strings.Contains(weaponLower, "bow"):
		return "Piercing"
	case strings.Contains(weaponLower, "mace") ||
		strings.Contains(weaponLower, "club") ||
		strings.Contains(weaponLower, "hammer") ||
		strings.Contains(weaponLower, "maul"):
		return "Bludgeoning"
	default:
		return "Bludgeoning" // Default damage type
	}
}
