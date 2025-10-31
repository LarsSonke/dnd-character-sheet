package infrastructure

import (
	"DnD-sheet/internal/equipment/domain"
	"encoding/csv"
	"errors"
	"os"
	"strconv"
	"strings"
)

// CSVEquipmentRepository implements equipment persistence using CSV files
type CSVEquipmentRepository struct {
	csvPath   string
	equipment []domain.Equipment
	loaded    bool
}

// NewCSVEquipmentRepository creates a new CSV equipment repository
func NewCSVEquipmentRepository(csvPath string) *CSVEquipmentRepository {
	return &CSVEquipmentRepository{csvPath: csvPath}
}

// LoadAll loads all equipment from the CSV file
func (r *CSVEquipmentRepository) LoadAll() ([]domain.Equipment, error) {
	if err := r.ensureLoaded(); err != nil {
		return nil, err
	}
	return r.equipment, nil
}

// FindByName finds equipment by name (case-insensitive)
func (r *CSVEquipmentRepository) FindByName(name string) (*domain.Equipment, error) {
	if err := r.ensureLoaded(); err != nil {
		return nil, err
	}

	name = strings.ToLower(strings.TrimSpace(name))
	for _, eq := range r.equipment {
		eqName := strings.ToLower(strings.TrimSpace(eq.Name))
		if eqName == name {
			return &eq, nil
		}
	}
	return nil, errors.New("equipment not found: " + name)
}

// FindByCategory returns all equipment in a specific category
func (r *CSVEquipmentRepository) FindByCategory(category string) ([]domain.Equipment, error) {
	if err := r.ensureLoaded(); err != nil {
		return nil, err
	}

	var result []domain.Equipment
	category = strings.ToLower(strings.TrimSpace(category))
	for _, eq := range r.equipment {
		eqCategory := strings.ToLower(strings.TrimSpace(eq.Category))
		if eqCategory == category {
			result = append(result, eq)
		}
	}
	return result, nil
}

// ensureLoaded loads the CSV data if not already loaded
func (r *CSVEquipmentRepository) ensureLoaded() error {
	if r.loaded {
		return nil
	}

	file, err := os.Open(r.csvPath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	var equipmentList []domain.Equipment
	for i, rec := range records {
		if i == 0 {
			continue // skip header
		}
		if len(rec) < 2 {
			continue // skip incomplete rows
		}

		ac := domain.ArmorClass{Base: 10} // default armor class
		if len(rec) > 3 {
			if base, err := strconv.Atoi(strings.TrimSpace(rec[2])); err == nil {
				ac.Base = base
			}
			if dexBonus := strings.TrimSpace(rec[3]); strings.EqualFold(dexBonus, "yes") {
				ac.DexBonus = true
			}
		}

		equipmentList = append(equipmentList, domain.Equipment{
			Name:       strings.TrimSpace(rec[0]),
			Category:   strings.TrimSpace(rec[1]),
			ArmorClass: ac,
		})
	}

	r.equipment = equipmentList
	r.loaded = true
	return nil
}
