package cli

import (
	"DnD-sheet/internal/equipment"
	"DnD-sheet/internal/spell"
	"fmt"
)

// APITestCommand tests the D&D 5e API integration
type APITestCommand struct {
	*BaseCommand

	// Flags
	spellTest     *bool
	equipmentTest *bool
	spellClass    *string
	equipmentType *string
	limit         *int
	query         *string
}

// NewAPITestCommand creates a new API test command
func NewAPITestCommand() *APITestCommand {
	cmd := &APITestCommand{
		BaseCommand: NewBaseCommand("api-test"),
	}

	// Define flags
	cmd.spellTest = cmd.flagSet.Bool("spells", false, "test spell enrichment")
	cmd.equipmentTest = cmd.flagSet.Bool("equipment", false, "test equipment enrichment")
	cmd.spellClass = cmd.flagSet.String("class", "wizard", "class for spell testing")
	cmd.equipmentType = cmd.flagSet.String("type", "weapon", "equipment type (weapon/armor)")
	cmd.limit = cmd.flagSet.Int("limit", 8, "limit number of items to test (max 10 for development)")
	cmd.query = cmd.flagSet.String("query", "", "search query for items")

	return cmd
}

// Name returns the command name
func (c *APITestCommand) Name() string {
	return "api-test"
}

// Execute runs the API test
func (c *APITestCommand) Execute() error {
	// Limit to max 10 items for development testing
	if *c.limit > 10 {
		*c.limit = 10
		fmt.Printf("Limiting to 10 items for development testing\n")
	}

	if *c.spellTest {
		return c.testSpellEnrichment()
	}

	if *c.equipmentTest {
		return c.testEquipmentEnrichment()
	}

	// Default: show usage
	fmt.Printf("API Test Command - Test D&D 5e API integration\n")
	fmt.Printf("Usage examples:\n")
	fmt.Printf("  ./dndcsg api-test -spells -class wizard -limit 5\n")
	fmt.Printf("  ./dndcsg api-test -equipment -type weapon -limit 8\n")
	fmt.Printf("  ./dndcsg api-test -spells -query \"magic missile\"\n")
	fmt.Printf("\nThis tests API integration with small batches (respecting 5-10 req/sec rate limit)\n")

	return nil
}

func (c *APITestCommand) testSpellEnrichment() error {
	fmt.Printf("Testing Spell Enrichment (D&D 5e API)\n")
	fmt.Printf("Rate limiting: 8 requests per second\n")
	fmt.Printf("Limit: %d spells\n\n", *c.limit)

	// Create spell enrichment service
	spellService := spell.NewEnrichmentService()
	defer spellService.Close()

	csvPath := "internal/spell/5e-SRD-Spells.csv"

	var enrichedSpells []spell.EnrichedSpell
	var err error

	if *c.query != "" {
		// Search for specific spells
		fmt.Printf("Searching for spells matching '%s'...\n", *c.query)
		enrichedSpells, err = spellService.SearchSpells(csvPath, *c.query, *c.limit)
	} else {
		// Get spells by class
		fmt.Printf("Getting %s spells...\n", *c.spellClass)
		enrichedSpells, err = spellService.GetSpellsByClass(csvPath, *c.spellClass, *c.limit)
	}

	if err != nil {
		return fmt.Errorf("failed to enrich spells: %w", err)
	}

	if len(enrichedSpells) == 0 {
		fmt.Printf("No spells found\n")
		return nil
	}

	fmt.Printf("Successfully enriched %d spells:\n\n", len(enrichedSpells))

	for i, spell := range enrichedSpells {
		fmt.Printf("=== Spell %d: %s ===\n", i+1, spell.Name)
		fmt.Printf("Level: %d\n", spell.LevelInt)
		fmt.Printf("Class: %s\n", spell.Class)

		if spell.School != "" {
			fmt.Printf("School: %s\n", spell.School)
		}
		if spell.Range != "" {
			fmt.Printf("Range: %s\n", spell.Range)
		}
		if spell.CastingTime != "" {
			fmt.Printf("Casting Time: %s\n", spell.CastingTime)
		}
		if spell.Duration != "" {
			fmt.Printf("Duration: %s\n", spell.Duration)
		}
		if len(spell.Components) > 0 {
			fmt.Printf("Components: %v\n", spell.Components)
		}
		if spell.Ritual {
			fmt.Printf("Ritual: Yes\n")
		}
		if spell.Concentration {
			fmt.Printf("Concentration: Yes\n")
		}
		if len(spell.Description) > 0 && len(spell.Description[0]) > 0 {
			desc := spell.Description[0]
			if len(desc) > 100 {
				desc = desc[:100] + "..."
			}
			fmt.Printf("Description: %s\n", desc)
		}
		fmt.Println()
	}

	// Show summary
	apiEnrichedCount := 0
	for _, spell := range enrichedSpells {
		if spell.School != "" || spell.Range != "" {
			apiEnrichedCount++
		}
	}

	fmt.Printf("Summary:\n")
	fmt.Printf("- Total spells: %d\n", len(enrichedSpells))
	fmt.Printf("- API enriched: %d\n", apiEnrichedCount)
	fmt.Printf("- Success rate: %.1f%%\n", float64(apiEnrichedCount)/float64(len(enrichedSpells))*100)

	return nil
}

func (c *APITestCommand) testEquipmentEnrichment() error {
	fmt.Printf("Testing Equipment Enrichment (D&D 5e API)\n")
	fmt.Printf("Rate limiting: 8 requests per second\n")
	fmt.Printf("Limit: %d items\n\n", *c.limit)

	// Create equipment enrichment service
	equipmentService := equipment.NewEnrichmentService()
	defer equipmentService.Close()

	csvPath := "internal/equipment/5e-SRD-Equipment.csv"

	var enrichedEquipment []equipment.EnrichedEquipment
	var err error

	if *c.query != "" {
		// Search for specific equipment
		fmt.Printf("Searching for equipment matching '%s'...\n", *c.query)
		enrichedEquipment, err = equipmentService.SearchEquipment(csvPath, *c.query, *c.limit)
	} else if *c.equipmentType == "weapon" {
		// Get weapons
		fmt.Printf("Getting weapons...\n")
		enrichedEquipment, err = equipmentService.GetWeapons(csvPath, *c.limit)
	} else if *c.equipmentType == "armor" {
		// Get armor
		fmt.Printf("Getting armor...\n")
		enrichedEquipment, err = equipmentService.GetArmor(csvPath, *c.limit)
	} else {
		// Search by type
		enrichedEquipment, err = equipmentService.SearchEquipment(csvPath, *c.equipmentType, *c.limit)
	}

	if err != nil {
		return fmt.Errorf("failed to enrich equipment: %w", err)
	}

	if len(enrichedEquipment) == 0 {
		fmt.Printf("No equipment found\n")
		return nil
	}

	fmt.Printf("Successfully processed %d equipment items:\n\n", len(enrichedEquipment))

	for i, eq := range enrichedEquipment {
		fmt.Printf("=== Equipment %d: %s ===\n", i+1, eq.Name)
		fmt.Printf("Category: %s\n", eq.Category)

		// Show weapon-specific data
		if eq.WeaponCategory != "" {
			fmt.Printf("Weapon Category: %s\n", eq.WeaponCategory)
			fmt.Printf("Weapon Range: %s\n", eq.WeaponRange)
			if eq.Range.Normal > 0 {
				fmt.Printf("Range: %d", eq.Range.Normal)
				if eq.Range.Long > 0 {
					fmt.Printf("/%d", eq.Range.Long)
				}
				fmt.Println(" ft")
			}
			if eq.TwoHanded {
				fmt.Printf("Two-Handed: Yes\n")
			}
			if eq.Damage != "" {
				fmt.Printf("Damage: %s", eq.Damage)
				if eq.DamageType != "" {
					fmt.Printf(" %s", eq.DamageType)
				}
				fmt.Println()
			}
			if len(eq.Properties) > 0 {
				fmt.Printf("Properties: %v\n", eq.Properties)
			}
		}

		// Show armor-specific data
		if eq.ArmorCategory != "" {
			fmt.Printf("Armor Category: %s\n", eq.ArmorCategory)
			fmt.Printf("Armor Class: %d", eq.ArmorClass.Base)
			if eq.ArmorClass.DexBonus {
				if eq.ArmorClass.MaxBonus > 0 {
					fmt.Printf(" + Dex (max %d)", eq.ArmorClass.MaxBonus)
				} else {
					fmt.Printf(" + Dex")
				}
			}
			fmt.Println()
			if eq.StrMinimum > 0 {
				fmt.Printf("Strength Requirement: %d\n", eq.StrMinimum)
			}
			if eq.StealthDisadvantage {
				fmt.Printf("Stealth Disadvantage: Yes\n")
			}
		}

		fmt.Println()
	}

	// Show summary
	apiEnrichedCount := 0
	for _, eq := range enrichedEquipment {
		if eq.WeaponCategory != "" || eq.ArmorCategory != "" {
			apiEnrichedCount++
		}
	}

	fmt.Printf("Summary:\n")
	fmt.Printf("- Total equipment: %d\n", len(enrichedEquipment))
	fmt.Printf("- API enriched: %d\n", apiEnrichedCount)
	fmt.Printf("- Success rate: %.1f%%\n", float64(apiEnrichedCount)/float64(len(enrichedEquipment))*100)

	return nil
}
