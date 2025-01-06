# Task CLI

A powerful and intuitive command-line task manager built with Go, designed to help you organize and track your tasks efficiently.

## Features

### Task Management

- Create new tasks with titles
- View tasks in a beautiful tabular format
- Update task titles and completion status
- Delete tasks when no longer needed
- Mark tasks as completed
- Get detailed information about specific tasks

### Data Handling

- Automatic data persistence using JSON
- Safe and efficient data storage
- Data stored in user's home directory

## Installation

### Prerequisites

- Go 1.16 or higher
- Git

### Installation Steps

```bash
# Clone the repository
git clone https://github.com/yourusername/task-cli.git
cd task-cli

# Build the application
# For Windows
go build -o task.exe
# For Linux/Mac
go build -o task

# Add to PATH (optional)
# Move the binary to a directory in your PATH
```

## Usage

### Basic Commands

```bash
# Add a new task
task add -title "Complete project documentation"

# List all tasks
task list

# Get detailed information about a specific task
task get <id>

# Update a task
task update <id> -title "New title"     # Update title
task update <id> -done                  # Mark as completed
task update <id> -title "New title" -done   # Update both

# Delete a task
task delete <id>
```

### Command Details

| Command  | Description                                         | Example                            |
| -------- | --------------------------------------------------- | ---------------------------------- |
| `add`    | Creates a new task                                  | `task add -title "Buy groceries"`  |
| `list`   | Shows all tasks in a table format                   | `task list`                        |
| `get`    | Displays detailed information about a specific task | `task get 1`                       |
| `update` | Modifies an existing task                           | `task update 1 -title "New title"` |
| `delete` | Removes a task                                      | `task delete 1`                    |

## Roadmap

### Task Priorities

- High, medium, and low priority levels
- Priority-based sorting

### Time Management

- Due dates
- Reminders
- Task scheduling

### Organization

- Categories and tags
- Multiple task lists
- Nested tasks

### Enhanced Features

- Advanced search and filtering
- Task statistics and analytics
- Rich task descriptions
- Export to various formats (CSV, PDF)
- Task archiving

## Contributing

Contributions are welcome! Feel free to:

1. Fork the repository
2. Create a new branch
3. Make your changes
4. Submit a pull request

## Contact

If you have any questions or suggestions, feel free to:

- Open an issue
- Submit a pull request
- Contact the maintainers
