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
}

type TaskManager struct {
	tasks  []Task
	nextID int
}

// NewTaskManager creates a new TaskManager
func NewTaskManager() *TaskManager {
	return &TaskManager{
		tasks:  make([]Task, 0),
		nextID: 1,
	}
}

// GetTasks returns all tasks in the TaskManager
func (tm *TaskManager) GetTasks() []Task {
	return tm.tasks
}

// GetTaskByID returns a task with a specific ID
func (tm *TaskManager) GetTaskByID(id int) (Task, error) {
	for _, task := range tm.tasks {
		if task.ID == id {
			return task, nil
		}
	}
	return Task{}, fmt.Errorf("task with id %d not found", id)

}

// AddTask adds a new task to the TaskManager
func (tm *TaskManager) AddTask(title string, priority TaskPriority) Task {
	// If priority is not set, use the default priority (Medium)
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

	tm.tasks = append(tm.tasks, task)
	tm.nextID++
	return task
}

// UpdateTask updates a task in the TaskManager
func (tm *TaskManager) UpdateTask(id int, title string, done bool, priority *TaskPriority) error {
	for i, task := range tm.tasks {
		if task.ID == id {
			// Update title
			if title != "" {
				tm.tasks[i].Title = title
			}

			//  Update priority
			if priority != nil {
				tm.tasks[i].Priority = *priority
			}

			// Update state
			tm.tasks[i].Done = done
			if done {
				tm.tasks[i].CompletedAt = time.Now()
			} else {
				tm.tasks[i].CompletedAt = time.Time{}
			}
			return nil
		}
	}
	return fmt.Errorf("task with ID %d not found", id)
}

// DeleteTask deletes a task from the TaskManager
func (tm *TaskManager) DeleteTask(id int) error {
	for i, task := range tm.tasks {
		if task.ID == id {
			// Delete task from slice using append
			tm.tasks = append(tm.tasks[:i], tm.tasks[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("task with id  %d not found", id)
}

// GetTasksSortedByPriority returns all tasks sorted by priority (high to low)
func (tm *TaskManager) GetTasksSorted(byPriority bool) []Task {
	sorted := make([]Task, len(tm.tasks))
	copy(sorted, tm.tasks)

	if byPriority {
		// Order by priority and then by ID
		sort.Slice(sorted, func(i, j int) bool {
			if sorted[i].Priority != sorted[j].Priority {
				return sorted[i].Priority > sorted[j].Priority
			}
			return sorted[i].ID < sorted[j].ID
		})
	} else {
		// Order only by ID
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].ID < sorted[j].ID
		})
	}

	return sorted
}
