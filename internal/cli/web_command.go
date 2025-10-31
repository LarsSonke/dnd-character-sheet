package cli

import (
	"DnD-sheet/internal/character/service"
	"DnD-sheet/internal/web"
	"fmt"
)

// WebCommand starts a web server for viewing character sheets
type WebCommand struct {
	*BaseCommand
	characterService *service.CharacterService

	// Flags
	port *int
}

// NewWebCommand creates a new WebCommand instance
func NewWebCommand(characterService *service.CharacterService) *WebCommand {
	cmd := &WebCommand{
		BaseCommand:      NewBaseCommand("web"),
		characterService: characterService,
	}

	// Define flags
	cmd.port = cmd.flagSet.Int("port", 8080, "port to listen on")

	return cmd
}

// Name returns the command name
func (c *WebCommand) Name() string {
	return "web"
}

// Execute starts the web server
func (c *WebCommand) Execute() error {
	// Create web server
	server := web.NewServer(c.characterService.GetRepository())

	// Load templates
	templateDir := "web/templates"
	if err := server.LoadTemplates(templateDir); err != nil {
		return fmt.Errorf("failed to load templates: %w", err)
	}

	fmt.Printf("Starting D&D Character Sheet web server...\n")
	fmt.Printf("Open your browser and go to http://localhost:%d\n", *c.port)
	fmt.Printf("Press Ctrl+C to stop the server\n")

	// Start server (this will block)
	return server.Start(*c.port)
}

// Usage prints usage information for the web command
func (c *WebCommand) Usage() {
	println("  web [-port PORT] - start web server for character sheets")
}
