package commands

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
	"task-cli/internal/task"
	"time"
)

type Commander struct {
	tm task.ITaskManager
}

func NewCommander(tm task.ITaskManager) *Commander {
	return &Commander{tm: tm}
}

func (c *Commander) Add(args []string) error {
	cmd := flag.NewFlagSet("add", flag.ExitOnError)
	title := cmd.String("title", "", "Task title")
	priorityFlag := cmd.String("priority", task.DefaultPriority.String(), "Task priority (low, medium, high)")
	dueDate := cmd.String("due", "", "Due date (format: YYYY-MM-DD HH:MM)")
	reminder := cmd.String("reminder", "", "Reminder time (format: YYYY-MM-DD HH:MM)")

	if err := cmd.Parse(args); err != nil {
		return err
	}

	if *title == "" {
		return fmt.Errorf("task title is required")
	}

	priority, err := task.ParsePriority(*priorityFlag)
	if err != nil {
		return err
	}

	// Create the task
	newTask := c.tm.AddTask(*title, priority)

	// Set due date if provided
	if *dueDate != "" {
		due, err := task.ParseDateTime(*dueDate)
		if err != nil {
			return fmt.Errorf("invalid due date: %v", err)
		}
		if err := c.tm.SetDueDate(newTask.ID, due); err != nil {
			return err
		}
	}

	// Set reminder if provided
	if *reminder != "" {
		rem, err := task.ParseDateTime(*reminder)
		if err != nil {
			return fmt.Errorf("invalid reminder time: %v", err)
		}
		if err := c.tm.SetReminder(newTask.ID, rem); err != nil {
			return err
		}
	}

	if err := c.tm.SaveTasks(); err != nil {
		return err
	}

	fmt.Printf("Task added with ID: %d (Priority: %s%s%s)\n",
		newTask.ID,
		newTask.Priority.Color(),
		newTask.Priority.String(),
		"\033[0m")
	return nil
}

// timeFilter represents the different time filter options available
type timeFilter struct {
	value    string
	dueDate  *time.Time
	status   *task.TimeStatus
	upcoming bool
}

// parseTimeFilter processes the specified time filter string
func parseTimeFilter(filter string) (*timeFilter, error) {
	if filter == "" {
		return nil, nil
	}

	tf := &timeFilter{value: filter}
	now := time.Now()

	switch strings.ToLower(filter) {
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
		// Try to parse as specific date
		if date, err := task.ParseDateTime(filter); err == nil {
			tf.dueDate = &date
		} else {
			return nil, fmt.Errorf("invalid time filter: %s", filter)
		}
	}

	return tf, nil
}

func (c *Commander) List(args []string) error {
	cmd := flag.NewFlagSet("list", flag.ExitOnError)

	// Sorting flags
	byPriority := cmd.Bool("priority", false, "Sort tasks by priority")
	byDueDate := cmd.Bool("by-due", false, "Sort tasks by due date")

	// Time filtering flags
	dueFilter := cmd.String("due", "", `Filter tasks by due date. Options:
		- "today": Due today
		- "tomorrow": Due tomorrow
		- "thisweek": Due this week
		- "nextweek": Due next week
		- "overdue": Overdue tasks
		- "duesoon": Tasks due soon
		- "upcoming": Tasks with upcoming reminders
		- Or specify a date: "2024-01-20 15:00"`)

	// Other flags
	showCompleted := cmd.Bool("all", false, "Show completed tasks")
	format := cmd.String("format", "table", "Output format: table or list")

	if err := cmd.Parse(args); err != nil {
		return err
	}

	// Process time filter
	tf, err := parseTimeFilter(*dueFilter)
	if err != nil {
		return err
	}

	// Get and filter tasks
	tasks := c.tm.GetTasksSorted(*byPriority, *byDueDate)
	var filteredTasks []task.Task

	for _, t := range tasks {
		// Filter completed tasks
		if t.Done && !*showCompleted {
			continue
		}

		// Apply time filters if they exist
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

		filteredTasks = append(filteredTasks, t)
	}

	if len(filteredTasks) == 0 {
		fmt.Println("No tasks found matching the criteria")
		return nil
	}

	// Print tasks in the specified format
	if *format == "list" {
		return c.printTasksList(filteredTasks)
	}
	return c.printTasksTable(filteredTasks)
}
func (c *Commander) Update(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("task ID is required")
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
		return err
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid task ID: %v", err)
	}

	currentTask, err := c.tm.GetTaskByID(id)
	if err != nil {
		return fmt.Errorf("error getting task: %v", err)
	}

	// Update title if provided
	newTitle := currentTask.Title
	if *title != "" {
		newTitle = *title
	}

	// Update priority if provided
	var priority *task.TaskPriority
	if *priorityFlag != "" {
		p, err := task.ParsePriority(*priorityFlag)
		if err != nil {
			return err
		}
		priority = &p
	}

	// Update completion status
	newDone := *done || currentTask.Done

	// Update the basic task properties
	if err := c.tm.UpdateTask(id, newTitle, newDone, priority); err != nil {
		return fmt.Errorf("error updating task: %v", err)
	}

	// Handle due date updates
	if *removeDue {
		if err := c.tm.RemoveDueDate(id); err != nil {
			return err
		}
	} else if *dueDate != "" {
		due, err := task.ParseDateTime(*dueDate)
		if err != nil {
			return fmt.Errorf("invalid due date: %v", err)
		}
		if err := c.tm.SetDueDate(id, due); err != nil {
			return err
		}
	}

	// Handle reminder updates
	if *removeReminder {
		if err := c.tm.RemoveReminder(id); err != nil {
			return err
		}
	} else if *reminder != "" {
		rem, err := task.ParseDateTime(*reminder)
		if err != nil {
			return fmt.Errorf("invalid reminder time: %v", err)
		}
		if err := c.tm.SetReminder(id, rem); err != nil {
			return err
		}
	}

	if err := c.tm.SaveTasks(); err != nil {
		return err
	}

	fmt.Printf("Task %d updated successfully\n", id)
	return nil
}

func (c *Commander) Get(args []string) error {
	cmd := flag.NewFlagSet("get", flag.ExitOnError)
	if err := cmd.Parse(args); err != nil {
		return err
	}

	if len(cmd.Args()) == 0 {
		return fmt.Errorf("task ID is required")
	}

	id, err := strconv.Atoi(cmd.Args()[0])
	if err != nil {
		return fmt.Errorf("invalid task ID: %v", err)
	}

	t, err := c.tm.GetTaskByID(id)
	if err != nil {
		return err
	}

	fmt.Printf("Task #%d\n", t.ID)
	fmt.Printf("Title: %s\n", t.Title)
	fmt.Printf("Status: %s\n", getStatusString(t))
	fmt.Printf("Priority: %s%s%s\n", t.Priority.Color(), t.Priority.String(), "\033[0m")
	fmt.Printf("Created: %s\n", t.CreatedAt.Format("2006-01-02 15:04:05"))

	if t.DueDate != nil {
		dueDate := task.FormatDateTime(t.DueDate)
		if t.IsOverdue() {
			dueDate = task.TimeStatusOverdue.Color() + dueDate + "\033[0m"
		}
		fmt.Printf("Due Date: %s\n", dueDate)
	}

	if t.Reminder != nil {
		reminder := task.FormatDateTime(t.Reminder)
		if t.IsUpcoming() {
			reminder = task.TimeStatusUpcoming.Color() + reminder + "\033[0m"
		}
		fmt.Printf("Reminder: %s\n", reminder)
	}

	if t.Done {
		fmt.Printf("Completed: %s\n", t.CompletedAt.Format("2006-01-02 15:04:05"))
	}

	return nil
}

func (c *Commander) Delete(args []string) error {
	cmd := flag.NewFlagSet("delete", flag.ExitOnError)
	if err := cmd.Parse(args); err != nil {
		return err
	}

	if len(cmd.Args()) == 0 {
		return fmt.Errorf("task ID is required")
	}

	id, err := strconv.Atoi(cmd.Args()[0])
	if err != nil {
		return fmt.Errorf("invalid task ID: %v", err)
	}

	if err := c.tm.DeleteTask(id); err != nil {
		return err
	}

	if err := c.tm.SaveTasks(); err != nil {
		return err
	}

	fmt.Printf("Task %d deleted successfully\n", id)
	return nil
}

func (c *Commander) Execute(command string, args []string) error {
	switch command {
	case "add":
		return c.Add(args)
	case "list":
		return c.List(args)
	case "delete":
		return c.Delete(args)
	case "update":
		return c.Update(args)
	case "get":
		return c.Get(args)
	case "help":
		return c.Help(args)
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}
func (c *Commander) printTasksTable(tasks []task.Task) error {
	// Define column widths
	idWidth := 4
	statusWidth := 12
	priorityWidth := 10
	titleWidth := 30
	dueDateWidth := 18
	reminderWidth := 18
	createdWidth := 18

	// Print table header
	totalWidth := idWidth + statusWidth + priorityWidth + titleWidth + dueDateWidth + reminderWidth + createdWidth + 15
	fmt.Println(strings.Repeat("-", totalWidth))
	fmt.Printf("| %-*s | %-*s | %-*s | %-*s | %-*s | %-*s | %-*s |\n",
		idWidth-2, "ID",
		statusWidth-2, "Status",
		priorityWidth-2, "Priority",
		titleWidth-2, "Title",
		dueDateWidth-2, "Due Date",
		reminderWidth-2, "Reminder",
		createdWidth-2, "Created")
	fmt.Println(strings.Repeat("-", totalWidth))

	for _, t := range tasks {
		status := getStatusString(t)
		dueDate := task.FormatDateTime(t.DueDate)
		if t.IsOverdue() {
			dueDate = task.TimeStatusOverdue.Color() + dueDate + "\033[0m"
		}

		reminder := task.FormatDateTime(t.Reminder)
		if t.IsUpcoming() {
			reminder = task.TimeStatusUpcoming.Color() + reminder + "\033[0m"
		}

		title := t.Title
		if len(title) > titleWidth-2 {
			title = title[:titleWidth-5] + "..."
		}

		fmt.Printf("| %-*d | %-*s | %s%-*s%s | %-*s | %-*s | %-*s | %-*s |\n",
			idWidth-2, t.ID,
			statusWidth-2, status,
			t.Priority.Color(), priorityWidth-2, t.Priority.String(), "\033[0m",
			titleWidth-2, title,
			dueDateWidth-2, dueDate,
			reminderWidth-2, reminder,
			createdWidth-2, t.CreatedAt.Format("2006-01-02 15:04"))
	}

	fmt.Println(strings.Repeat("-", totalWidth))
	return nil
}

func (c *Commander) printTasksList(tasks []task.Task) error {
	for _, t := range tasks {
		// Print task header with ID and title
		fmt.Printf("\n%s Task #%d: %s%s%s - %s\n",
			getStatusIcon(t),
			t.ID,
			t.Priority.Color(),
			t.Priority.String(),
			"\033[0m",
			t.Title)

		// Print dates and times
		fmt.Printf("   Created: %s\n", t.CreatedAt.Format("2006-01-02 15:04"))

		if t.DueDate != nil {
			dueStr := fmt.Sprintf("   Due: %s", task.FormatDateTime(t.DueDate))
			if t.IsOverdue() {
				dueStr = task.TimeStatusOverdue.Color() + dueStr + "\033[0m"
			}
			fmt.Println(dueStr)
		}

		if t.Reminder != nil {
			reminderStr := fmt.Sprintf("   Reminder: %s", task.FormatDateTime(t.Reminder))
			if t.IsUpcoming() {
				reminderStr = task.TimeStatusUpcoming.Color() + reminderStr + "\033[0m"
			}
			fmt.Println(reminderStr)
		}

		if t.Done {
			fmt.Printf("   Completed: %s\n", t.CompletedAt.Format("2006-01-02 15:04"))
		}

		// Print separator between tasks
		fmt.Println(strings.Repeat("-", 50))
	}
	return nil
}

func getStatusString(t task.Task) string {
	if t.Done {
		return "✓ Done"
	}
	switch t.GetTimeStatus() {
	case task.TimeStatusOverdue:
		return "! Overdue"
	case task.TimeStatusDueSoon:
		return "→ DueSoon"
	case task.TimeStatusUpcoming:
		return "⏰ Upcoming"
	default:
		return "Pending"
	}
}

func getStatusIcon(t task.Task) string {
	if t.Done {
		return "[✓]"
	}
	switch t.GetTimeStatus() {
	case task.TimeStatusOverdue:
		return "[!]"
	case task.TimeStatusDueSoon:
		return "[→]"
	case task.TimeStatusUpcoming:
		return "[⏰]"
	default:
		return "[ ]"
	}
}
func (c *Commander) Help(args []string) error {
	cmd := flag.NewFlagSet("help", flag.ExitOnError)
	if err := cmd.Parse(args); err != nil {
		return err
	}

	// Si se proporciona un comando específico, mostrar ayuda detallada para ese comando
	if len(cmd.Args()) > 0 {
		return c.showCommandHelp(cmd.Args()[0])
	}

	// Mostrar ayuda general
	fmt.Println("Task CLI - A simple task manager")
	fmt.Println("\nUsage:")
	fmt.Println("  task <command> [flags]")
	fmt.Println("\nCommands:")
	fmt.Println("  add         Create a new task")
	fmt.Println("  list        List and filter tasks")
	fmt.Println("  update      Update an existing task")
	fmt.Println("  delete      Remove a task")
	fmt.Println("  get         Show detailed task information")
	fmt.Println("  help        Show help about any command")
	fmt.Println("\nRun 'task help <command>' for detailed usage of each command")
	return nil
}

func (c *Commander) showCommandHelp(command string) error {
	switch command {
	case "add":
		fmt.Println("Create a new task")
		fmt.Println("\nUsage:")
		fmt.Println("  task add [flags]")
		fmt.Println("\nFlags:")
		fmt.Println("  -title string      Task title (required)")
		fmt.Println("  -priority string   Task priority: low, medium, high (default: medium)")
		fmt.Println("  -due string        Due date (format: YYYY-MM-DD HH:MM)")
		fmt.Println("  -reminder string   Reminder time (format: YYYY-MM-DD HH:MM)")

	case "list":
		fmt.Println("List and filter tasks")
		fmt.Println("\nUsage:")
		fmt.Println("  task list [flags]")
		fmt.Println("\nFlags:")
		fmt.Println("  -priority          Sort by priority")
		fmt.Println("  -by-due           Sort by due date")
		fmt.Println("  -due string       Filter by time: today, tomorrow, thisweek, nextweek")
		fmt.Println("                    overdue, duesoon, upcoming")
		fmt.Println("  -all              Show completed tasks")
		fmt.Println("  -format string    Output format: table, list (default: table)")

	case "update":
		fmt.Println("Update an existing task")
		fmt.Println("\nUsage:")
		fmt.Println("  task update <id> [flags]")
		fmt.Println("\nFlags:")
		fmt.Println("  -title string      New task title")
		fmt.Println("  -priority string   Change priority: low, medium, high")
		fmt.Println("  -done              Mark as completed")
		fmt.Println("  -due string        Set due date (YYYY-MM-DD HH:MM)")
		fmt.Println("  -reminder string   Set reminder (YYYY-MM-DD HH:MM)")
		fmt.Println("  -remove-due        Remove due date")
		fmt.Println("  -remove-reminder   Remove reminder")

	case "delete":
		fmt.Println("Remove a task")
		fmt.Println("\nUsage:")
		fmt.Println("  task delete <id>")

	case "get":
		fmt.Println("Show detailed task information")
		fmt.Println("\nUsage:")
		fmt.Println("  task get <id>")

	case "help":
		fmt.Println("Show help about any command")
		fmt.Println("\nUsage:")
		fmt.Println("  task help [command]")

	default:
		return fmt.Errorf("unknown command: %s", command)
	}
	return nil
}
