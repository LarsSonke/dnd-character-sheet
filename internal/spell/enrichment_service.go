package spell

import (
	"DnD-sheet/internal/api"
	"fmt"
	"log"
	"strings"
)

// EnrichmentService handles enriching spells with API data
type EnrichmentService struct {
	apiClient *api.Client
}

// NewEnrichmentService creates a new spell enrichment service
func NewEnrichmentService() *EnrichmentService {
	return &EnrichmentService{
		apiClient: api.NewClient(),
	}
}

// Close closes the API client
func (s *EnrichmentService) Close() {
	if s.apiClient != nil {
		s.apiClient.Close()
	}
}

// EnrichSpell enriches a single spell with API data
func (s *EnrichmentService) EnrichSpell(spell Spell) EnrichedSpell {
	enriched := spell.ToEnriched()

	// Fetch additional data from API
	spellDetails, err := s.apiClient.GetSpell(spell.Name)
	if err != nil {
		log.Printf("Failed to enrich spell '%s': %v", spell.Name, err)
		return enriched
	}

	// Map API data to enriched spell
	enriched.School = spellDetails.School.Name
	enriched.Range = spellDetails.Range
	enriched.Components = spellDetails.Components
	enriched.Duration = spellDetails.Duration
	enriched.CastingTime = spellDetails.CastingTime
	enriched.Description = spellDetails.Description
	enriched.HigherLevel = spellDetails.HigherLevel
	enriched.Ritual = spellDetails.Ritual
	enriched.Concentration = spellDetails.Concentration
	enriched.LevelInt = spellDetails.Level

	return enriched
}

// EnrichSpellsBatch enriches multiple spells concurrently
func (s *EnrichmentService) EnrichSpellsBatch(spells []Spell) []EnrichedSpell {
	if len(spells) == 0 {
		return nil
	}

	// Extract spell names for batch request
	spellNames := make([]string, len(spells))
	spellMap := make(map[string]Spell)

	for i, spell := range spells {
		spellNames[i] = spell.Name
		spellMap[spell.Name] = spell
	}

	// Make batch API request
	results := s.apiClient.GetSpellsBatch(spellNames)

	// Process results
	enrichedSpells := make([]EnrichedSpell, 0, len(spells))

	for _, result := range results {
		baseSpell := spellMap[result.Name]
		enriched := baseSpell.ToEnriched()

		if result.Error != nil {
			log.Printf("Failed to enrich spell '%s': %v", result.Name, result.Error)
			enrichedSpells = append(enrichedSpells, enriched)
			continue
		}

		if spellDetails, ok := result.Data.(*api.SpellDetails); ok {
			// Map API data to enriched spell
			enriched.School = spellDetails.School.Name
			enriched.Range = spellDetails.Range
			enriched.Components = spellDetails.Components
			enriched.Duration = spellDetails.Duration
			enriched.CastingTime = spellDetails.CastingTime
			enriched.Description = spellDetails.Description
			enriched.HigherLevel = spellDetails.HigherLevel
			enriched.Ritual = spellDetails.Ritual
			enriched.Concentration = spellDetails.Concentration
			enriched.LevelInt = spellDetails.Level
		}

		enrichedSpells = append(enrichedSpells, enriched)
	}

	return enrichedSpells
}

// SearchSpells searches for spells by name or class
func (s *EnrichmentService) SearchSpells(csvPath string, query string, limit int) ([]EnrichedSpell, error) {
	// Load spells from CSV
	spells, err := LoadSpellsFromCSV(csvPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load spells: %w", err)
	}

	// Filter spells based on query
	var matchedSpells []Spell
	queryLower := strings.ToLower(query)

	for _, spell := range spells {
		if strings.Contains(strings.ToLower(spell.Name), queryLower) ||
			strings.Contains(strings.ToLower(spell.Class), queryLower) {
			matchedSpells = append(matchedSpells, spell)

			// Limit results for testing
			if len(matchedSpells) >= limit {
				break
			}
		}
	}

	// Enrich matched spells
	return s.EnrichSpellsBatch(matchedSpells), nil
}

// GetSpellsByClass returns spells for a specific class with enrichment
func (s *EnrichmentService) GetSpellsByClass(csvPath string, className string, limit int) ([]EnrichedSpell, error) {
	// Load spells from CSV
	spells, err := LoadSpellsFromCSV(csvPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load spells: %w", err)
	}

	// Filter spells by class
	var classSpells []Spell
	classLower := strings.ToLower(className)

	for _, spell := range spells {
		if strings.Contains(strings.ToLower(spell.Class), classLower) {
			classSpells = append(classSpells, spell)

			// Limit results for testing
			if len(classSpells) >= limit {
				break
			}
		}
	}

	// Enrich class spells
	return s.EnrichSpellsBatch(classSpells), nil
}
