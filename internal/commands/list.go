package commands

import (
	"flag"
	"task-cli/internal/task"
	"time"
)

type ListCommand struct {
	tm        task.ITaskManager
	presenter Presenter
}

// timeFilter represents a filter for tasks based on time criteria
type timeFilter struct {
	value    string
	dueDate  *time.Time
	status   *task.TimeStatus
	upcoming bool
}

// NewListCommand creates a new instance of ListCommand
func NewListCommand(tm task.ITaskManager, p Presenter) *ListCommand {
	return &ListCommand{
		tm:        tm,
		presenter: p,
	}
}

// Execute executes the list command
func (c *ListCommand) Execute(args []string) error {
	cmd := flag.NewFlagSet("list", flag.ExitOnError)

	// Order flags
	byPriority := cmd.Bool("priority", false, "Sort tasks by priority")
	byDueDate := cmd.Bool("by-due", false, "Sort by due date")

	// Due time flags
	dueFilter := cmd.String("due", "",
		`Filter tasks by due date. Options:
         - "today":     Due today
         - "tomorrow":  Due tomorrow
         - "thisweek":  Due this week
         - "nextweek":  Due next week
         - "overdue":   Overdue tasks
         - "duesoon":   Tasks due soon
         - "upcoming":  Tasks with upcoming reminders
         - Or specify a date: "2024-01-20 15:00"`)

	// Other flags
	showCompleted := cmd.Bool("all", false, "Show completed tasks")
	format := cmd.String("format", "table", "Output format: table or list")

	if err := cmd.Parse(args); err != nil {
		return c.presenter.PrintError("error parsing arguments: %v", err)
	}

	// Process time filter
	tf, err := c.parseTimeFilter(*dueFilter)
	if err != nil {
		return c.presenter.PrintError("invalid time filter: %v", err)
	}

	// Get and filter tasks based on flags
	tasks := c.tm.GetTasksSorted(*byPriority, *byDueDate)
	filteredTasks := c.filterTasks(tasks, tf, *showCompleted)

	if len(filteredTasks) == 0 {
		c.presenter.PrintSuccess("No tasks found matching the criteria")
		return nil
	}

	// Show tasks in the selected format
	if *format == "list" {
		return c.presenter.PrintTaskList(filteredTasks)
	}
	return c.presenter.PrintTaskTable(filteredTasks)
}

// filterTasks filters tasks based on the provided criteria
func (c *ListCommand) filterTasks(tasks []task.Task, tf *timeFilter, showCompleted bool) []task.Task {
	var filtered []task.Task

	for _, t := range tasks {
		// Filter completed tasks
		if t.Done && !showCompleted {
			continue
		}

		// Aply time filter
		if tf != nil {
			if tf.status != nil && t.GetTimeStatus() != *tf.status {
				continue
			}
			if tf.dueDate != nil {
				if t.DueDate == nil || t.DueDate.After(*tf.dueDate) {
					continue
				}
			}
			if tf.upcoming && !t.IsUpcoming() {
				continue
			}
		}

		filtered = append(filtered, t)
	}

	return filtered
}

// parseTimeFilter parses the time filter flag
func (c *ListCommand) parseTimeFilter(filter string) (*timeFilter, error) {
	if filter == "" {
		return nil, nil
	}

	tf := &timeFilter{value: filter}
	now := time.Now()

	switch filter {
	case "today":
		endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
		tf.dueDate = &endOfDay
	case "tomorrow":
		tomorrow := now.AddDate(0, 0, 1)
		endOfTomorrow := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 23, 59, 59, 0, now.Location())
		tf.dueDate = &endOfTomorrow
	case "thisweek":
		daysUntilEndOfWeek := 7 - int(now.Weekday())
		if daysUntilEndOfWeek == 0 {
			daysUntilEndOfWeek = 7
		}
		endOfWeek := time.Date(now.Year(), now.Month(), now.Day()+daysUntilEndOfWeek, 23, 59, 59, 0, now.Location())
		tf.dueDate = &endOfWeek
	case "nextweek":
		startOfNextWeek := now.AddDate(0, 0, 7-int(now.Weekday())+1)
		endOfNextWeek := startOfNextWeek.AddDate(0, 0, 6)
		endOfDay := time.Date(endOfNextWeek.Year(), endOfNextWeek.Month(), endOfNextWeek.Day(), 23, 59, 59, 0, now.Location())
		tf.dueDate = &endOfDay
	case "overdue":
		status := task.TimeStatusOverdue
		tf.status = &status
	case "upcoming":
		tf.upcoming = true
	case "duesoon":
		status := task.TimeStatusDueSoon
		tf.status = &status
	default:
		// Try to parse a specific date
		if date, err := task.ParseDateTime(filter); err == nil {
			tf.dueDate = &date
		} else {
			return nil, err
		}
	}

	return tf, nil
}

// Help returns the help message for the list command
func (c *ListCommand) Help() string {
	return `List and filter tasks

Usage:
  task list [flags]

Flags:
  -priority          Sort by priority
  -by-due           Sort by due date
  -due string       Filter by time: today, tomorrow, thisweek, nextweek,
                    overdue, duesoon, upcoming, or specify date (YYYY-MM-DD HH:MM)
  -all              Show completed tasks
  -format string    Output format: table or list (default: table)`
}
