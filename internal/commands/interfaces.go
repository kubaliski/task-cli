package commands

import "task-cli/internal/task"

// Command defind the base interface for all commands
type Command interface {
	// Execute execute the command with the given arguments
	Execute(args []string) error
	// Help returns the help message for the command
	Help() string
}

// Commands is a map of commands to their names
type Commands map[string]Command

// Presenter definds the interface for presenting information to the user
type Presenter interface {
	// PrintTaskTable shows the tasks in a table format
	PrintTaskTable(tasks []task.Task) error
	// PrintTaskList shows the tasks in a list format
	PrintTaskList(tasks []task.Task) error
	// PrintTask shows an individual task in a detailed format
	PrintTask(t task.Task) error
	// PrintSuccess shows a success message
	PrintSuccess(format string, a ...interface{})
	// PrintError shows an error message
	PrintError(format string, a ...interface{}) error
}

// TaskFilter defines the interface for filtering tasks
type TaskFilter interface {
	// Apply filters the tasks based on the filter criteria
	Apply(tasks []task.Task) []task.Task
}
