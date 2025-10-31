package cli

import (
	"DnD-sheet/internal/character/domain"
	"fmt"
	"strings"
)

// printCharacterInfo prints character information in the expected format
func (c *ViewCommand) printCharacterInfo(char *domain.Character) {
	// Print basic info
	fmt.Printf("Name: %s\n", char.Name)
	fmt.Printf("Class: %s\n", strings.ToLower(char.Class))
	fmt.Printf("Race: %s\n", strings.ToLower(char.Race))
	fmt.Printf("Background: %s\n", strings.ToLower(char.Background))
	fmt.Printf("Level: %d\n", char.Level)

	// Print ability scores
	fmt.Println("Ability scores:")
	fmt.Printf("  STR: %d (%+d)\n", char.Str, domain.Modifier(char.Str))
	fmt.Printf("  DEX: %d (%+d)\n", char.Dex, domain.Modifier(char.Dex))
	fmt.Printf("  CON: %d (%+d)\n", char.Con, domain.Modifier(char.Con))
	fmt.Printf("  INT: %d (%+d)\n", char.Int, domain.Modifier(char.Int))
	fmt.Printf("  WIS: %d (%+d)\n", char.Wis, domain.Modifier(char.Wis))
	fmt.Printf("  CHA: %d (%+d)\n", char.Cha, domain.Modifier(char.Cha))

	// Print proficiency bonus (lowercase 'bonus')
	fmt.Printf("Proficiency bonus: %+d\n", char.ProficiencyBonus)

	// Print skill proficiencies (just the names, comma-separated)
	if len(char.SkillProficiencies) > 0 {
		fmt.Printf("Skill proficiencies: %s\n", strings.Join(char.SkillProficiencies, ", "))
	}

	// Print spell slots if the character has any
	if len(char.SpellSlots) > 0 {
		fmt.Println("Spell slots:")
		// Sort spell levels for consistent output
		levels := make([]int, 0, len(char.SpellSlots))
		for level := range char.SpellSlots {
			levels = append(levels, level)
		}
		// Simple sort
		for i := 0; i < len(levels)-1; i++ {
			for j := i + 1; j < len(levels); j++ {
				if levels[i] > levels[j] {
					levels[i], levels[j] = levels[j], levels[i]
				}
			}
		}
		for _, level := range levels {
			if char.SpellSlots[level] > 0 {
				fmt.Printf("  Level %d: %d\n", level, char.SpellSlots[level])
			}
		}

		// Print spellcasting stats if character can cast spells
		if c.canCastSpells(char.Class) {
			spellcastingAbility := c.getSpellcastingAbility(char.Class)
			spellcastingMod := c.getSpellcastingModifier(char, spellcastingAbility)
			spellSaveDC := 8 + char.ProficiencyBonus + spellcastingMod
			spellAttackBonus := char.ProficiencyBonus + spellcastingMod

			fmt.Printf("Spellcasting ability: %s\n", spellcastingAbility)
			fmt.Printf("Spell save DC: %d\n", spellSaveDC)
			fmt.Printf("Spell attack bonus: +%d\n", spellAttackBonus)
		}
	}

	// Print equipment information
	if char.Weapon != "" {
		weaponSlot := "main hand"
		if char.WeaponSlot != "" {
			weaponSlot = char.WeaponSlot
		}
		// Capitalize only the first letter, not every word
		capitalizedSlot := strings.ToUpper(string(weaponSlot[0])) + weaponSlot[1:]
		fmt.Printf("%s: %s\n", capitalizedSlot, char.Weapon)
	}
	if char.Armor != "" {
		fmt.Printf("Armor: %s\n", char.Armor)
	}
	if char.Shield != "" {
		fmt.Printf("Shield: %s\n", char.Shield)
	}

	// Print calculated stats
	fmt.Printf("Armor class: %d\n", c.calculateArmorClass(char))
	fmt.Printf("Initiative bonus: %d\n", c.calculateInitiative(char))
	fmt.Printf("Passive perception: %d\n", c.calculatePassivePerception(char))
}

// Helper methods for calculating character stats

// getAbilityModifier calculates the ability modifier from an ability score
func (c *ViewCommand) getAbilityModifier(abilityScore int) int {
	modifier := (abilityScore - 10) / 2
	// Handle negative modifiers correctly (D&D rounds down, not towards zero)
	if abilityScore < 10 && (abilityScore-10)%2 != 0 {
		modifier-- // Round down for odd negative numbers
	}
	return modifier
}

// calculateArmorClass calculates AC based on armor and dexterity
func (c *ViewCommand) calculateArmorClass(char *domain.Character) int {
	baseAC := 10 // Base AC without armor
	dexMod := c.getAbilityModifier(char.Dex)

	// Calculate AC based on armor type
	if char.Armor != "" {
		switch strings.ToLower(char.Armor) {
		// Light Armor (Dex modifier applies fully)
		case "padded", "padded armor":
			baseAC = 11 + dexMod
		case "leather", "leather armor":
			baseAC = 11 + dexMod
		case "studded leather", "studded leather armor":
			baseAC = 12 + dexMod

		// Medium Armor (Dex modifier max +2)
		case "hide", "hide armor":
			baseAC = 12 + min(dexMod, 2)
		case "chain shirt":
			baseAC = 13 + min(dexMod, 2)
		case "scale mail", "scale mail armor":
			baseAC = 14 + min(dexMod, 2)
		case "breastplate":
			baseAC = 14 + min(dexMod, 2)
		case "half plate", "half plate armor":
			baseAC = 15 + min(dexMod, 2)

		// Heavy Armor (No Dex modifier)
		case "ring mail", "ring mail armor":
			baseAC = 14
		case "chain mail", "chain mail armor":
			baseAC = 16
		case "splint", "splint armor":
			baseAC = 17
		case "plate", "plate armor":
			baseAC = 18

		// Natural Armor (for some races/classes)
		case "natural armor":
			baseAC = 13 + dexMod

		// Magical or Special Armor
		case "mage armor":
			baseAC = 13 + dexMod
		case "barkskin":
			baseAC = max(16, 10+dexMod) // Barkskin sets AC to 16 minimum

		default:
			// Try to guess based on name patterns
			armorLower := strings.ToLower(char.Armor)
			if strings.Contains(armorLower, "leather") {
				baseAC = 11 + dexMod // Default to leather armor
			} else if strings.Contains(armorLower, "chain") || strings.Contains(armorLower, "mail") {
				baseAC = 16 // Default to chain mail
			} else if strings.Contains(armorLower, "plate") {
				baseAC = 18 // Default to plate
			} else {
				baseAC = 11 + dexMod // Default light armor
			}
		}
	} else {
		// No armor - check for class-based AC calculations
		classLower := strings.ToLower(char.Class)
		switch classLower {
		case "barbarian":
			// Unarmored Defense: 10 + Dex + Con
			conMod := c.getAbilityModifier(char.Con)
			baseAC = 10 + dexMod + conMod
		case "monk":
			// Unarmored Defense: 10 + Dex + Wis
			wisMod := c.getAbilityModifier(char.Wis)
			baseAC = 10 + dexMod + wisMod
		case "sorcerer":
			// Draconic Bloodline gets 13 + Dex (simplified - assume draconic)
			baseAC = max(13+dexMod, 10+dexMod)
		default:
			baseAC = 10 + dexMod // Standard unarmored AC
		}
	}

	// Add shield bonus
	if char.Shield != "" {
		shieldLower := strings.ToLower(char.Shield)
		switch shieldLower {
		case "shield":
			baseAC += 2
		case "buckler":
			baseAC += 1
		case "+1 shield", "shield +1":
			baseAC += 3
		case "+2 shield", "shield +2":
			baseAC += 4
		case "+3 shield", "shield +3":
			baseAC += 5
		default:
			baseAC += 2 // Default shield bonus
		}
	}

	// Consider weapon-based AC bonuses (rare but possible)
	if char.Weapon != "" {
		weaponLower := strings.ToLower(char.Weapon)
		if strings.Contains(weaponLower, "defensive") || strings.Contains(weaponLower, "+1 ac") {
			baseAC += 1
		}
	}

	return baseAC
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
} // calculateInitiative calculates initiative bonus (Dex mod + potential class bonuses)
func (c *ViewCommand) calculateInitiative(char *domain.Character) int {
	dexMod := c.getAbilityModifier(char.Dex)
	initiative := dexMod

	// Class-specific initiative bonuses
	classLower := strings.ToLower(char.Class)
	switch classLower {
	case "rogue":
		// Some rogue subclasses get initiative bonuses, but we'll keep it simple
		// Could add half proficiency bonus for certain levels
	case "barbarian":
		// Feral Instinct at higher levels, but keeping simple
	case "bard":
		// Jack of All Trades adds half proficiency to initiative (simplified)
		if char.Level >= 2 {
			initiative += char.ProficiencyBonus / 2
		}
	}

	// Alert feat would add +5, but we don't track feats yet

	return initiative
}

// calculatePassivePerception calculates passive perception (10 + Wis mod + proficiency if proficient)
func (c *ViewCommand) calculatePassivePerception(char *domain.Character) int {
	wisMod := c.getAbilityModifier(char.Wis)
	passivePerception := 10 + wisMod

	// Check if character is proficient in Perception
	for _, skill := range char.SkillProficiencies {
		if strings.ToLower(skill) == "perception" {
			passivePerception += char.ProficiencyBonus
			break
		}
	}

	// Observant feat would add +5, but we don't track feats yet

	return passivePerception
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// canCastSpells returns true if the given class can cast spells
func (c *ViewCommand) canCastSpells(class string) bool {
	spellcastingClasses := []string{
		"wizard", "sorcerer", "warlock", "cleric", "druid", "paladin", "ranger",
		"bard", "artificer", "eldritch knight", "arcane trickster",
	}

	classLower := strings.ToLower(class)
	for _, spellClass := range spellcastingClasses {
		if classLower == spellClass {
			return true
		}
	}
	return false
}

// getSpellcastingAbility returns the spellcasting ability for a class
func (c *ViewCommand) getSpellcastingAbility(class string) string {
	classLower := strings.ToLower(class)
	switch classLower {
	case "wizard", "eldritch knight", "arcane trickster", "artificer":
		return "intelligence"
	case "cleric", "druid", "ranger":
		return "wisdom"
	case "sorcerer", "bard", "warlock", "paladin":
		return "charisma"
	default:
		return "intelligence" // Default
	}
}

// getSpellcastingModifier returns the spellcasting modifier for a character
func (c *ViewCommand) getSpellcastingModifier(char *domain.Character, ability string) int {
	switch strings.ToLower(ability) {
	case "intelligence":
		return c.getAbilityModifier(char.Int)
	case "wisdom":
		return c.getAbilityModifier(char.Wis)
	case "charisma":
		return c.getAbilityModifier(char.Cha)
	default:
		return 0
	}
}
