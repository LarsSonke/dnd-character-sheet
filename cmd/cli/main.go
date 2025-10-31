package main

import (
	"DnD-sheet/internal/character/infrastructure"
	"DnD-sheet/internal/character/service"
	"DnD-sheet/internal/cli"
	"fmt"
	"os"
)

const dataDir = "../data"

func main() {
	// Initialize dependencies
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
	cliApp.Register(cli.NewSheetCommand(characterService))
	cliApp.Register(cli.NewWebCommand(characterService))
	cliApp.Register(cli.NewAPITestCommand())

	// Run CLI
	if err := cliApp.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
