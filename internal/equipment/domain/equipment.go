package domain

import "strings"

// ArmorClass represents armor class statistics
type ArmorClass struct {
	Base     int  `json:"base"`
	DexBonus bool `json:"dex_bonus"`
}

// Equipment represents a piece of equipment in D&D
type Equipment struct {
	Name       string     `json:"name"`
	Category   string     `json:"category"`
	ArmorClass ArmorClass `json:"armor_class"`
	// Additional fields can be added here as needed
}

// EquipmentRepository defines the interface for equipment data access
type EquipmentRepository interface {
	// LoadAll loads all equipment from the data source
	LoadAll() ([]Equipment, error)

	// FindByName finds equipment by name (case-insensitive)
	FindByName(name string) (*Equipment, error)

	// FindByCategory returns all equipment in a specific category
	FindByCategory(category string) ([]Equipment, error)
}

// CalculateAC calculates the armor class with dexterity modifier according to D&D 5e armor rules
func (e *Equipment) CalculateAC(dexModifier int) int {
	ac := e.ArmorClass.Base

	if e.ArmorClass.DexBonus {
		// Apply dex bonus based on armor category (D&D 5e rules)
		switch strings.ToLower(e.Category) {
		case "medium armor":
			// Medium armor: max +2 dex bonus
			if dexModifier > 2 {
				ac += 2
			} else {
				ac += dexModifier
			}
		case "light armor":
			// Light armor: full dex bonus
			ac += dexModifier
		case "heavy armor":
			// Heavy armor: no dex bonus
		default:
			// Default: full dex bonus (for shields, etc.)
			ac += dexModifier
		}
	}

	return ac
}
