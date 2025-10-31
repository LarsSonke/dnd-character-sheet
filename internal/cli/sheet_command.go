package cli

import (
	"DnD-sheet/internal/character/service"
	"fmt"
)

// SheetCommand handles character sheet export
type SheetCommand struct {
	*BaseCommand
	characterService *service.CharacterService

	// Flags
	name   *string
	format *string
}

// NewSheetCommand creates a new sheet command
func NewSheetCommand(characterService *service.CharacterService) *SheetCommand {
	cmd := &SheetCommand{
		BaseCommand:      NewBaseCommand("sheet"),
		characterService: characterService,
	}

	// Define flags
	cmd.name = cmd.flagSet.String("name", "", "name of the character")
	cmd.format = cmd.flagSet.String("format", "markdown", "output format (markdown)")

	return cmd
}

// Name returns the command name
func (c *SheetCommand) Name() string {
	return "sheet"
}

// Execute exports the character sheet
func (c *SheetCommand) Execute() error {
	if *c.name == "" {
		return fmt.Errorf("character name is required")
	}

	// Load the character
	char, err := c.characterService.GetCharacter(*c.name)
	if err != nil {
		return fmt.Errorf("failed to load character: %w", err)
	}

	// Check format
	if *c.format != "markdown" {
		return fmt.Errorf("only markdown format is currently supported")
	}

	// Create markdown formatter and export
	formatter := service.NewMarkdownFormatter()
	markdownOutput := formatter.FormatCharacter(char)

	// Print to stdout
	fmt.Print(markdownOutput)

	return nil
}

// Usage prints usage information for the sheet command
func (c *SheetCommand) Usage() {
	fmt.Println("  sheet -name CHARACTER_NAME [-format FORMAT] - export character sheet")
}
