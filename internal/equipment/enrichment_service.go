package equipment

import (
	"DnD-sheet/internal/api"
	"fmt"
	"log"
	"strings"
)

// EnrichmentService handles enriching equipment with API data
type EnrichmentService struct {
	apiClient *api.Client
}

// NewEnrichmentService creates a new equipment enrichment service
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

// EnrichEquipment enriches a single equipment item with API data
func (s *EnrichmentService) EnrichEquipment(equipment Equipment) EnrichedEquipment {
	enriched := equipment.ToEnriched()

	// Fetch additional data from API
	equipmentData, err := s.apiClient.GetEquipment(equipment.Name)
	if err != nil {
		log.Printf("Failed to enrich equipment '%s': %v", equipment.Name, err)
		return enriched
	}

	// Map API data based on equipment type
	switch data := equipmentData.(type) {
	case *api.WeaponDetails:
		// Enrich weapon data
		enriched.WeaponCategory = data.WeaponCategory
		enriched.WeaponRange = data.WeaponRange
		enriched.Range = WeaponRange{
			Normal: data.Range.Normal,
			Long:   data.Range.Long,
		}
		enriched.Damage = data.Damage.DamageDice
		enriched.DamageType = data.Damage.DamageType.Name

		// Extract properties
		properties := make([]string, len(data.Properties))
		for i, prop := range data.Properties {
			properties[i] = prop.Name
			// Check for two-handed property
			if strings.ToLower(prop.Name) == "two-handed" {
				enriched.TwoHanded = true
			}
		}
		enriched.Properties = properties

	case *api.ArmorDetails:
		// Enrich armor data
		enriched.ArmorCategory = data.ArmorCategory
		enriched.ArmorClass = ArmorClass{
			Base:     data.ArmorClass.Base,
			DexBonus: data.ArmorClass.DexBonus,
			MaxBonus: data.ArmorClass.MaxBonus,
		}
		enriched.StrMinimum = data.StrMinimum
		enriched.StealthDisadvantage = data.StealthDisadvantage

	default:
		log.Printf("Unknown equipment type for '%s'", equipment.Name)
	}

	return enriched
}

// EnrichEquipmentBatch enriches multiple equipment items concurrently
func (s *EnrichmentService) EnrichEquipmentBatch(equipment []Equipment) []EnrichedEquipment {
	if len(equipment) == 0 {
		return nil
	}

	// Extract equipment names for batch request
	equipmentNames := make([]string, len(equipment))
	equipmentMap := make(map[string]Equipment)

	for i, eq := range equipment {
		equipmentNames[i] = eq.Name
		equipmentMap[eq.Name] = eq
	}

	// Make batch API request
	results := s.apiClient.GetEquipmentBatch(equipmentNames)

	// Process results
	enrichedEquipment := make([]EnrichedEquipment, 0, len(equipment))

	for _, result := range results {
		baseEquipment := equipmentMap[result.Name]
		enriched := baseEquipment.ToEnriched()

		if result.Error != nil {
			log.Printf("Failed to enrich equipment '%s': %v", result.Name, result.Error)
			enrichedEquipment = append(enrichedEquipment, enriched)
			continue
		}

		// Map API data based on equipment type
		switch data := result.Data.(type) {
		case *api.WeaponDetails:
			// Enrich weapon data
			enriched.WeaponCategory = data.WeaponCategory
			enriched.WeaponRange = data.WeaponRange
			enriched.Range = WeaponRange{
				Normal: data.Range.Normal,
				Long:   data.Range.Long,
			}
			enriched.Damage = data.Damage.DamageDice
			enriched.DamageType = data.Damage.DamageType.Name

			// Extract properties
			properties := make([]string, len(data.Properties))
			for i, prop := range data.Properties {
				properties[i] = prop.Name
				// Check for two-handed property
				if strings.ToLower(prop.Name) == "two-handed" {
					enriched.TwoHanded = true
				}
			}
			enriched.Properties = properties

		case *api.ArmorDetails:
			// Enrich armor data
			enriched.ArmorCategory = data.ArmorCategory
			enriched.ArmorClass = ArmorClass{
				Base:     data.ArmorClass.Base,
				DexBonus: data.ArmorClass.DexBonus,
				MaxBonus: data.ArmorClass.MaxBonus,
			}
			enriched.StrMinimum = data.StrMinimum
			enriched.StealthDisadvantage = data.StealthDisadvantage
		}

		enrichedEquipment = append(enrichedEquipment, enriched)
	}

	return enrichedEquipment
}

// SearchEquipment searches for equipment by name or category
func (s *EnrichmentService) SearchEquipment(csvPath string, query string, limit int) ([]EnrichedEquipment, error) {
	// Load equipment from CSV
	equipment, err := LoadEquipmentFromCSV(csvPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load equipment: %w", err)
	}

	// Filter equipment based on query
	var matchedEquipment []Equipment
	queryLower := strings.ToLower(query)

	for _, eq := range equipment {
		if strings.Contains(strings.ToLower(eq.Name), queryLower) ||
			strings.Contains(strings.ToLower(eq.Category), queryLower) {
			matchedEquipment = append(matchedEquipment, eq)

			// Limit results for testing
			if len(matchedEquipment) >= limit {
				break
			}
		}
	}

	// Enrich matched equipment
	return s.EnrichEquipmentBatch(matchedEquipment), nil
}

// GetWeapons returns weapons with enrichment
func (s *EnrichmentService) GetWeapons(csvPath string, limit int) ([]EnrichedEquipment, error) {
	// Load equipment from CSV
	equipment, err := LoadEquipmentFromCSV(csvPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load equipment: %w", err)
	}

	// Filter weapons
	var weapons []Equipment
	for _, eq := range equipment {
		if strings.Contains(strings.ToLower(eq.Category), "weapon") {
			weapons = append(weapons, eq)

			// Limit results for testing
			if len(weapons) >= limit {
				break
			}
		}
	}

	// Enrich weapons
	return s.EnrichEquipmentBatch(weapons), nil
}

// GetArmor returns armor with enrichment
func (s *EnrichmentService) GetArmor(csvPath string, limit int) ([]EnrichedEquipment, error) {
	// Load equipment from CSV
	equipment, err := LoadEquipmentFromCSV(csvPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load equipment: %w", err)
	}

	// Filter armor
	var armor []Equipment
	for _, eq := range equipment {
		if strings.Contains(strings.ToLower(eq.Category), "armor") {
			armor = append(armor, eq)

			// Limit results for testing
			if len(armor) >= limit {
				break
			}
		}
	}

	// Enrich armor
	return s.EnrichEquipmentBatch(armor), nil
}
