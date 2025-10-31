package infrastructure

import (
	"DnD-sheet/internal/character/domain"
	"encoding/json"
	"os"
	"path/filepath"
)

// JSONCharacterRepository implements character persistence using JSON files
type JSONCharacterRepository struct {
	dataDir string
}

// NewJSONCharacterRepository creates a new JSON repository
func NewJSONCharacterRepository(dataDir string) *JSONCharacterRepository {
	return &JSONCharacterRepository{dataDir: dataDir}
}

// Save persists a character to a JSON file
func (r *JSONCharacterRepository) Save(character *domain.Character) error {
	if err := os.MkdirAll(r.dataDir, 0755); err != nil {
		return err
	}
	path := filepath.Join(r.dataDir, character.Name+".json")
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(character)
}

// Load retrieves a character by name from a JSON file
func (r *JSONCharacterRepository) Load(name string) (*domain.Character, error) {
	path := filepath.Join(r.dataDir, name+".json")
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var c domain.Character
	if err := json.NewDecoder(f).Decode(&c); err != nil {
		return nil, err
	}
	return &c, nil
}

// Delete removes a character's JSON file
func (r *JSONCharacterRepository) Delete(name string) error {
	path := filepath.Join(r.dataDir, name+".json")
	return os.Remove(path)
}

// List returns all character names by reading JSON files in the data directory
func (r *JSONCharacterRepository) List() ([]string, error) {
	var names []string
	files, err := os.ReadDir(r.dataDir)
	if err != nil {
		return names, err
	}
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			names = append(names, file.Name()[:len(file.Name())-5])
		}
	}
	return names, nil
}

// Exists checks if a character file exists
func (r *JSONCharacterRepository) Exists(name string) bool {
	path := filepath.Join(r.dataDir, name+".json")
	_, err := os.Stat(path)
	return err == nil
}
