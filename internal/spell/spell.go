package spell

import (
	"encoding/csv"
	"os"
	"strconv"
	"strings"
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

// GetSpellLevel returns the spell level for common D&D 5e spells
// This is a reference table based on D&D 5e System Reference Document
func GetSpellLevel(spellName string) int {
	spellLevels := map[string]int{
		// Cantrips (Level 0)
		"fire bolt": 0, "mage hand": 0, "prestidigitation": 0, "light": 0,

		// Level 1 spells
		"burning hands": 1, "magic missile": 1, "cure wounds": 1, "shield": 1, "false life": 1, "feather fall": 1, "command": 1,

		// Level 2 spells
		"scorching ray": 2, "misty step": 2, "web": 2, "hold person": 2,

		// Level 3 spells
		"fireball": 3, "lightning bolt": 3, "counterspell": 3, "fly": 3, "beacon of hope": 3,

		// Level 4 spells
		"wall of fire": 4, "dimension door": 4, "polymorph": 4,

		// Level 5 spells
		"cone of cold": 5, "teleport": 5, "wall of stone": 5,

		// Level 6 spells
		"disintegrate": 6, "mass suggestion": 6, "wall of ice": 6,

		// Level 7 spells
		"etherealness": 7, "fire storm": 7, "plane shift": 7,

		// Level 8 spells
		"power word stun": 8, "maze": 8, "sunburst": 8,

		// Level 9 spells
		"wish": 9, "meteor swarm": 9, "time stop": 9,
	}

	spellLower := strings.ToLower(spellName)
	if level, exists := spellLevels[spellLower]; exists {
		return level
	}

	// Default to level 1 for unknown spells to be safe
	return 1
}
