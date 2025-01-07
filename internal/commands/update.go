package commands

import (
	"flag"
	"strconv"
	"task-cli/internal/task"
)

type UpdateCommand struct {
	tm        task.ITaskManager
	presenter Presenter
}

// NewUpdateCommand creates a new instance of UpdateCommand
func NewUpdateCommand(tm task.ITaskManager, p Presenter) *UpdateCommand {
	return &UpdateCommand{
		tm:        tm,
		presenter: p,
	}
}

// Execute executes the update command
func (c *UpdateCommand) Execute(args []string) error {
	if len(args) < 1 {
		return c.presenter.PrintError("task ID is required")
	}

	// Parse task ID before parsing flags
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return c.presenter.PrintError("invalid task ID: %v", err)
	}

	cmd := flag.NewFlagSet("update", flag.ExitOnError)
	title := cmd.String("title", "", "New task title")
	done := cmd.Bool("done", false, "Mark task as done")
	priorityFlag := cmd.String("priority", "", "Task priority (low, medium, high)")
	dueDate := cmd.String("due", "", "Due date (format: YYYY-MM-DD HH:MM)")
	reminder := cmd.String("reminder", "", "Reminder time (format: YYYY-MM-DD HH:MM)")
	removeDue := cmd.Bool("remove-due", false, "Remove due date")
	removeReminder := cmd.Bool("remove-reminder", false, "Remove reminder")

	if err := cmd.Parse(args[1:]); err != nil {
		return c.presenter.PrintError("error parsing arguments: %v", err)
	}

	// Get current task
	currentTask, err := c.tm.GetTaskByID(id)
	if err != nil {
		return c.presenter.PrintError("error getting task: %v", err)
	}

	// Update basic fields
	if err := c.updateBasicFields(id, currentTask, *title, *done, *priorityFlag); err != nil {
		return err
	}

	// Manage due date updates
	if err := c.updateDueDate(id, *dueDate, *removeDue); err != nil {
		return err
	}

	// Manage reminder updates
	if err := c.updateReminder(id, *reminder, *removeReminder); err != nil {
		return err
	}

	// Save Changes
	if err := c.tm.SaveTasks(); err != nil {
		return c.presenter.PrintError("error saving changes: %v", err)
	}

	c.presenter.PrintSuccess("Task %d updated successfully", id)
	return nil
}

// updateBasicFields updates the basic fields of a task
func (c *UpdateCommand) updateBasicFields(id int, currentTask task.Task, newTitle string, done bool, priorityFlag string) error {
	// Update title if provided
	title := currentTask.Title
	if newTitle != "" {
		title = newTitle
	}

	// Update Priority if provided
	var priority *task.TaskPriority
	if priorityFlag != "" {
		p, err := task.ParsePriority(priorityFlag)
		if err != nil {
			return c.presenter.PrintError("invalid priority: %v", err)
		}
		priority = &p
	}

	// Update done status
	newDone := done || currentTask.Done

	// Update basic fields
	if err := c.tm.UpdateTask(id, title, newDone, priority); err != nil {
		return c.presenter.PrintError("error updating task: %v", err)
	}

	return nil
}

// updateDueDate updates the due date of a task
func (c *UpdateCommand) updateDueDate(id int, dueDate string, removeDue bool) error {
	if removeDue {
		if err := c.tm.RemoveDueDate(id); err != nil {
			return c.presenter.PrintError("error removing due date: %v", err)
		}
	} else if dueDate != "" {
		due, err := task.ParseDateTime(dueDate)
		if err != nil {
			return c.presenter.PrintError("invalid due date: %v", err)
		}
		if err := c.tm.SetDueDate(id, due); err != nil {
			return c.presenter.PrintError("error setting due date: %v", err)
		}
	}
	return nil
}

// updateReminder updates the reminder of a task
func (c *UpdateCommand) updateReminder(id int, reminder string, removeReminder bool) error {
	if removeReminder {
		if err := c.tm.RemoveReminder(id); err != nil {
			return c.presenter.PrintError("error removing reminder: %v", err)
		}
	} else if reminder != "" {
		rem, err := task.ParseDateTime(reminder)
		if err != nil {
			return c.presenter.PrintError("invalid reminder time: %v", err)
		}
		if err := c.tm.SetReminder(id, rem); err != nil {
			return c.presenter.PrintError("error setting reminder: %v", err)
		}
	}
	return nil
}

// Help returns the help message for the update command
func (c *UpdateCommand) Help() string {
	return `Update an existing task

Usage:
  task update <id> [flags]

Flags:
  -title string      New task title
  -priority string   Change priority: low, medium, high
  -done              Mark as completed
  -due string        Set due date (YYYY-MM-DD HH:MM)
  -reminder string   Set reminder (YYYY-MM-DD HH:MM)
  -remove-due        Remove due date
  -remove-reminder   Remove reminder`
}
