package cli

import (
	"errors"
	"flag"
)

// Command represents a CLI command interface
type Command interface {
	Name() string
	Parse(args []string) error
	Execute() error
	Usage()
}

// CLI manages all CLI commands
type CLI struct {
	commands map[string]Command
}

// NewCLI creates a new CLI instance
func NewCLI() *CLI {
	return &CLI{
		commands: make(map[string]Command),
	}
}

// Register adds a command to the CLI
func (c *CLI) Register(cmd Command) {
	c.commands[cmd.Name()] = cmd
}

// Run executes the CLI with provided arguments
func (c *CLI) Run(args []string) error {
	if len(args) < 2 {
		c.usage()
		return errors.New("no command provided")
	}

	cmdName := args[1]
	cmd, exists := c.commands[cmdName]
	if !exists {
		c.usage()
		return errors.New("unknown command: " + cmdName)
	}

	if err := cmd.Parse(args[2:]); err != nil {
		return err
	}

	return cmd.Execute()
}

// usage prints general CLI usage information
func (c *CLI) usage() {
	// Print usage information for all commands
	println("Usage:")
	for _, cmd := range c.commands {
		cmd.Usage()
	}
}

// BaseCommand provides common functionality for all commands
type BaseCommand struct {
	flagSet *flag.FlagSet
}

// NewBaseCommand creates a new base command
func NewBaseCommand(name string) *BaseCommand {
	return &BaseCommand{
		flagSet: flag.NewFlagSet(name, flag.ExitOnError),
	}
}

// Parse parses command line arguments
func (bc *BaseCommand) Parse(args []string) error {
	return bc.flagSet.Parse(args)
}

// Usage prints command usage
func (bc *BaseCommand) Usage() {
	bc.flagSet.Usage()
}
