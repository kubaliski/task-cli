// commands/commander.go
package commands

import (
	"task-cli/internal/task"
)

// Commander cordinates all the commands
type Commander struct {
	tm        task.ITaskManager
	presenter Presenter
	commands  map[string]Command
}

// NewCommander create a new instance of Commander
func NewCommander(tm task.ITaskManager) *Commander {
	c := &Commander{
		tm:        tm,
		presenter: NewDefaultPresenter(),
		commands:  make(map[string]Command),
	}
	c.registerCommands()
	return c
}

// registerCommands registers all the commands
func (c *Commander) registerCommands() {
	c.commands = map[string]Command{
		"add":    NewAddCommand(c.tm, c.presenter),
		"list":   NewListCommand(c.tm, c.presenter),
		"update": NewUpdateCommand(c.tm, c.presenter),
		"delete": NewDeleteCommand(c.tm, c.presenter),
		"get":    NewGetCommand(c.tm, c.presenter),
	}
	// Help command needs the list of commands
	c.commands["help"] = NewHelpCommand(c.commands, c.presenter)
}

// Execute executes the command with the given name
func (c *Commander) Execute(cmdName string, args []string) error {
	cmd, exists := c.commands[cmdName]
	if !exists {
		return c.presenter.PrintError("unknown command: %s", cmdName)
	}
	return cmd.Execute(args)
}

// SetPresenter allows to set a custom presenter
func (c *Commander) SetPresenter(p Presenter) {
	c.presenter = p
}

// GetCommands returns the list of registered commands
func (c *Commander) GetCommands() map[string]Command {
	return c.commands
}
