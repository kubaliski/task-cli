package task

import (
	"time"
)

type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Done        bool      `json:"done"`
	CreatedAt   time.Time `json:"created_at"`
	CompletedAt time.Time `json:"completed_at"`
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

// AddTask adds a new task to the TaskManager
func (tm *TaskManager) AddTask(title string) Task {
	task := Task{
		ID:        tm.nextID,
		Title:     title,
		Done:      false,
		CreatedAt: time.Now(),
	}

	tm.tasks = append(tm.tasks, task)
	tm.nextID++
	return task
}

// GetTasks returns all tasks in the TaskManager
func (tm *TaskManager) GetTasks() []Task {
	return tm.tasks
}
