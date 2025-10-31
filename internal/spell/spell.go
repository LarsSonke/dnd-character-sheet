package spell

import (
	"encoding/csv"
	"os"
	"strconv"
)

// Spell represents basic spell information from CSV
type Spell struct {
	Name  string
	Level string
	Class string
}

// EnrichedSpell represents spell with additional API data
type EnrichedSpell struct {
	Spell
	// API-enriched fields
	School        string   `json:"school,omitempty"`
	Range         string   `json:"range,omitempty"`
	Components    []string `json:"components,omitempty"`
	Duration      string   `json:"duration,omitempty"`
	CastingTime   string   `json:"casting_time,omitempty"`
	Description   []string `json:"description,omitempty"`
	HigherLevel   []string `json:"higher_level,omitempty"`
	Ritual        bool     `json:"ritual,omitempty"`
	Concentration bool     `json:"concentration,omitempty"`
	LevelInt      int      `json:"level_int,omitempty"`
}

// ToEnriched converts a basic Spell to EnrichedSpell
func (s Spell) ToEnriched() EnrichedSpell {
	enriched := EnrichedSpell{Spell: s}

	// Convert string level to int if possible
	if level, err := strconv.Atoi(s.Level); err == nil {
		enriched.LevelInt = level
	}

	return enriched
}

func LoadSpellsFromCSV(path string) ([]Spell, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var spells []Spell
	for i, rec := range records {
		if i == 0 {
			continue // skip header
		}
		if len(rec) < 3 {
			continue // skip incomplete or empty rows
		}
		spells = append(spells, Spell{
			Name:  rec[0],
			Level: rec[1],
			Class: rec[2],
		})
	}
	return spells, nil
}
