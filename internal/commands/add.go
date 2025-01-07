package commands

import (
	"flag"
	"task-cli/internal/task"
)

type AddCommand struct {
	tm        task.ITaskManager
	presenter Presenter
}

// NewAddCommand creates a new instance of AddCommand
func NewAddCommand(tm task.ITaskManager, p Presenter) *AddCommand {
	return &AddCommand{
		tm:        tm,
		presenter: p,
	}
}

// Execute executes the add command
func (c *AddCommand) Execute(args []string) error {
	cmd := flag.NewFlagSet("add", flag.ExitOnError)
	title := cmd.String("title", "", "Task title")
	priorityFlag := cmd.String("priority", task.DefaultPriority.String(), "Task priority (low, medium, high)")
	dueDate := cmd.String("due", "", "Due date (format: YYYY-MM-DD HH:MM)")
	reminder := cmd.String("reminder", "", "Reminder time (format: YYYY-MM-DD HH:MM)")

	if err := cmd.Parse(args); err != nil {
		return c.presenter.PrintError("error parsing arguments: %v", err)
	}

	if *title == "" {
		return c.presenter.PrintError("task title is required")
	}

	priority, err := task.ParsePriority(*priorityFlag)
	if err != nil {
		return c.presenter.PrintError("invalid priority: %v", err)
	}

	// Create the task
	newTask := c.tm.AddTask(*title, priority)

	// Set due date if provided
	if *dueDate != "" {
		due, err := task.ParseDateTime(*dueDate)
		if err != nil {
			return c.presenter.PrintError("invalid due date: %v", err)
		}
		if err := c.tm.SetDueDate(newTask.ID, due); err != nil {
			return c.presenter.PrintError("error setting due date: %v", err)
		}
	}

	// Set reminder if provided
	if *reminder != "" {
		rem, err := task.ParseDateTime(*reminder)
		if err != nil {
			return c.presenter.PrintError("invalid reminder time: %v", err)
		}
		if err := c.tm.SetReminder(newTask.ID, rem); err != nil {
			return c.presenter.PrintError("error setting reminder: %v", err)
		}
	}

	if err := c.tm.SaveTasks(); err != nil {
		return c.presenter.PrintError("error saving task: %v", err)
	}

	c.presenter.PrintSuccess("Task added with ID: %d (Priority: %s%s%s)",
		newTask.ID,
		newTask.Priority.Color(),
		newTask.Priority.String(),
		"\033[0m")

	return nil
}

// Help returns information about the add command
func (c *AddCommand) Help() string {
	return `Create a new task
    
Usage:
  task add [flags]

Flags:
  -title string      Task title (required)
  -priority string   Task priority: low, medium, high (default: medium)
  -due string        Due date (format: YYYY-MM-DD HH:MM)
  -reminder string   Reminder time (format: YYYY-MM-DD HH:MM)`
}
