package domain

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
