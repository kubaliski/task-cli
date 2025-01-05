package task

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const fileName = "tasks.json"

func (tm *TaskManager) SaveTasks() error {
	data, err := json.MarshalIndent(tm.tasks, "", "  ")
	if err != nil {
		return err
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	dataDir := filepath.Join(homeDir, ".task-cli")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(dataDir, fileName), data, 0644)
}

func (tm *TaskManager) LoadTasks() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	filepath := filepath.Join(homeDir, ".task-cli", fileName)
	data, err := os.ReadFile(filepath)
	if os.IsNotExist(err) {
		return nil // No file, no tasks we start with an empty list
	}

	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &tm.tasks); err != nil {
		return err
	}

	// Update nextID base on the highest ID find in the tasks
	for _, task := range tm.tasks {
		if task.ID >= tm.nextID {
			tm.nextID = task.ID + 1
		}
	}

	return nil
}
