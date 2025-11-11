package service

import (
	"DnD-sheet/internal/equipment/domain"
	"errors"
)

// EquipmentService handles equipment business logic
type EquipmentService struct {
	repo domain.EquipmentRepository
}

// NewEquipmentService creates a new equipment service
func NewEquipmentService(repo domain.EquipmentRepository) *EquipmentService {
	return &EquipmentService{repo: repo}
}

// GetAllEquipment returns all available equipment
func (s *EquipmentService) GetAllEquipment() ([]domain.Equipment, error) {
	return s.repo.LoadAll()
}

// FindEquipmentByName finds equipment by name (case-insensitive)
func (s *EquipmentService) FindEquipmentByName(name string) (*domain.Equipment, error) {
	if name == "" {
		return nil, errors.New("equipment name cannot be empty")
	}
	return s.repo.FindByName(name)
}

// GetEquipmentByCategory returns all equipment in a specific category
func (s *EquipmentService) GetEquipmentByCategory(category string) ([]domain.Equipment, error) {
	if category == "" {
		return nil, errors.New("category cannot be empty")
	}
	return s.repo.FindByCategory(category)
}

// CalculateArmorClass calculates AC for a character based on equipped armor and dexterity
func (s *EquipmentService) CalculateArmorClass(armorName string, dexModifier int) (int, error) {
	if armorName == "" {
		// No armor: 10 + Dex modifier
		return 10 + dexModifier, nil
	}

	armor, err := s.repo.FindByName(armorName)
	if err != nil {
		return 10, err
	}

	// Use domain logic to calculate AC
	return armor.CalculateAC(dexModifier), nil
}
