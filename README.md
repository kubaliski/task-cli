# Task CLI

A simple command-line task manager application built with Go.

## Current Features

- Add new tasks
- List existing tasks
- Data persistence using JSON
- Task completion status

## Upcoming Features

- Mark tasks as done
- Delete tasks
- Task priorities (high, medium, low)
- Due dates
- Categories/tags
- Task filtering and search
- Basic statistics (completed/pending tasks)
- Task descriptions
- Export tasks to different formats (CSV, PDF)

## Usage

```bash
# Add a new task
./task add -title "Task name"

# List all tasks
./task list
```

```

## Installation

1. Clone the repository
2. Run `go build -o task ./cmd/main.go`
3. Use the `task` executable

## Requirements

- Go 1.16 or higher

```
