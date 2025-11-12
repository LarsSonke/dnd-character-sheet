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
		if char.IsSpellcaster() {
			spellAbility := char.SpellcastingAbility()
			// Convert INT/WIS/CHA to full name for display
			var abilityName string
			switch spellAbility {
			case "INT":
				abilityName = "intelligence"
			case "WIS":
				abilityName = "wisdom"
			case "CHA":
				abilityName = "charisma"
			default:
				abilityName = "intelligence"
			}

			fmt.Printf("Spellcasting ability: %s\n", abilityName)
			fmt.Printf("Spell save DC: %d\n", char.SpellSaveDC())
			fmt.Printf("Spell attack bonus: +%d\n", char.SpellAttackBonus())
		}
		
		// Print known spells if the character has any
		if len(char.KnownSpells) > 0 {
			fmt.Println("Known spells:")
			for _, spell := range char.KnownSpells {
				fmt.Printf("  - %s\n", spell)
			}
		}
		
		// Print prepared spells if the character has any
		if len(char.PreparedSpells) > 0 {
			fmt.Println("Prepared spells:")
			for _, spell := range char.PreparedSpells {
				fmt.Printf("  - %s\n", spell)
			}
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
	fmt.Printf("Armor class: %d\n", char.ArmorClass())
	fmt.Printf("Initiative bonus: %d\n", char.Initiative())
	fmt.Printf("Passive perception: %d\n", char.PassivePerception())
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
