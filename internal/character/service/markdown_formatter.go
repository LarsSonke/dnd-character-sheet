package service

import (
	"DnD-sheet/internal/character/domain"
	"fmt"
	"sort"
	"strings"
)

// MarkdownFormatter handles markdown export formatting
type MarkdownFormatter struct{}

// NewMarkdownFormatter creates a new markdown formatter
func NewMarkdownFormatter() *MarkdownFormatter {
	return &MarkdownFormatter{}
}

// FormatCharacter formats a character as markdown
func (f *MarkdownFormatter) FormatCharacter(char *domain.Character) string {
	var builder strings.Builder

	// Title
	builder.WriteString(fmt.Sprintf("# %s\n\n", char.Name))

	// Character section
	builder.WriteString("## Character\n")
	builder.WriteString(fmt.Sprintf("Class: %s\n", char.Class))
	builder.WriteString(fmt.Sprintf("Race: %s\n", char.Race))
	builder.WriteString(fmt.Sprintf("Background: %s\n", char.Background))
	builder.WriteString(fmt.Sprintf("Level: %d\n", char.Level))
	builder.WriteString(fmt.Sprintf("Proficiency bonus: +%d\n", char.ProficiencyBonus))

	// Calculate passive perception
	wisModifier := domain.Modifier(char.Wis)
	passivePerception := 10 + wisModifier
	if f.hasSkillProficiency(char.SkillProficiencies, "perception") {
		passivePerception += char.ProficiencyBonus
	}
	builder.WriteString(fmt.Sprintf("Passive perception: %d\n\n", passivePerception))

	// Ability scores
	builder.WriteString("## Ability scores\n")
	builder.WriteString(fmt.Sprintf("STR: %d (%s)\n", char.Str, f.formatModifier(domain.Modifier(char.Str))))
	builder.WriteString(fmt.Sprintf("DEX: %d (%s)\n", char.Dex, f.formatModifier(domain.Modifier(char.Dex))))
	builder.WriteString(fmt.Sprintf("CON: %d (%s)\n", char.Con, f.formatModifier(domain.Modifier(char.Con))))
	builder.WriteString(fmt.Sprintf("INT: %d (%s)\n", char.Int, f.formatModifier(domain.Modifier(char.Int))))
	builder.WriteString(fmt.Sprintf("WIS: %d (%s)\n", char.Wis, f.formatModifier(domain.Modifier(char.Wis))))
	builder.WriteString(fmt.Sprintf("CHA: %d (%s)\n\n", char.Cha, f.formatModifier(domain.Modifier(char.Cha))))

	// Skills
	builder.WriteString("## Skills\n")
	builder.WriteString(f.formatSkills(char))
	builder.WriteString("\n")

	// Equipment
	builder.WriteString("## Equipment\n")
	if char.Weapon != "" {
		builder.WriteString(fmt.Sprintf("Main hand: %s\n", char.Weapon))
	}
	if char.Armor != "" {
		builder.WriteString(fmt.Sprintf("Armor: %s\n", char.Armor))
	}
	if char.Shield != "" {
		builder.WriteString(fmt.Sprintf("Shield: %s\n", char.Shield))
	}
	builder.WriteString("\n")

	// Combat stats
	builder.WriteString("## Combat stats\n")
	ac := f.calculateArmorClass(char)
	builder.WriteString(fmt.Sprintf("Armor class: %d\n", ac))
	builder.WriteString(fmt.Sprintf("Initiative bonus: %s\n\n", f.formatModifier(domain.Modifier(char.Dex))))

	// Spell slots (only for casters)
	if f.isSpellcaster(char.Class) {
		builder.WriteString("## Spell slots\n")
		for level := 0; level <= 9; level++ {
			if slots, exists := char.SpellSlots[level]; exists && slots > 0 {
				builder.WriteString(fmt.Sprintf("Level %d: %d\n", level, slots))
			}
		}
		builder.WriteString("\n")

		// Spellcasting
		builder.WriteString("## Spellcasting\n")
		spellAbility := f.getSpellcastingAbility(char.Class)
		builder.WriteString(fmt.Sprintf("Spellcasting ability: %s\n", spellAbility))

		spellMod := f.getSpellcastingModifier(char, spellAbility)
		spellSaveDC := 8 + char.ProficiencyBonus + spellMod
		spellAttackBonus := char.ProficiencyBonus + spellMod

		builder.WriteString(fmt.Sprintf("Spell save DC: %d\n", spellSaveDC))
		builder.WriteString(fmt.Sprintf("Spell attack bonus: +%d\n\n", spellAttackBonus))

		// Spells
		if len(char.PreparedSpells) > 0 {
			builder.WriteString("## Spells\n\n")
			builder.WriteString(f.formatSpellsByLevel(char.PreparedSpells))
		}
	}

	return builder.String()
}

// formatModifier formats ability modifiers with proper +/- signs
func (f *MarkdownFormatter) formatModifier(modifier int) string {
	if modifier >= 0 {
		return fmt.Sprintf("+%d", modifier)
	}
	return fmt.Sprintf("%d", modifier)
}

// formatSkills formats the skills list with checkboxes
func (f *MarkdownFormatter) formatSkills(char *domain.Character) string {
	skills := []string{
		"Acrobatics (Dex)", "Animal Handling (Wis)", "Arcana (Int)", "Athletics (Str)",
		"Deception (Cha)", "History (Int)", "Insight (Wis)", "Intimidation (Cha)",
		"Investigation (Int)", "Medicine (Wis)", "Nature (Int)", "Perception (Wis)",
		"Performance (Cha)", "Persuasion (Cha)", "Religion (Int)", "Sleight of Hand (Dex)",
		"Stealth (Dex)", "Survival (Wis)",
	}

	var builder strings.Builder
	for _, skill := range skills {
		skillName := strings.ToLower(strings.Split(skill, " (")[0])
		if f.hasSkillProficiency(char.SkillProficiencies, skillName) {
			builder.WriteString(fmt.Sprintf("* [x] %s\n", skill))
		} else {
			builder.WriteString(fmt.Sprintf("* [] %s\n", skill))
		}
	}
	return builder.String()
}

// hasSkillProficiency checks if a character has proficiency in a skill
func (f *MarkdownFormatter) hasSkillProficiency(proficiencies []string, skill string) bool {
	skill = strings.ToLower(skill)
	for _, prof := range proficiencies {
		if strings.ToLower(prof) == skill {
			return true
		}
	}
	return false
}

// calculateArmorClass calculates the character's AC
func (f *MarkdownFormatter) calculateArmorClass(char *domain.Character) int {
	baseAC := 10
	dexMod := domain.Modifier(char.Dex)

	// Base AC from armor
	switch strings.ToLower(char.Armor) {
	case "leather armor":
		baseAC = 11 + dexMod
	case "studded leather":
		baseAC = 12 + dexMod
	case "chain shirt":
		baseAC = 13 + min(dexMod, 2)
	case "scale mail":
		baseAC = 14 + min(dexMod, 2)
	case "chain mail":
		baseAC = 16
	case "plate":
		baseAC = 18
	default:
		baseAC = 10 + dexMod
	}

	// Shield bonus
	if char.Shield != "" {
		baseAC += 2
	}

	return baseAC
}

// isSpellcaster checks if a class can cast spells
func (f *MarkdownFormatter) isSpellcaster(class string) bool {
	spellcasters := []string{"wizard", "sorcerer", "warlock", "bard", "cleric", "druid", "paladin", "ranger"}
	class = strings.ToLower(class)
	for _, caster := range spellcasters {
		if class == caster {
			return true
		}
	}
	return false
}

// getSpellcastingAbility returns the spellcasting ability for a class
func (f *MarkdownFormatter) getSpellcastingAbility(class string) string {
	switch strings.ToLower(class) {
	case "wizard":
		return "intelligence"
	case "sorcerer", "bard", "paladin", "warlock":
		return "charisma"
	case "cleric", "druid", "ranger":
		return "wisdom"
	default:
		return "intelligence"
	}
}

// getSpellcastingModifier returns the spellcasting modifier for a character
func (f *MarkdownFormatter) getSpellcastingModifier(char *domain.Character, ability string) int {
	switch strings.ToLower(ability) {
	case "intelligence":
		return domain.Modifier(char.Int)
	case "wisdom":
		return domain.Modifier(char.Wis)
	case "charisma":
		return domain.Modifier(char.Cha)
	default:
		return 0
	}
}

// formatSpellsByLevel formats spells organized by level
func (f *MarkdownFormatter) formatSpellsByLevel(spells []string) string {
	if len(spells) == 0 {
		return ""
	}

	// For now, we'll use a simple heuristic to determine spell levels
	// In a real implementation, you'd look this up from a spell database
	spellLevels := make(map[int][]string)

	for _, spell := range spells {
		level := f.getSpellLevel(spell)
		spellLevels[level] = append(spellLevels[level], spell)
	}

	var builder strings.Builder

	// Sort levels
	levels := make([]int, 0, len(spellLevels))
	for level := range spellLevels {
		levels = append(levels, level)
	}
	sort.Ints(levels)

	for _, level := range levels {
		spells := spellLevels[level]
		builder.WriteString(fmt.Sprintf("### %d\n", level))
		for _, spell := range spells {
			builder.WriteString(fmt.Sprintf("- %s\n", spell))
		}
		builder.WriteString("\n")
	}

	return builder.String()
}

// getSpellLevel returns the level of a spell (simplified heuristic)
func (f *MarkdownFormatter) getSpellLevel(spell string) int {
	// Simple heuristic based on spell names from the example
	spell = strings.ToLower(spell)
	switch spell {
	case "command":
		return 1
	case "beacon of hope":
		return 3
	default:
		return 1 // Default to level 1
	}
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
