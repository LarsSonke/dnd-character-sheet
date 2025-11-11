package service

import (
	"DnD-sheet/internal/character/domain"
	"DnD-sheet/internal/spell"
	"errors"
	"fmt"
	"sort"
	"strings"
)

// CharacterService handles character business logic
type CharacterService struct {
	repo domain.CharacterRepository
}

// NewCharacterService creates a new character service
func NewCharacterService(repo domain.CharacterRepository) *CharacterService {
	return &CharacterService{repo: repo}
}

// GetRepository returns the character repository (for web server access)
func (s *CharacterService) GetRepository() domain.CharacterRepository {
	return s.repo
}

// validSkills maps skill names for validation
var validSkills = map[string]bool{
	"acrobatics": true, "animal handling": true, "arcana": true, "athletics": true,
	"deception": true, "history": true, "insight": true, "intimidation": true,
	"investigation": true, "medicine": true, "nature": true, "perception": true,
	"performance": true, "persuasion": true, "religion": true, "sleight of hand": true,
	"stealth": true, "survival": true,
}

// CreateCharacterRequest contains parameters for creating a character
type CreateCharacterRequest struct {
	Name       string
	Race       string
	Class      string
	Level      int
	Str        int
	Dex        int
	Con        int
	Int        int
	Wis        int
	Cha        int
	Background string
}

// CreateCharacter creates a new character with racial bonuses and skill proficiencies
func (s *CharacterService) CreateCharacter(req CreateCharacterRequest) (*domain.Character, error) {
	// Validate input
	if req.Name == "" {
		return nil, errors.New("character name is required")
	}

	// Check if character already exists
	if s.repo.Exists(req.Name) {
		return nil, errors.New("character already exists")
	}

	// Set default background if not provided
	if req.Background == "" {
		req.Background = "acolyte"
	}

	// Apply racial bonuses using domain logic
	race := domain.NewRace(req.Race)
	bonuses := race.GetAbilityBonuses()
	req.Str += bonuses["str"]
	req.Dex += bonuses["dex"]
	req.Con += bonuses["con"]
	req.Int += bonuses["int"]
	req.Wis += bonuses["wis"]
	req.Cha += bonuses["cha"]

	// Automatically assign skill proficiencies based on D&D rules
	skills := s.generateSkillProficiencies(req.Background, req.Class)

	// Validate skills
	if err := s.validateSkills(skills); err != nil {
		return nil, err
	}

	// Create character
	c := domain.NewCharacter(
		req.Name, req.Race, req.Class, req.Level,
		req.Str, req.Dex, req.Con, req.Int, req.Wis, req.Cha,
		req.Background, skills,
	)

	// Assign spell slots based on class and level using domain logic
	c.SpellSlots = c.GetSpellSlots()

	// Save character
	if err := s.repo.Save(c); err != nil {
		return nil, err
	}

	return c, nil
}

// GetCharacter retrieves a character by name
func (s *CharacterService) GetCharacter(name string) (*domain.Character, error) {
	return s.repo.Load(name)
}

// ListCharacters returns all character names
func (s *CharacterService) ListCharacters() ([]string, error) {
	return s.repo.List()
}

// DeleteCharacter removes a character
func (s *CharacterService) DeleteCharacter(name string) error {
	return s.repo.Delete(name)
}

// UpdateLevel updates a character's level and recalculates dependent stats
func (s *CharacterService) UpdateLevel(name string, newLevel int) error {
	c, err := s.repo.Load(name)
	if err != nil {
		return err
	}

	oldLevel := c.Level
	c.Level = newLevel
	c.ProficiencyBonus = domain.ProficiencyBonus(newLevel)
	c.ApplySRDAbilityScoreImprovements(oldLevel, newLevel)
	c.SpellSlots = c.GetSpellSlots()

	return s.repo.Save(c)
}

// EquipCharacter equips a character with weapons, armor, and shields
func (s *CharacterService) EquipCharacter(name, weapon, armor, shield, weaponSlot string) error {
	c, err := s.repo.Load(name)
	if err != nil {
		return err
	}

	if weapon != "" {
		// Check if character already has a weapon equipped
		if c.Weapon != "" {
			// Determine the target slot
			targetSlot := weaponSlot
			if targetSlot == "" {
				targetSlot = "main hand" // Default slot
			}

			// Check if the current weapon is in the same slot we're trying to equip to
			currentSlot := c.WeaponSlot
			if currentSlot == "" {
				currentSlot = "main hand" // Default slot for existing weapon
			}

			if currentSlot == targetSlot {
				return fmt.Errorf("%s already occupied", targetSlot)
			}
		}

		c.Weapon = weapon
		if weaponSlot != "" {
			c.WeaponSlot = weaponSlot
		} else {
			c.WeaponSlot = "main hand" // Default slot
		}
	}
	if armor != "" {
		c.Armor = armor
	}
	if shield != "" {
		c.Shield = shield
	}

	return s.repo.Save(c)
}

// LearnSpell adds a spell to a character's known spells
func (s *CharacterService) LearnSpell(name, spell string) error {
	c, err := s.repo.Load(name)
	if err != nil {
		return err
	}

	// Check if class can cast spells
	if !c.IsSpellcaster() {
		return errors.New("this class can't cast spells")
	}

	// Check if this is a prepared caster (they can't learn spells, only prepare them)
	if c.IsPreparedCaster() {
		return errors.New("this class prepares spells and can't learn them")
	}

	// Check if spell is already known
	for _, knownSpell := range c.KnownSpells {
		if strings.EqualFold(knownSpell, spell) {
			return errors.New("spell already known")
		}
	}

	c.KnownSpells = append(c.KnownSpells, spell)
	return s.repo.Save(c)
}

// PrepareSpell adds a spell to a character's prepared spells
func (s *CharacterService) PrepareSpell(name, spellName string) error {
	c, err := s.repo.Load(name)
	if err != nil {
		return err
	}

	// Check if class can cast spells
	if !c.IsSpellcaster() {
		return errors.New("this class can't cast spells")
	}

	// Check if this is a known caster (they can't prepare spells, only learn them)
	if !c.IsPreparedCaster() {
		return errors.New("this class learns spells and can't prepare them")
	}

	// Check if character has spell slots for this spell level
	spellLevel := spell.GetSpellLevel(spellName)
	if spellLevel > 0 { // Only check for leveled spells (not cantrips)
		if slots, hasSlots := c.SpellSlots[spellLevel]; !hasSlots || slots == 0 {
			return errors.New("the spell has higher level than the available spell slots")
		}
	}

	// Check if spell is already prepared
	for _, preparedSpell := range c.PreparedSpells {
		if strings.EqualFold(preparedSpell, spellName) {
			return errors.New("spell already prepared")
		}
	}

	c.PreparedSpells = append(c.PreparedSpells, spellName)
	return s.repo.Save(c)
}

// generateSkillProficiencies creates skill list based on background and class using domain logic
func (s *CharacterService) generateSkillProficiencies(background, class string) []string {
	bg := domain.NewBackground(background)
	bgSkills := bg.GetSkillProficiencies()

	cl := domain.NewClass(class)
	classSkills := cl.GetAvailableSkills()
	nClassSkills := cl.GetSkillCount()

	// Start with background skills
	skillList := []string{}
	skillList = append(skillList, bgSkills...)

	// Add up to N class skills, even if they duplicate background skills
	count := 0
	for _, skill := range classSkills {
		skillList = append(skillList, skill)
		count++
		if count >= nClassSkills {
			break
		}
	}

	sort.Strings(skillList)
	return skillList
}

// validateSkills checks if all skills are valid
func (s *CharacterService) validateSkills(skills []string) error {
	for _, skill := range skills {
		if !validSkills[strings.ToLower(skill)] {
			return errors.New("invalid skill proficiency: " + skill)
		}
	}
	return nil
}

// IsStandardArray validates if ability scores use the standard array
func (s *CharacterService) IsStandardArray(str, dex, con, int_, wis, cha int) bool {
	standard := []int{16, 14, 13, 12, 10, 8}
	input := []int{str, dex, con, int_, wis, cha}
	used := make([]bool, 6)
	for _, val := range input {
		found := false
		for i, std := range standard {
			if val == std && !used[i] {
				used[i] = true
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
