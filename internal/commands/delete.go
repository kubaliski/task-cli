package commands

import (
	"flag"
	"strconv"
	"task-cli/internal/task"
)

type DeleteCommand struct {
	tm        task.ITaskManager
	presenter Presenter
}

// NewDeleteCommand creates a new instance of DeleteCommand
func NewDeleteCommand(tm task.ITaskManager, p Presenter) *DeleteCommand {
	return &DeleteCommand{
		tm:        tm,
		presenter: p,
	}
}

// Execute executes the delete command
func (c *DeleteCommand) Execute(args []string) error {
	cmd := flag.NewFlagSet("delete", flag.ExitOnError)
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

	// Verify if the task exists
	if _, err := c.tm.GetTaskByID(id); err != nil {
		return c.presenter.PrintError("task not found: %v", err)
	}

	if err := c.tm.DeleteTask(id); err != nil {
		return c.presenter.PrintError("error deleting task: %v", err)
	}

	if err := c.tm.SaveTasks(); err != nil {
		return c.presenter.PrintError("error saving changes: %v", err)
	}

	c.presenter.PrintSuccess("Task %d deleted successfully", id)
	return nil
}

// Help returns the help message for the delete command
func (c *DeleteCommand) Help() string {
	return `Remove a task

Usage:
  task delete <id>

Arguments:
  <id>    The ID of the task to delete`
}
