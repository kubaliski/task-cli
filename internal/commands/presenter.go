package commands

import (
	"fmt"
	"strings"
	"task-cli/internal/task"
	"unicode/utf8"
)

// TableColumn define a column in the table
type TableColumn struct {
	Header string
	Width  int
	Get    func(t task.Task) string
}

// DefaultPresenter implement the default presenter
type DefaultPresenter struct {
	columns []TableColumn
}

// NewDefaultPresenter create a new instance of DefaultPresenter
func NewDefaultPresenter() *DefaultPresenter {
	p := &DefaultPresenter{}
	p.initializeColumns()
	return p
}

// initializeColumns initialize the columns for the table
func (p *DefaultPresenter) initializeColumns() {
	p.columns = []TableColumn{
		{
			Header: "ID",
			Width:  2,
			Get: func(t task.Task) string {
				return fmt.Sprintf("%d", t.ID)
			},
		},
		{
			Header: "Status",
			Width:  9,
			Get: func(t task.Task) string {
				return getStatusString(t)
			},
		},
		{
			Header: "Priority",
			Width:  7,
			Get: func(t task.Task) string {
				priority := t.Priority.String()
				if priority == "" {
					return ""
				}
				return fmt.Sprintf("%s%s%s",
					t.Priority.Color(),
					priority,
					"\033[0m")
			},
		},
		{
			Header: "Title",
			Width:  40,
			Get: func(t task.Task) string {
				return truncateString(t.Title, 37)
			},
		},
		{
			Header: "Due Date",
			Width:  16,
			Get: func(t task.Task) string {
				date := task.FormatDateTime(t.DueDate)
				if t.IsOverdue() {
					return task.TimeStatusOverdue.Color() + date + "\033[0m"
				}
				return date
			},
		},
		{
			Header: "Reminder",
			Width:  16,
			Get: func(t task.Task) string {
				reminder := task.FormatDateTime(t.Reminder)
				if t.IsUpcoming() {
					return task.TimeStatusUpcoming.Color() + reminder + "\033[0m"
				}
				return reminder
			},
		},
		{
			Header: "Created",
			Width:  16,
			Get: func(t task.Task) string {
				return t.CreatedAt.Format("2006-01-02 15:04")
			},
		},
	}
}

// centerText center selected text in a string
func centerText(text string, width int) string {
	cleanText := stripANSI(text)
	textLen := utf8.RuneCountInString(cleanText)
	if textLen >= width {
		return text
	}

	leftPad := (width - textLen) / 2
	rightPad := width - textLen - leftPad

	return strings.Repeat(" ", leftPad) + text + strings.Repeat(" ", rightPad)
}

// truncateString truncate a string to a maximum length (utf8 safe)
func truncateString(s string, maxWidth int) string {
	if utf8.RuneCountInString(s) <= maxWidth {
		return s
	}
	runes := []rune(s)
	return string(runes[:maxWidth-3]) + "..."
}

// stripANSI remove ANSI escape codes from a string (for text width calculation)
func stripANSI(str string) string {
	var result strings.Builder
	inEscape := false

	for _, ch := range str {
		if ch == '\033' {
			inEscape = true
			continue
		}
		if inEscape {
			if ch == 'm' {
				inEscape = false
			}
			continue
		}
		result.WriteRune(ch)
	}

	return result.String()
}

// PrintTaskTable implement the table format view
func (p *DefaultPresenter) PrintTaskTable(tasks []task.Task) error {
	if len(tasks) == 0 {
		return nil
	}

	// Prepare separator
	separator := "+"
	for _, col := range p.columns {
		separator += strings.Repeat("-", col.Width+2) + "+"
	}

	// Print superior separator
	fmt.Println(separator)

	// Print Headers
	for i, col := range p.columns {
		if i == 0 {
			fmt.Print("|")
		}
		fmt.Printf(" %s |", centerText(col.Header, col.Width))
	}
	fmt.Println()

	// Print separator after headers
	fmt.Println(separator)

	// Print rows
	for _, t := range tasks {
		for i, col := range p.columns {
			if i == 0 {
				fmt.Print("|")
			}
			value := col.Get(t)
			// Make sure the value fits in the column
			cleanValue := stripANSI(value)
			if utf8.RuneCountInString(cleanValue) > col.Width {
				value = truncateString(cleanValue, col.Width)
			}
			fmt.Printf(" %s |", centerText(value, col.Width))
		}
		fmt.Println()
	}

	// Print inferior separator
	fmt.Println(separator)

	return nil
}

// PrintTaskList implement the list format view
func (p *DefaultPresenter) PrintTaskList(tasks []task.Task) error {
	for _, t := range tasks {
		if err := p.PrintTask(t); err != nil {
			return err
		}
		fmt.Println(strings.Repeat("-", 50))
	}
	return nil
}

// PrintTask implement the task format view (detailed)
func (p *DefaultPresenter) PrintTask(t task.Task) error {
	priorityStr := fmt.Sprintf("%s%s%s",
		t.Priority.Color(),
		t.Priority.String(),
		"\033[0m")

	fmt.Printf("\n%s Task #%d: %s - %s\n",
		getStatusIcon(t),
		t.ID,
		priorityStr,
		t.Title)

	fmt.Printf("   Created: %s\n", t.CreatedAt.Format("2006-01-02 15:04:05"))

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
		fmt.Printf("   Completed: %s\n", t.CompletedAt.Format("2006-01-02 15:04:05"))
	}

	return nil
}

// PrintSuccess print a success message
func (p *DefaultPresenter) PrintSuccess(format string, a ...interface{}) {
	fmt.Printf(format+"\n", a...)
}

// PrintError print an error message
func (p *DefaultPresenter) PrintError(format string, a ...interface{}) error {
	return fmt.Errorf(format, a...)
}

// getStatusString returns formatted status string
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

// getStatusIcon returns the status icon for a task
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
