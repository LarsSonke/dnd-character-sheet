package service

import (
	"DnD-sheet/helpers"
	"DnD-sheet/internal/character/domain"
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

	// Apply racial bonuses
	req = s.applyRacialBonuses(req)

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

	// Assign spell slots based on class and level
	c.SpellSlots = s.generateSpellSlots(req.Class, req.Level)

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
	c.SpellSlots = s.generateSpellSlots(c.Class, newLevel)

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
	if !s.canCastSpells(c.Class) {
		return errors.New("this class can't cast spells")
	}

	// Check if this is a prepared caster (they can't learn spells, only prepare them)
	if s.isPreparedCaster(c.Class) {
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
func (s *CharacterService) PrepareSpell(name, spell string) error {
	c, err := s.repo.Load(name)
	if err != nil {
		return err
	}

	// Check if class can cast spells
	if !s.canCastSpells(c.Class) {
		return errors.New("this class can't cast spells")
	}

	// Check if this is a known caster (they can't prepare spells, only learn them)
	if !s.isPreparedCaster(c.Class) {
		return errors.New("this class learns spells and can't prepare them")
	}

	// Check if character has spell slots for this spell level
	spellLevel := s.getSpellLevel(spell)
	if spellLevel > 0 { // Only check for leveled spells (not cantrips)
		if slots, hasSlots := c.SpellSlots[spellLevel]; !hasSlots || slots == 0 {
			return errors.New("the spell has higher level than the available spell slots")
		}
	}

	// Check if spell is already prepared
	for _, preparedSpell := range c.PreparedSpells {
		if strings.EqualFold(preparedSpell, spell) {
			return errors.New("spell already prepared")
		}
	}

	c.PreparedSpells = append(c.PreparedSpells, spell)
	return s.repo.Save(c)
}

// applyRacialBonuses applies racial ability score bonuses
func (s *CharacterService) applyRacialBonuses(req CreateCharacterRequest) CreateCharacterRequest {
	switch strings.ToLower(req.Race) {
	case "dwarf":
		req.Con += 2
	case "elf":
		req.Dex += 2
	case "halfling":
		req.Dex += 2
	case "lightfoot halfling":
		req.Dex += 2
		req.Cha += 1
	case "stout halfling":
		req.Dex += 2
		req.Con += 1
	case "human":
		req.Str += 1
		req.Dex += 1
		req.Con += 1
		req.Int += 1
		req.Wis += 1
		req.Cha += 1
	case "dragonborn":
		req.Str += 2
		req.Cha += 1
	case "gnome":
		req.Int += 2
	case "half elf":
		req.Cha += 2
	case "half orc":
		req.Str += 2
		req.Con += 1
	case "tiefling":
		req.Int += 1
		req.Cha += 2
	case "hill dwarf":
		req.Con += 2
		req.Wis += 1
	}
	return req
}

// generateSkillProficiencies creates skill list based on background and class
func (s *CharacterService) generateSkillProficiencies(background, class string) []string {
	bgSkills := helpers.BackgroundSkillProficiencies[strings.ToLower(background)]
	classSkills := helpers.ClassSkillProficiencies[strings.ToLower(class)]

	// Start with background skills
	skillList := []string{}
	skillList = append(skillList, bgSkills...)

	// Add up to N class skills, even if they duplicate background skills
	nClassSkills := helpers.ClassSkillCount[strings.ToLower(class)]
	count := 0
	for _, s := range classSkills {
		skillList = append(skillList, s)
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

// generateSpellSlots creates spell slots based on class and level
func (s *CharacterService) generateSpellSlots(class string, level int) map[int]int {
	switch strings.ToLower(class) {
	case "wizard", "cleric", "druid", "bard", "sorcerer":
		slots := domain.FullCasterSpellSlots(level)
		// Add cantrips (Level 0) for full casters
		slots[0] = domain.FullCasterCantrips(level)
		return slots
	case "paladin", "ranger":
		return domain.HalfCasterSpellSlots(level)
	case "warlock":
		return domain.PactMagicSpellSlots(level)
	default:
		return map[int]int{}
	}
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

// canCastSpells returns true if the given class can cast spells
func (s *CharacterService) canCastSpells(class string) bool {
	spellcastingClasses := []string{
		"wizard", "sorcerer", "warlock", "cleric", "druid", "paladin", "ranger",
		"bard", "artificer", "eldritch knight", "arcane trickster",
	}

	classLower := strings.ToLower(class)
	for _, spellClass := range spellcastingClasses {
		if classLower == spellClass {
			return true
		}
	}
	return false
}

// isPreparedCaster returns true if the class prepares spells (vs learning them)
func (s *CharacterService) isPreparedCaster(class string) bool {
	// Known caster classes that can learn spells
	knownClasses := []string{
		"sorcerer", "warlock", "bard", "eldritch knight", "arcane trickster",
	}

	classLower := strings.ToLower(class)
	// If it's a known caster, it's NOT a prepared caster
	for _, knownClass := range knownClasses {
		if classLower == knownClass {
			return false
		}
	}

	// All other spellcasters are prepared casters (including wizard)
	return s.canCastSpells(class)
}

// getSpellLevel returns the spell level for common D&D spells
func (s *CharacterService) getSpellLevel(spell string) int {
	spellLevels := map[string]int{
		// Cantrips (Level 0)
		"fire bolt": 0, "mage hand": 0, "prestidigitation": 0, "light": 0,

		// Level 1 spells
		"burning hands": 1, "magic missile": 1, "cure wounds": 1, "shield": 1, "false life": 1, "feather fall": 1,

		// Level 2 spells
		"scorching ray": 2, "misty step": 2, "web": 2, "hold person": 2,

		// Level 3 spells
		"fireball": 3, "lightning bolt": 3, "counterspell": 3, "fly": 3,

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

	spellLower := strings.ToLower(spell)
	if level, exists := spellLevels[spellLower]; exists {
		return level
	}

	// Default to level 1 for unknown spells to be safe
	return 1
}
