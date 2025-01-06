# Task CLI

A powerful and intuitive command-line task manager built with Go, designed to help you organize and track your tasks efficiently.

## Features

### Task Management

- Create new tasks with titles and priority levels
- Set due dates and reminders for tasks
- View tasks in a beautiful tabular format with color-coded priorities and statuses
- Sort tasks by ID, priority, or due date
- Filter tasks by time status (today, this week, overdue, etc.)
- Update task titles, completion status, priority levels, due dates, and reminders
- Delete tasks when no longer needed
- Mark tasks as completed
- Get detailed information about specific tasks

### Time Management

- Set due dates for tasks
- Add reminders before due dates
- Automatic status tracking:
  - Overdue tasks
  - Tasks due soon (within 24 hours)
  - Tasks with upcoming reminders
  - Normal tasks
- Time-based filtering options
- Visual indicators for task status

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
task add -title "Team meeting" -priority high -due "2024-01-10 15:00" -reminder "2024-01-10 14:00"

# List tasks
task list                  # Show pending tasks (default)
task list -all            # Show all tasks including completed
task list -priority       # Sort by priority
task list -by-due         # Sort by due date
task list -format list    # Show in detailed list format

# Filter tasks by time
task list -due today      # Tasks due today
task list -due tomorrow   # Tasks due tomorrow
task list -due thisweek   # Tasks due this week
task list -due nextweek   # Tasks due next week
task list -due overdue    # Show overdue tasks
task list -due duesoon    # Show tasks due soon
task list -due upcoming   # Show tasks with upcoming reminders

# Get detailed information about a specific task
task get <id>

# Update a task
task update <id> -title "New title"                    # Update title
task update <id> -done                                # Mark as completed
task update <id> -priority high                       # Change priority
task update <id> -due "2024-01-10 15:00"             # Set due date
task update <id> -reminder "2024-01-10 14:00"        # Set reminder
task update <id> -remove-due                         # Remove due date
task update <id> -remove-reminder                    # Remove reminder

# Delete a task
task delete <id>
```

### Command Details

| Command  | Flags                                                                                                                                          | Description                                                                 | Example                                                            |
| -------- | ---------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------- | ------------------------------------------------------------------ |
| `add`    | `-title` (required)<br>`-priority` (optional, default: medium)<br>`-due` (optional)<br>`-reminder` (optional)                                  | Creates a new task                                                          | `task add -title "Meeting" -priority high -due "2024-01-10 15:00"` |
| `list`   | `-priority` (sort by priority)<br>`-by-due` (sort by due date)<br>`-due` (filter by time)<br>`-all` (show completed)<br>`-format` (table/list) | Shows tasks in table/list format with various sorting and filtering options | `task list -due today -priority`                                   |
| `get`    | `<id>` (required)                                                                                                                              | Displays detailed information about a specific task                         | `task get 1`                                                       |
| `update` | `<id>` (required)<br>`-title`<br>`-done`<br>`-priority`<br>`-due`<br>`-reminder`<br>`-remove-due`<br>`-remove-reminder`                        | Modifies an existing task                                                   | `task update 1 -title "New title" -due "2024-01-10 15:00"`         |
| `delete` | `<id>` (required)                                                                                                                              | Removes a task                                                              | `task delete 1`                                                    |

### Status and Colors

Tasks can have different statuses, each with its own visual indicator:

- `✓ Done`: Task is completed
- `! Overdue`: Task's due date has passed
- `→ DueSoon`: Task is due within 24 hours
- `⏰ Upcoming`: Task has an upcoming reminder
- `Pending`: Normal task status

### Priority Levels

Tasks can be assigned one of three priority levels:

- `high`: For urgent and important tasks (shown in red)
- `medium`: Default priority level (shown in yellow)
- `low`: For less urgent tasks (shown in green)

## Roadmap

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
