package main

import (
	"DnD-sheet/internal/character/infrastructure"
	"DnD-sheet/internal/character/service"
	"DnD-sheet/internal/cli"
	"fmt"
	"os"
	"strings"
)

const dataDir = "../data"

func main() {
	// Initialize dependencies using the new refactored architecture
	characterRepo := infrastructure.NewJSONCharacterRepository(dataDir)
	characterService := service.NewCharacterService(characterRepo)

	// Create CLI instance
	cliApp := cli.NewCLI()

	// Register commands
	cliApp.Register(cli.NewCreateCommand(characterService))
	cliApp.Register(cli.NewViewCommand(characterService))
	cliApp.Register(cli.NewListCommand(characterService))
	cliApp.Register(cli.NewDeleteCommand(characterService))
	cliApp.Register(cli.NewUpdateCommand(characterService))
	cliApp.Register(cli.NewEquipCommand(characterService))
	cliApp.Register(cli.NewPrepareSpellCommand(characterService))
	cliApp.Register(cli.NewLearnSpellCommand(characterService))

	// Run CLI
	if err := cliApp.Run(os.Args); err != nil {
		fmt.Printf("%v\n", err)
		// Exit with code 2 for validation errors (missing required fields)
		if err.Error() == "name is required" ||
			err.Error() == "name and level (>=1) are required" {
			os.Exit(2)
		}
		// Exit with code 1 for "not found" errors and equipment conflicts
		if strings.Contains(err.Error(), "not found") ||
			strings.Contains(err.Error(), "already occupied") {
			os.Exit(1)
		}
		os.Exit(1)
	}
}
