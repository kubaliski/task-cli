package main

import (
	"flag"
	"fmt"
	"os"
	"task-cli/internal/task"
)

func main() {
	tm := task.NewTaskManager()
	if err := tm.LoadTasks(); err != nil {
		fmt.Fprintf(os.Stderr, "error loading tasks: %v\n", err)
		os.Exit(1)
	}

	//subcommands
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)

	//Add flags
	addTitle := addCmd.String("title", "", "task title")

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "subcommand is required")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		addCmd.Parse(os.Args[2:])
		if *addTitle == "" {
			fmt.Fprintln(os.Stderr, "title is required")
			os.Exit(1)
		}
		task := tm.AddTask(*addTitle)
		if err := tm.SaveTasks(); err != nil {
			fmt.Fprintf(os.Stderr, "error saving tasks: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("task added with ID %d\n", task.ID)

	case "list":
		listCmd.Parse(os.Args[2:])
		tasks := tm.GetTasks()
		if len(tasks) == 0 {
			fmt.Println("No tasks")
			return
		}
		for _, task := range tasks {
			status := " "
			if task.Done {
				status = "✓"
			}
			fmt.Printf("%d %s %s\n", task.ID, status, task.Title)
		}
	default:
		fmt.Printf("unknown subcommand: %s\n", os.Args[1])
		os.Exit(1)
	}
}
