package commands

import (
	"flag"
	"strings"
)

type HelpCommand struct {
	commands  map[string]Command
	presenter Presenter
}

// NewHelpCommand creates a new instance of HelpCommand
func NewHelpCommand(commands map[string]Command, p Presenter) *HelpCommand {
	return &HelpCommand{
		commands:  commands,
		presenter: p,
	}
}

// Execute executes the help command
func (c *HelpCommand) Execute(args []string) error {
	cmd := flag.NewFlagSet("help", flag.ExitOnError)
	if err := cmd.Parse(args); err != nil {
		return c.presenter.PrintError("error parsing arguments: %v", err)
	}

	// If a command is provided, show help for that command
	if len(cmd.Args()) > 0 {
		return c.showCommandHelp(cmd.Args()[0])
	}

	// Show general help
	return c.showGeneralHelp()
}

// ShowCommandHelp shows help for a specific command
func (c *HelpCommand) showCommandHelp(commandName string) error {
	cmd, exists := c.commands[commandName]
	if !exists {
		return c.presenter.PrintError("unknown command: %s", commandName)
	}

	help := cmd.Help()
	c.presenter.PrintSuccess(help)
	return nil
}

func (c *HelpCommand) showGeneralHelp() error {
	var sb strings.Builder

	sb.WriteString("Task CLI - A simple task manager\n\n")
	sb.WriteString("Usage:\n")
	sb.WriteString("  task <command> [flags]\n\n")
	sb.WriteString("Available Commands:\n")

	// Lista de comandos ordenada
	commands := []struct {
		name        string
		description string
	}{
		{"add", "Create a new task"},
		{"list", "List and filter tasks"},
		{"update", "Update an existing task"},
		{"delete", "Remove a task"},
		{"get", "Show detailed task information"},
		{"help", "Show help about any command"},
	}

	// Find the maximum width of the command names
	maxWidth := 0
	for _, cmd := range commands {
		if len(cmd.name) > maxWidth {
			maxWidth = len(cmd.name)
		}
	}

	// Print each command with padding
	for _, cmd := range commands {
		sb.WriteString(c.formatCommandHelp(cmd.name, cmd.description, maxWidth))
	}

	sb.WriteString("\nUse 'task help <command>' for more information about a command\n")

	c.presenter.PrintSuccess(sb.String())
	return nil
}

// formatCommandHelp formats the help message for a command
func (c *HelpCommand) formatCommandHelp(name, description string, width int) string {
	padding := width - len(name)
	if padding < 0 {
		padding = 0
	}
	return "  " + name + strings.Repeat(" ", padding+2) + description + "\n"
}

// Help returns the help message for the help command
func (c *HelpCommand) Help() string {
	return `Show help about commands

Usage:
  task help [command]

Arguments:
  [command]    Optional command name to get detailed help for`
}
