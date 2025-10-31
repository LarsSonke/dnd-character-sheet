package cli

import (
	"DnD-sheet/internal/character/service"
	"fmt"
)

// CreateCommand handles character creation
type CreateCommand struct {
	*BaseCommand
	characterService *service.CharacterService

	// Flags
	name         *string
	race         *string
	class        *string
	level        *int
	str          *int
	dex          *int
	con          *int
	intelligence *int
	wis          *int
	cha          *int
	background   *string
}

// NewCreateCommand creates a new create command
func NewCreateCommand(characterService *service.CharacterService) *CreateCommand {
	cmd := &CreateCommand{
		BaseCommand:      NewBaseCommand("create"),
		characterService: characterService,
	}

	// Define flags
	cmd.name = cmd.flagSet.String("name", "", "character name (required)")
	cmd.race = cmd.flagSet.String("race", "", "race")
	cmd.class = cmd.flagSet.String("class", "", "class")
	cmd.level = cmd.flagSet.Int("level", 1, "level")
	cmd.str = cmd.flagSet.Int("str", 10, "strength")
	cmd.dex = cmd.flagSet.Int("dex", 10, "dexterity")
	cmd.con = cmd.flagSet.Int("con", 10, "constitution")
	cmd.intelligence = cmd.flagSet.Int("int", 10, "intelligence")
	cmd.wis = cmd.flagSet.Int("wis", 10, "wisdom")
	cmd.cha = cmd.flagSet.Int("cha", 10, "charisma")
	cmd.background = cmd.flagSet.String("background", "", "background")

	return cmd
}

// Name returns the command name
func (c *CreateCommand) Name() string {
	return "create"
}

// Execute creates a new character
func (c *CreateCommand) Execute() error {
	if *c.name == "" {
		return fmt.Errorf("name is required")
	}

	req := service.CreateCharacterRequest{
		Name:       *c.name,
		Race:       *c.race,
		Class:      *c.class,
		Level:      *c.level,
		Str:        *c.str,
		Dex:        *c.dex,
		Con:        *c.con,
		Int:        *c.intelligence,
		Wis:        *c.wis,
		Cha:        *c.cha,
		Background: *c.background,
	}

	character, err := c.characterService.CreateCharacter(req)
	if err != nil {
		return err
	}

	fmt.Printf("saved character %s\n", character.Name)
	return nil
}

// Usage prints create command usage
func (c *CreateCommand) Usage() {
	fmt.Println("  create -name CHARACTER_NAME -race RACE -class CLASS -level N -str N -dex N -con N -int N -wis N -cha N -background BACKGROUND")
}

// ViewCommand handles character viewing
type ViewCommand struct {
	*BaseCommand
	characterService *service.CharacterService

	name *string
}

// NewViewCommand creates a new view command
func NewViewCommand(characterService *service.CharacterService) *ViewCommand {
	cmd := &ViewCommand{
		BaseCommand:      NewBaseCommand("view"),
		characterService: characterService,
	}

	cmd.name = cmd.flagSet.String("name", "", "character name (required)")
	return cmd
}

// Name returns the command name
func (c *ViewCommand) Name() string {
	return "view"
}

// Execute views a character
func (c *ViewCommand) Execute() error {
	if *c.name == "" {
		return fmt.Errorf("name is required")
	}

	character, err := c.characterService.GetCharacter(*c.name)
	if err != nil {
		return fmt.Errorf("character \"%s\" not found", *c.name)
	}

	// Print character information
	c.printCharacterInfo(character)
	return nil
}

// Usage prints view command usage
func (c *ViewCommand) Usage() {
	fmt.Println("  view -name CHARACTER_NAME")
}

// ListCommand handles character listing
type ListCommand struct {
	*BaseCommand
	characterService *service.CharacterService
}

// NewListCommand creates a new list command
func NewListCommand(characterService *service.CharacterService) *ListCommand {
	return &ListCommand{
		BaseCommand:      NewBaseCommand("list"),
		characterService: characterService,
	}
}

// Name returns the command name
func (c *ListCommand) Name() string {
	return "list"
}

// Execute lists all characters
func (c *ListCommand) Execute() error {
	names, err := c.characterService.ListCharacters()
	if err != nil {
		return err
	}

	for _, name := range names {
		fmt.Println(name)
	}
	return nil
}

// Usage prints list command usage
func (c *ListCommand) Usage() {
	fmt.Println("  list")
}

// DeleteCommand handles character deletion
type DeleteCommand struct {
	*BaseCommand
	characterService *service.CharacterService

	name *string
}

// NewDeleteCommand creates a new delete command
func NewDeleteCommand(characterService *service.CharacterService) *DeleteCommand {
	cmd := &DeleteCommand{
		BaseCommand:      NewBaseCommand("delete"),
		characterService: characterService,
	}

	cmd.name = cmd.flagSet.String("name", "", "character name (required)")
	return cmd
}

// Name returns the command name
func (c *DeleteCommand) Name() string {
	return "delete"
}

// Execute deletes a character
func (c *DeleteCommand) Execute() error {
	if *c.name == "" {
		return fmt.Errorf("name is required")
	}

	if err := c.characterService.DeleteCharacter(*c.name); err != nil {
		return err
	}

	fmt.Printf("deleted %s\n", *c.name)
	return nil
}

// Usage prints delete command usage
func (c *DeleteCommand) Usage() {
	fmt.Println("  delete -name CHARACTER_NAME")
}

// UpdateCommand handles character level updates
type UpdateCommand struct {
	*BaseCommand
	characterService *service.CharacterService

	name  *string
	level *int
}

// NewUpdateCommand creates a new update command
func NewUpdateCommand(characterService *service.CharacterService) *UpdateCommand {
	cmd := &UpdateCommand{
		BaseCommand:      NewBaseCommand("update"),
		characterService: characterService,
	}

	cmd.name = cmd.flagSet.String("name", "", "character name (required)")
	cmd.level = cmd.flagSet.Int("level", 0, "new level")
	return cmd
}

// Name returns the command name
func (c *UpdateCommand) Name() string {
	return "update"
}

// Execute updates a character's level
func (c *UpdateCommand) Execute() error {
	if *c.name == "" || *c.level < 1 {
		return fmt.Errorf("name and level (>=1) are required")
	}

	if err := c.characterService.UpdateLevel(*c.name, *c.level); err != nil {
		return err
	}

	character, err := c.characterService.GetCharacter(*c.name)
	if err != nil {
		return err
	}

	fmt.Printf("Updated %s to level %d (Proficiency Bonus: %d)\n", character.Name, character.Level, character.ProficiencyBonus)
	return nil
}

// Usage prints update command usage
func (c *UpdateCommand) Usage() {
	fmt.Println("  update -name CHARACTER_NAME -level N")
}

// EquipCommand handles character equipment
type EquipCommand struct {
	*BaseCommand
	characterService *service.CharacterService

	name   *string
	weapon *string
	armor  *string
	shield *string
	slot   *string
}

// NewEquipCommand creates a new equip command
func NewEquipCommand(characterService *service.CharacterService) *EquipCommand {
	cmd := &EquipCommand{
		BaseCommand:      NewBaseCommand("equip"),
		characterService: characterService,
	}

	cmd.name = cmd.flagSet.String("name", "", "character name (required)")
	cmd.weapon = cmd.flagSet.String("weapon", "", "weapon name")
	cmd.armor = cmd.flagSet.String("armor", "", "armor name")
	cmd.shield = cmd.flagSet.String("shield", "", "shield name")
	cmd.slot = cmd.flagSet.String("slot", "", "equipment slot")
	return cmd
}

// Name returns the command name
func (c *EquipCommand) Name() string {
	return "equip"
}

// Execute equips a character with weapons, armor, or shields
func (c *EquipCommand) Execute() error {
	if *c.name == "" {
		return fmt.Errorf("name is required")
	}

	if err := c.characterService.EquipCharacter(*c.name, *c.weapon, *c.armor, *c.shield, *c.slot); err != nil {
		return err
	}

	// Print equipment messages like the original
	if *c.weapon != "" {
		slot := "main hand"
		if *c.slot != "" {
			slot = *c.slot
		}
		fmt.Printf("Equipped weapon %s to %s\n", *c.weapon, slot)
	}
	if *c.armor != "" {
		fmt.Printf("Equipped armor %s\n", *c.armor)
	}
	if *c.shield != "" {
		fmt.Printf("Equipped shield %s\n", *c.shield)
	}

	return nil
}

// Usage prints equip command usage
func (c *EquipCommand) Usage() {
	fmt.Println("  equip -name CHARACTER_NAME -weapon WEAPON_NAME -armor ARMOR_NAME -shield SHIELD_NAME")
}

// PrepareSpellCommand handles spell preparation
type PrepareSpellCommand struct {
	*BaseCommand
	characterService *service.CharacterService

	// Flags
	name  *string
	spell *string
}

// NewPrepareSpellCommand creates a new prepare-spell command
func NewPrepareSpellCommand(characterService *service.CharacterService) *PrepareSpellCommand {
	cmd := &PrepareSpellCommand{
		BaseCommand:      NewBaseCommand("prepare-spell"),
		characterService: characterService,
	}

	// Define flags
	cmd.name = cmd.flagSet.String("name", "", "character name (required)")
	cmd.spell = cmd.flagSet.String("spell", "", "spell name (required)")

	return cmd
}

// Name returns the command name
func (c *PrepareSpellCommand) Name() string {
	return "prepare-spell"
}

// Execute runs the prepare-spell command
func (c *PrepareSpellCommand) Execute() error {
	if *c.name == "" || *c.spell == "" {
		return fmt.Errorf("name and spell are required")
	}

	err := c.characterService.PrepareSpell(*c.name, *c.spell)
	if err != nil {
		return err
	}

	fmt.Printf("Prepared spell %s\n", *c.spell)
	return nil
}

// Usage prints prepare-spell command usage
func (c *PrepareSpellCommand) Usage() {
	fmt.Println("  prepare-spell -name CHARACTER_NAME -spell SPELL_NAME")
}

// LearnSpellCommand handles spell learning
type LearnSpellCommand struct {
	*BaseCommand
	characterService *service.CharacterService

	// Flags
	name  *string
	spell *string
}

// NewLearnSpellCommand creates a new learn-spell command
func NewLearnSpellCommand(characterService *service.CharacterService) *LearnSpellCommand {
	cmd := &LearnSpellCommand{
		BaseCommand:      NewBaseCommand("learn-spell"),
		characterService: characterService,
	}

	// Define flags
	cmd.name = cmd.flagSet.String("name", "", "character name (required)")
	cmd.spell = cmd.flagSet.String("spell", "", "spell name (required)")

	return cmd
}

// Name returns the command name
func (c *LearnSpellCommand) Name() string {
	return "learn-spell"
}

// Execute runs the learn-spell command
func (c *LearnSpellCommand) Execute() error {
	if *c.name == "" || *c.spell == "" {
		return fmt.Errorf("name and spell are required")
	}

	err := c.characterService.LearnSpell(*c.name, *c.spell)
	if err != nil {
		return err
	}

	fmt.Printf("Learned spell %s\n", *c.spell)
	return nil
}

// Usage prints learn-spell command usage
func (c *LearnSpellCommand) Usage() {
	fmt.Println("  learn-spell -name CHARACTER_NAME -spell SPELL_NAME")
}
