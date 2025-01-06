package commands

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
	"task-cli/internal/task"
)

type Commander struct {
	tm *task.TaskManager
}

func NewCommander(tm *task.TaskManager) *Commander {
	return &Commander{tm: tm}
}

func (c *Commander) Add(args []string) error {
	cmd := flag.NewFlagSet("add", flag.ExitOnError)
	title := cmd.String("title", "", "Task title")
	priorityFlag := cmd.String("priority", task.DefaultPriority.String(), "Task priority (low, medium, high)")

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

	newTask := c.tm.AddTask(*title, priority)
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

func (c *Commander) List(args []string) error {
	cmd := flag.NewFlagSet("list", flag.ExitOnError)
	byPriority := cmd.Bool("priority", false, "Sort tasks by priority")

	if err := cmd.Parse(args); err != nil {
		return err
	}

	tasks := c.tm.GetTasksSorted(*byPriority)
	if len(tasks) == 0 {
		fmt.Println("No tasks found")
		return nil
	}

	// Define column widths
	idWidth := 4
	statusWidth := 10
	priorityWidth := 10
	titleWidth := 40
	createdWidth := 20
	completedWidth := 20

	// Print table header
	totalWidth := idWidth + statusWidth + priorityWidth + titleWidth + createdWidth + completedWidth + 7
	fmt.Println(strings.Repeat("-", totalWidth))
	fmt.Printf("| %-*s | %-*s | %-*s | %-*s | %-*s | %-*s |\n",
		idWidth-2, "ID",
		statusWidth-2, "Status",
		priorityWidth-2, "Priority",
		titleWidth-2, "Title",
		createdWidth-2, "Created At",
		completedWidth-2, "Completed At")
	fmt.Println(strings.Repeat("-", totalWidth))

	for _, t := range tasks {
		status := "Pending"
		completedAt := "---"
		if t.Done {
			status = "✓ Done"
			completedAt = t.CompletedAt.Format("2006-01-02 15:04")
		}

		title := t.Title
		if len(title) > titleWidth-2 {
			title = title[:titleWidth-5] + "..."
		}

		fmt.Printf("| %-*d | %-*s | %s%-*s%s | %-*s | %-*s | %-*s |\n",
			idWidth-2, t.ID,
			statusWidth-2, status,
			t.Priority.Color(),
			priorityWidth-2, t.Priority.String(),
			"\033[0m",
			titleWidth-2, title,
			createdWidth-2, t.CreatedAt.Format("2006-01-02 15:04"),
			completedWidth-2, completedAt)
	}

	fmt.Println(strings.Repeat("-", totalWidth))
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

	task, err := c.tm.GetTaskByID(id)
	if err != nil {
		return err
	}

	status := " "
	if task.Done {
		status = "✓"
	}

	fmt.Printf("Task #%d\n", task.ID)
	fmt.Printf("Status: [%s]\n", status)
	fmt.Printf("Priority: %s%s%s\n", task.Priority.Color(), task.Priority.String(), "\033[0m")
	fmt.Printf("Title: %s\n", task.Title)
	fmt.Printf("Created: %s\n", task.CreatedAt.Format("2006-01-02 15:04:05"))
	if task.Done {
		fmt.Printf("Completed: %s\n", task.CompletedAt.Format("2006-01-02 15:04:05"))
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

func (c *Commander) Update(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("task ID is required")
	}

	cmd := flag.NewFlagSet("update", flag.ExitOnError)
	title := cmd.String("title", "", "New task title")
	done := cmd.Bool("done", false, "Mark task as done")
	priorityFlag := cmd.String("priority", "", "Task priority (low, medium, high)")

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

	newTitle := currentTask.Title
	if *title != "" {
		newTitle = *title
	}

	var priority *task.TaskPriority
	if *priorityFlag != "" {
		p, err := task.ParsePriority(*priorityFlag)
		if err != nil {
			return err
		}
		priority = &p
	}

	newDone := *done || currentTask.Done

	if err := c.tm.UpdateTask(id, newTitle, newDone, priority); err != nil {
		return fmt.Errorf("error updating task: %v", err)
	}

	if err := c.tm.SaveTasks(); err != nil {
		return err
	}

	fmt.Printf("Task %d updated successfully\n", id)
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
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}
