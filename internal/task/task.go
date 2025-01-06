package task

import (
	"fmt"
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

// UpdateTask updates a task in the TaskManager
func (tm *TaskManager) UpdateTask(id int, title string, done bool) error {
	for i, task := range tm.tasks {
		if task.ID == id {
			// Actualizar el título
			if title != "" {
				tm.tasks[i].Title = title
			}

			// Actualizar el estado
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
			// Eliminar la tarea usando slice tricks
			tm.tasks = append(tm.tasks[:i], tm.tasks[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("task with id  %d not found", id)
}
