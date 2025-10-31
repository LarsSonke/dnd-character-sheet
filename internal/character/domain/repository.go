package domain

// CharacterRepository defines the interface for character data persistence
type CharacterRepository interface {
	// Save persists a character to storage
	Save(character *Character) error

	// Load retrieves a character by name from storage
	Load(name string) (*Character, error)

	// Delete removes a character from storage
	Delete(name string) error

	// List returns all character names in storage
	List() ([]string, error)

	// Exists checks if a character with the given name exists
	Exists(name string) bool
}
