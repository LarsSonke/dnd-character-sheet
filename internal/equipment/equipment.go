package equipment

import (
	"encoding/csv"
	"os"
	"strconv"
	"strings"
)

type ArmorClass struct {
	Base     int  `json:"base"`
	DexBonus bool `json:"dex_bonus"`
	MaxBonus int  `json:"max_bonus,omitempty"`
}

type WeaponRange struct {
	Normal int `json:"normal,omitempty"`
	Long   int `json:"long,omitempty"`
}

type Equipment struct {
	Name       string     `json:"name"`
	Category   string     `json:"category"`
	ArmorClass ArmorClass `json:"armor_class"`
}

// EnrichedEquipment represents equipment with additional API data
type EnrichedEquipment struct {
	Equipment
	// API-enriched fields for weapons
	WeaponCategory string      `json:"weapon_category,omitempty"`
	WeaponRange    string      `json:"weapon_range,omitempty"`
	Range          WeaponRange `json:"range,omitempty"`
	TwoHanded      bool        `json:"two_handed,omitempty"`
	Damage         string      `json:"damage,omitempty"`
	DamageType     string      `json:"damage_type,omitempty"`
	Properties     []string    `json:"properties,omitempty"`

	// API-enriched fields for armor
	ArmorCategory       string `json:"armor_category,omitempty"`
	StrMinimum          int    `json:"str_minimum,omitempty"`
	StealthDisadvantage bool   `json:"stealth_disadvantage,omitempty"`
}

// ToEnriched converts basic Equipment to EnrichedEquipment
func (e Equipment) ToEnriched() EnrichedEquipment {
	return EnrichedEquipment{Equipment: e}
}

func LoadEquipmentFromCSV(path string) ([]Equipment, error) {
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
	var equipmentList []Equipment
	for i, rec := range records {
		if i == 0 {
			continue // skip header
		}
		if len(rec) < 2 {
			continue // skip incomplete rows
		}
		ac := ArmorClass{Base: 10} // default armor class
		if len(rec) > 3 {
			if base, err := strconv.Atoi(strings.TrimSpace(rec[2])); err == nil {
				ac.Base = base
			}
			if dexBonus := strings.TrimSpace(rec[3]); strings.EqualFold(dexBonus, "yes") {
				ac.DexBonus = true
			}
		}
		equipmentList = append(equipmentList, Equipment{
			Name:       strings.TrimSpace(rec[0]),
			Category:   strings.TrimSpace(rec[1]),
			ArmorClass: ac, // parsed from rec[2] and rec[3]
		})
	}
	return equipmentList, nil
}

func FindEquipmentByName(equipmentList []Equipment, name string) *Equipment {
	name = strings.ToLower(strings.TrimSpace(name))
	for _, eq := range equipmentList {
		eqName := strings.ToLower(strings.TrimSpace(eq.Name))
		if eqName == name {
			return &eq
		}
	}
	return nil
}
