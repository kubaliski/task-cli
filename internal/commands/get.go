package commands

import (
	"flag"
	"strconv"
	"task-cli/internal/task"
)

type GetCommand struct {
	tm        task.ITaskManager
	presenter Presenter
}

// NewGetCommand creates a new instance of GetCommand
func NewGetCommand(tm task.ITaskManager, p Presenter) *GetCommand {
	return &GetCommand{
		tm:        tm,
		presenter: p,
	}
}

// Execute executes the get command
func (c *GetCommand) Execute(args []string) error {
	cmd := flag.NewFlagSet("get", flag.ExitOnError)
	if err := cmd.Parse(args); err != nil {
		return c.presenter.PrintError("error parsing arguments: %v", err)
	}

	if len(cmd.Args()) == 0 {
		return c.presenter.PrintError("task ID is required")
	}

	id, err := strconv.Atoi(cmd.Args()[0])
	if err != nil {
		return c.presenter.PrintError("invalid task ID: %v", err)
	}

	t, err := c.tm.GetTaskByID(id)
	if err != nil {
		return c.presenter.PrintError("error getting task: %v", err)
	}

	return c.presenter.PrintTask(t)
}

// Help returns the help message for the get command
func (c *GetCommand) Help() string {
	return `Show detailed task information

Usage:
  task get <id>

Arguments:
  <id>    The ID of the task to display`
}
