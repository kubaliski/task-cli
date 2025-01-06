package task

import (
	"fmt"
	"sort"
	"time"
)

type Task struct {
	ID          int          `json:"id"`
	Title       string       `json:"title"`
	Done        bool         `json:"done"`
	Priority    TaskPriority `json:"priority"`
	CreatedAt   time.Time    `json:"created_at"`
	CompletedAt time.Time    `json:"completed_at"`
	DueDate     *time.Time   `json:"due_date,omitempty"`
	Reminder    *time.Time   `json:"reminder,omitempty"`
	timeStatus  TimeStatus   `json:"-"` // Is calculated but it won't be shown in the JSON
}

// GetTimeStatus returns the time status of a task based on its due date and reminder
func (t *Task) GetTimeStatus() TimeStatus {
	if t.Done {
		return TimeStatusNormal
	}
	return GetTimeStatus(t.DueDate, t.Reminder)
}

// UpdateTimeStatus update the time status of a task
func (t *Task) UpdateTimeStatus() {
	t.timeStatus = t.GetTimeStatus()
}

// IsUpcoming shows if the task is upcoming
func (t *Task) IsUpcoming() bool {
	return t.timeStatus == TimeStatusUpcoming
}

// IsDueSoon shows if the task is due soon
func (t *Task) IsDueSoon() bool {
	return t.timeStatus == TimeStatusDueSoon
}

// IsOverdue show if the task is overdue
func (t *Task) IsOverdue() bool {
	return t.timeStatus == TimeStatusOverdue
}

type TaskManager struct {
	tasks  []Task
	nextID int
}

func NewTaskManager() *TaskManager {
	return &TaskManager{
		tasks:  make([]Task, 0),
		nextID: 1,
	}
}

// AddTask creates a new task and adds it to the task manager
func (tm *TaskManager) AddTask(title string, priority TaskPriority) Task {
	if priority == 0 {
		priority = DefaultPriority
	}

	task := Task{
		ID:        tm.nextID,
		Title:     title,
		Done:      false,
		Priority:  priority,
		CreatedAt: time.Now(),
	}

	task.UpdateTimeStatus()
	tm.tasks = append(tm.tasks, task)
	tm.nextID++
	return task
}

// SetDueDate stablish a due date for a task
func (tm *TaskManager) SetDueDate(id int, dueDate time.Time) error {
	for i, task := range tm.tasks {
		if task.ID == id {
			// Validar que el recordatorio (si existe) sea anterior a la fecha límite
			if err := ValidateTimeOrder(&dueDate, task.Reminder); err != nil {
				return err
			}
			tm.tasks[i].DueDate = &dueDate
			tm.tasks[i].UpdateTimeStatus()
			return nil
		}
	}
	return fmt.Errorf("task with ID %d not found", id)
}

// SetReminder stablish a reminder for a task
func (tm *TaskManager) SetReminder(id int, reminder time.Time) error {
	for i, task := range tm.tasks {
		if task.ID == id {
			// Validate that the reminder (if exists) is before the due date
			if err := ValidateTimeOrder(task.DueDate, &reminder); err != nil {
				return err
			}
			tm.tasks[i].Reminder = &reminder
			tm.tasks[i].UpdateTimeStatus()
			return nil
		}
	}
	return fmt.Errorf("task with ID %d not found", id)
}

// RemoveDueDate remove the due date of a task
func (tm *TaskManager) RemoveDueDate(id int) error {
	for i, task := range tm.tasks {
		if task.ID == id {
			tm.tasks[i].DueDate = nil
			tm.tasks[i].UpdateTimeStatus()
			return nil
		}
	}
	return fmt.Errorf("task with ID %d not found", id)
}

// RemoveReminder remove the reminder of a task
func (tm *TaskManager) RemoveReminder(id int) error {
	for i, task := range tm.tasks {
		if task.ID == id {
			tm.tasks[i].Reminder = nil
			tm.tasks[i].UpdateTimeStatus()
			return nil
		}
	}
	return fmt.Errorf("task with ID %d not found", id)
}

// GetTasksSorted returns the tasks sorted by priority and due date
func (tm *TaskManager) GetTasksSorted(byPriority bool, byDueDate bool) []Task {
	sorted := make([]Task, len(tm.tasks))
	copy(sorted, tm.tasks)

	// Update time status for each task
	for i := range sorted {
		sorted[i].UpdateTimeStatus()
	}

	sort.Slice(sorted, func(i, j int) bool {
		// First order by temporal state (overdue first)
		if sorted[i].timeStatus != sorted[j].timeStatus {
			return sorted[i].timeStatus > sorted[j].timeStatus
		}

		if byDueDate && sorted[i].DueDate != nil && sorted[j].DueDate != nil {
			// If both have a due date, order by due date
			return sorted[i].DueDate.Before(*sorted[j].DueDate)
		}

		if byPriority {
			// If both have the same priority, order by ID
			if sorted[i].Priority != sorted[j].Priority {
				return sorted[i].Priority > sorted[j].Priority
			}
		}

		// By default, order by ID
		return sorted[i].ID < sorted[j].ID
	})

	return sorted
}

// GetTasksByTimeStatus returns the tasks filtered by time status
func (tm *TaskManager) GetTasksByTimeStatus(status TimeStatus) []Task {
	var filtered []Task
	for _, task := range tm.tasks {
		task.UpdateTimeStatus()
		if task.timeStatus == status {
			filtered = append(filtered, task)
		}
	}
	return filtered
}

// UpdateTask update a task with new values
func (tm *TaskManager) UpdateTask(id int, title string, done bool, priority *TaskPriority) error {
	for i, task := range tm.tasks {
		if task.ID == id {
			if title != "" {
				tm.tasks[i].Title = title
			}
			if priority != nil {
				tm.tasks[i].Priority = *priority
			}

			// Update done status
			if done != task.Done {
				tm.tasks[i].Done = done
				if done {
					tm.tasks[i].CompletedAt = time.Now()
				} else {
					tm.tasks[i].CompletedAt = time.Time{}
				}
			}

			tm.tasks[i].UpdateTimeStatus()
			return nil
		}
	}
	return fmt.Errorf("task with ID %d not found", id)
}

// GetTaskByID returns a task by its ID
func (tm *TaskManager) GetTaskByID(id int) (Task, error) {
	for _, task := range tm.tasks {
		if task.ID == id {
			task.UpdateTimeStatus()
			return task, nil
		}
	}
	return Task{}, fmt.Errorf("task with ID %d not found", id)
}

// DeleteTask removes a task from the task manager
func (tm *TaskManager) DeleteTask(id int) error {
	for i, task := range tm.tasks {
		if task.ID == id {
			tm.tasks = append(tm.tasks[:i], tm.tasks[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("task with ID %d not found", id)
}
