package main

import (
	"fmt"
	"os"
	"task-cli/internal/commands"
	"task-cli/internal/task"
)

func main() {
	tm := task.NewTaskManager()
	if err := tm.LoadTasks(); err != nil {
		fmt.Fprintf(os.Stderr, "Error loading tasks: %v\n", err)
		os.Exit(1)
	}

	commander := commands.NewCommander(tm)

	// If no command is provided, show help
	if len(os.Args) < 2 {
		if err := commander.Help([]string{}); err != nil {
			fmt.Fprintf(os.Stderr, "Error showing help: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	command := os.Args[1]
	args := os.Args[2:]

	if err := commander.Execute(command, args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
