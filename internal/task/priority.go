package task

import (
	"fmt"
)

// TaskPriority represents the priority of a task
type TaskPriority int

// Task priorities constants default is Medium
const (
	PriorityLow TaskPriority = iota
	PriorityMedium
	PriorityHigh

	DefaultPriority = PriorityMedium
)

// String returns the string representation of a TaskPriority
func (p TaskPriority) String() string {
	switch p {
	case PriorityLow:
		return "Low"
	case PriorityMedium:
		return "Medium"
	case PriorityHigh:
		return "High"
	default:
		return "Unknown"
	}
}

// ParsePriority parses a string and returns the corresponding TaskPriority
func ParsePriority(s string) (TaskPriority, error) {
	switch s {
	case "low", "Low", "LOW":
		return PriorityLow, nil
	case "medium", "Medium", "MEDIUM":
		return PriorityMedium, nil
	case "high", "High", "HIGH":
		return PriorityHigh, nil
	default:
		return DefaultPriority, fmt.Errorf("unknown priority: %s", s)
	}
}

// Color returns the color code for a TaskPriority
func (p TaskPriority) Color() string {
	switch p {
	case PriorityLow:
		return "\033[0;32m" // Green
	case PriorityMedium:
		return "\033[0;33m" // Yellow
	case PriorityHigh:
		return "\033[0;31m" // Red
	default:
		return "\033[0m" // Reset
	}
}
