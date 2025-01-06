package task

import (
	"fmt"
	"time"
)

// TimeStatus represents the time status of a task
type TimeStatus int

const (
	TimeStatusNormal TimeStatus = iota
	TimeStatusUpcoming
	TimeStatusDueSoon
	TimeStatusOverdue
)

// DefaultReminderWindow is the default time window for reminders
const DefaultReminderWindow = 24 * time.Hour

// String returns the string representation of a TimeStatus
func (ts TimeStatus) String() string {
	switch ts {
	case TimeStatusNormal:
		return "Normal"
	case TimeStatusUpcoming:
		return "Upcoming"
	case TimeStatusDueSoon:
		return "Due Soon"
	case TimeStatusOverdue:
		return "Overdue"
	default:
		return "Unknown"
	}
}

// Color returns the color code for a TimeStatus
func (ts TimeStatus) Color() string {
	switch ts {
	case TimeStatusNormal:
		return "\033[0m" // Reset
	case TimeStatusUpcoming:
		return "\033[0;34m" // Blue
	case TimeStatusDueSoon:
		return "\033[0;33m" // Yellow
	case TimeStatusOverdue:
		return "\033[0;31m" // Red
	default:
		return "\033[0m" // Reset
	}
}

// ParseDateTime parses a string and returns the corresponding time.Time
func ParseDateTime(s string) (time.Time, error) {
	// multiple date time formats
	formats := []string{
		"2006-01-02 15:04",
		"2006-01-02T15:04",
		"2006/01/02 15:04",
		"02/01/2006 15:04",
		"02-01-2006 15:04",
	}

	var firstErr error
	for _, format := range formats {
		if t, err := time.Parse(format, s); err == nil {
			return t, nil
		} else if firstErr == nil {
			firstErr = err
		}
	}
	return time.Time{}, fmt.Errorf("invalid date time format: %v", firstErr)
}

// FormatDateTime returns a formatted string representation of a time.Time
func FormatDateTime(t *time.Time) string {
	if t == nil || t.IsZero() {
		return "---"
	}
	return t.Format("2006-01-02 15:04")
}

// GetTimeStatus returns the time status of a task based on its due date and reminder
func GetTimeStatus(dueDate *time.Time, reminder *time.Time) TimeStatus {
	if dueDate == nil {
		return TimeStatusNormal
	}

	now := time.Now()

	// If the task is overdue
	if dueDate.Before(now) {
		return TimeStatusOverdue
	}

	// If the task is due soon(less than 24 hours)
	if dueDate.Sub(now) <= DefaultReminderWindow {
		return TimeStatusDueSoon
	}

	// If the task have an upcoming reminder
	if reminder != nil && reminder.After(now) && reminder.Sub(now) <= DefaultReminderWindow {
		return TimeStatusUpcoming
	}

	return TimeStatusNormal

}

// ValidateTimeOrder validates the order of due date and reminder
func ValidateTimeOrder(dueDate, reminder *time.Time) error {
	if dueDate == nil || reminder == nil {
		return nil
	}

	if reminder.After(*dueDate) {
		return fmt.Errorf("reminder time cannot be after due date")
	}

	return nil
}
