package commands

import (
	"flag"
	"fmt"
	"strconv"
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

	if err := cmd.Parse(args); err != nil {
		return err
	}

	if *title == "" {
		return fmt.Errorf("task title is required")
	}

	task := c.tm.AddTask(*title)
	if err := c.tm.SaveTasks(); err != nil {
		return err
	}

	fmt.Printf("Task added with ID: %d\n", task.ID)
	return nil
}

func (c *Commander) List(_ []string) error {
	tasks := c.tm.GetTasks()
	if len(tasks) == 0 {
		fmt.Println("No tasks found")
		return nil
	}

	for _, t := range tasks {
		status := " "
		if t.Done {
			status = "✓"
		}
		fmt.Printf("[%s] %d: %s\n", status, t.ID, t.Title)
	}
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

	// Mostrar la información detallada de la tarea
	status := " "
	if task.Done {
		status = "✓"
	}

	fmt.Printf("Task #%d\n", task.ID)
	fmt.Printf("Status: [%s]\n", status)
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
	cmd := flag.NewFlagSet("update", flag.ExitOnError)
	title := cmd.String("title", "", "New task title")
	done := cmd.Bool("done", false, "Mark task as done")

	// Necesitamos procesar los flags antes de procesar el ID
	if err := cmd.Parse(args[1:]); err != nil {
		return err
	}

	// El ID debería ser el primer argumento
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid task ID: %v", err)
	}

	currentTask, err := c.tm.GetTaskByID(id)
	if err != nil {
		return fmt.Errorf("error getting task: %v", err)
	}

	// Usamos el nuevo título solo si se proporcionó uno
	newTitle := currentTask.Title
	if *title != "" {
		newTitle = *title
	}

	// Usamos el valor del flag done
	newDone := *done || currentTask.Done

	if err := c.tm.UpdateTask(id, newTitle, newDone); err != nil {
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
