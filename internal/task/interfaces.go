// task/interfaces.go
package task

import "time"

// TaskManager define la interfaz para manejar las operaciones de tareas
type ITaskManager interface {
	// Operaciones básicas
	AddTask(title string, priority TaskPriority) Task
	GetTaskByID(id int) (Task, error)
	DeleteTask(id int) error
	UpdateTask(id int, title string, done bool, priority *TaskPriority) error

	// Manejo de fechas y recordatorios
	SetDueDate(id int, date time.Time) error
	SetReminder(id int, date time.Time) error
	RemoveDueDate(id int) error
	RemoveReminder(id int) error

	// Consultas y listados
	GetTasksSorted(byPriority, byDueDate bool) []Task
	GetTasksByTimeStatus(status TimeStatus) []Task

	// Persistencia
	LoadTasks() error
	SaveTasks() error
}

// Verificamos que TaskManager implementa ITaskManager
var _ ITaskManager = (*TaskManager)(nil)
