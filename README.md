# Task CLI

A powerful and intuitive command-line task manager built with Go, designed to help you organize and track your tasks efficiently.

## Features

### Task Management

- Create new tasks with titles and priority levels
- View tasks in a beautiful tabular format with color-coded priorities
- Sort tasks by ID or priority
- Update task titles, completion status, and priority levels
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
task add -title "Complete project documentation" -priority high

# List all tasks
task list                # Sort by ID (default)
task list -priority     # Sort by priority

# Get detailed information about a specific task
task get <id>

# Update a task
task update <id> -title "New title"                    # Update title
task update <id> -done                                # Mark as completed
task update <id> -priority high                       # Change priority
task update <id> -title "New title" -done -priority low  # Update multiple fields

# Delete a task
task delete <id>
```

### Command Details

| Command  | Flags                                                                                    | Description                                                                     | Example                                           |
| -------- | ---------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------- | ------------------------------------------------- |
| `add`    | `-title` (required)<br>`-priority` (optional, default: medium)                           | Creates a new task                                                              | `task add -title "Buy groceries" -priority high`  |
| `list`   | `-priority` (optional)                                                                   | Shows all tasks in a table format.<br>Use `-priority` to sort by priority level | `task list`<br>`task list -priority`              |
| `get`    | `<id>` (required)                                                                        | Displays detailed information about a specific task                             | `task get 1`                                      |
| `update` | `<id>` (required)<br>`-title` (optional)<br>`-done` (optional)<br>`-priority` (optional) | Modifies an existing task.<br>Multiple flags can be combined                    | `task update 1 -title "New title" -priority high` |
| `delete` | `<id>` (required)                                                                        | Removes a task                                                                  | `task delete 1`                                   |

### Priority Levels

Tasks can be assigned one of three priority levels:

- `high`: For urgent and important tasks
- `medium`: Default priority level
- `low`: For less urgent tasks

## Roadmap

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

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

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
