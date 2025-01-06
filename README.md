# Task CLI

A simple command-line task manager application built with Go.

## Current Features

- Add new tasks
- List existing tasks
- Update tasks (title and status)
- Delete tasks
- Data persistence using JSON
- Task completion status

## Usage

```bash
# Add a new task
task add -title "Task name"

# List all tasks
task list

# Get a specific task
task get <id>

# Update a task's title
task update <id> -title "New title"

# Mark a task as done
task update <id> -done

# Update both title and status
task update <id> -title "New title" -done

# Delete a task
task delete <id>
Installation

Clone the repository
Run go build -o task.exe ./cmd/ (use task instead of task.exe on Linux/Mac)
Add the directory containing the executable to your PATH
Use the task command from anywhere

Requirements

Go 1.16 or higher

Upcoming Features

Task priorities (high, medium, low)
Due dates
Categories/tags
Task filtering and search
Basic statistics (completed/pending tasks)
Task descriptions
Export tasks to different formats (CSV, PDF)

```
